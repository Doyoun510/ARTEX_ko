#!/usr/bin/env bash
# ARTEX 업데이트 스크립트: ① Docker 업데이트(새 이미지를 가져와 다시 생성)  ② 로컬 컴파일 업데이트(바이너리 다시 생성)
# install.sh와 대응: install은 최초 설치를, update는 새 버전으로 업그레이드를 담당한다.
# DB 마이그레이션을 직접 실행할 필요 없이 artex가 시작할 때마다 schema.sql을 멱등하게 다시 실행(ADD COLUMN/CREATE
# INDEX IF NOT EXISTS 포함)하므로 “재시작 시 마이그레이션”된다. 데이터(pgdata 볼륨·./data·./skills)는 영향을 받지 않는다.
set -euo pipefail
cd "$(cd "$(dirname "$0")" && pwd)"

info(){ printf '\033[36m[*]\033[0m %s\n' "$*"; }
ok(){   printf '\033[32m[+]\033[0m %s\n' "$*"; }
warn(){ printf '\033[33m[!]\033[0m %s\n' "$*"; }
die(){  printf '\033[31m[x]\033[0m %s\n' "$*" >&2; exit 1; }
ask(){  local p="$1" d="${2:-}" a; read -rp "$p${d:+ [$d]}: " a; echo "${a:-$d}"; }

# ── 선택 사항: 저장소를 최신 코드로 동기화(compose/스크립트/로컬 컴파일 소스 코드를 모두 이 방식으로 업데이트)───────
sync_repo(){
  [ -d .git ] && command -v git >/dev/null 2>&1 || { warn "git 작업 사본이 아니므로 git pull을 건너뜁니다"; return; }
  [ "$(ask '최신 코드를 가져오시겠습니까 (git pull --ff-only)? (y/n)' y)" = y ] || return
  if ! git pull --ff-only; then
    warn "git pull을 fast-forward로 진행할 수 없습니다(로컬 변경 또는 브랜치 분기). 직접 처리한 후 다시 시도하세요. 이번에는 현재 코드를 사용합니다"
  fi
}

# ── ① Docker 업데이트 ───────────────────────────────
update_docker(){
  command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1 \
    || die "docker / docker compose를 찾을 수 없습니다. 먼저 ./install.sh로 설치·배포하세요"
  [ -f .env ] || die ".env를 찾을 수 없습니다. 먼저 ./install.sh를 실행해 최초 배포를 완료하세요"

  # 선택 사항: 지정한 버전 tag로 업그레이드(입력하지 않으면 .env의 ARTEX_TAG를 유지하며 기본값은 latest)
  local tag; tag="$(ask '대상 이미지 tag(엔터를 누르면 .env / latest 유지)' '')"
  if [ -n "$tag" ]; then
    if grep -q '^ARTEX_TAG=' .env; then
      sed -i.bak "s|^ARTEX_TAG=.*|ARTEX_TAG=${tag}|" .env && rm -f .env.bak
    else
      printf '\nARTEX_TAG=%s\n' "$tag" >> .env
    fi
    ok "ARTEX_TAG 설정 완료: ${tag}"
  fi

  # artex만 변경한다: postgres는 16-alpine으로 고정되어 함께 업그레이드할 필요가 없다(이를 가져오면 대역폭만 낭비하고,
  # 주요 버전 변경은 호환성 위험도 있다). artex에는 depends_on postgres가 선언되어 있어 서비스 이름과 함께
  # up을 실행하면 pg가 실행 중이 아닐 때 자동 시작하며, 이미 실행 중이면 그대로 유지하고 다시 생성하지 않는다.
  info "새 이미지 가져오기(artex만)…"
  docker compose pull artex
  info "다시 생성 후 시작(artex 재시작 시 schema 자동 마이그레이션)…"
  docker compose up -d artex
  ok "업데이트 완료 → http://localhost:8787"
  info "로그 확인: docker compose logs -f artex"
  info "이전 이미지 정리(선택 사항): docker image prune -f"
}

# ── ② 로컬 컴파일 업데이트 ──────────────────────────────
update_local(){
  command -v go >/dev/null 2>&1 || die "Go(>=1.26)를 찾을 수 없습니다: https://go.dev/dl/"
  [ -f config.json ] || warn "config.json을 찾을 수 없습니다. 최초 배포라면 ./install.sh를 사용하세요"
  ok "Go: $(go version)"

  if command -v npm >/dev/null 2>&1; then
    info "프런트엔드 정적 산출물 다시 빌드…"
    ( cd web && npm ci && npm run build:static )
    rm -rf server/webui/dist && cp -r web/out server/webui/dist
    info "임베드된 단일 바이너리 다시 컴파일…"
    CGO_ENABLED=0 go build -tags embedui -trimpath -o artex ./cmd/artex
  else
    warn "npm을 찾을 수 없습니다: **프런트엔드를 임베드하지 않은** 백엔드를 컴파일합니다(프런트엔드는 npm run dev로 별도 실행 필요)"
    CGO_ENABLED=0 go build -o artex ./cmd/artex
  fi
  ok "컴파일 완료 → ./artex"
  warn "적용하려면 실행 중인 artex 프로세스를 재시작하세요(재시작 시 schema 자동 마이그레이션)"
}

echo "=============================="
echo "  ARTEX 업데이트"
echo "  1) Docker 업데이트(새 이미지를 가져와 다시 생성)"
echo "  2) 로컬 업데이트(go 다시 컴파일)"
echo "=============================="
case "$(ask '선택' 1)" in
  1) sync_repo; update_docker ;;
  2) sync_repo; update_local ;;
  *) die "유효하지 않은 선택" ;;
esac
