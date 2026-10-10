#!/usr/bin/env bash
# norma(github.com/Autumn-27/norma)의 OpenAI 어댑터에 Gemini thought_signature
# 라운드트립 지원을 주입한다. norma는 공개 저장소가 아니라 업스트림 PR이 불가능하므로,
# go.mod가 가리키는 버전의 모듈을 받아 third_party/norma에 펼친 뒤 patches/의 패치를
# 적용하고, go.mod의 `replace github.com/Autumn-27/norma => ./third_party/norma`가
# 이 디렉터리를 쓰게 한다. third_party/norma는 .gitignore 대상이라(라이선스 미표기
# 3자 코드 비공개 유지) 빌드 전에 이 스크립트를 한 번 실행해야 한다.
#
# 사용: scripts/apply-norma-patch.sh
# 요구: go, git(또는 patch)
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

MOD="github.com/Autumn-27/norma"
PATCH="$ROOT/patches/norma-gemini-thought-signature.patch"
DEST="$ROOT/third_party/norma"

# go.mod가 요구하는 정확한 버전을 읽는다(패치는 해당 버전 소스에 맞춰져 있음).
VER="$(go list -m -f '{{.Version}}' "$MOD" 2>/dev/null || true)"
if [ -z "$VER" ]; then
  VER="$(awk '/[[:space:]]'"${MOD//\//\\/}"'[[:space:]]/{print $2}' go.mod | head -1)"
fi
[ -n "$VER" ] || { echo "ERROR: go.mod에서 $MOD 버전을 찾지 못함"; exit 1; }
echo "[norma-patch] 대상 $MOD $VER"

# 모듈 소스를 받아 캐시 경로를 찾는다.
GOFLAGS=-mod=mod go mod download "$MOD@$VER"
SRC="$(go env GOMODCACHE)/$(printf '%s' "$MOD@$VER" | sed 's#\([A-Z]\)#!\L\1#g')"
[ -d "$SRC" ] || SRC="$(find "$(go env GOMODCACHE)" -maxdepth 3 -type d -iname "norma@$VER" 2>/dev/null | head -1)"
[ -d "$SRC" ] || { echo "ERROR: 모듈 캐시에서 norma 소스를 찾지 못함"; exit 1; }

# third_party/norma를 깨끗이 재구성(캐시는 읽기전용이므로 복사 후 쓰기권한 부여).
rm -rf "$DEST"
mkdir -p "$DEST"
cp -a "$SRC"/. "$DEST"/
chmod -R u+w "$DEST"

# 패치 적용(이미 적용돼 있으면 조용히 넘어감).
if git apply --check -p1 --directory=third_party/norma "$PATCH" 2>/dev/null; then
  git apply -p1 --directory=third_party/norma "$PATCH"
elif (cd "$DEST" && patch -p1 --dry-run < "$PATCH" >/dev/null 2>&1); then
  (cd "$DEST" && patch -p1 < "$PATCH")
else
  if grep -q "ThoughtSignature" "$DEST/llm/openai.go" 2>/dev/null; then
    echo "[norma-patch] 이미 적용됨 — 건너뜀"
  else
    echo "ERROR: 패치 적용 실패(버전 불일치 가능). go.mod의 norma 버전과 patches/ 확인"; exit 1
  fi
fi

grep -q "ThoughtSignature" "$DEST/llm/openai.go" \
  && echo "[norma-patch] 완료: third_party/norma 준비됨 (Gemini thought_signature 지원)" \
  || { echo "ERROR: 적용 검증 실패"; exit 1; }
