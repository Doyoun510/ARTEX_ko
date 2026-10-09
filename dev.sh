#!/usr/bin/env bash
# 개발 모드: 백엔드(:8787) + 트래픽 프록시(:8788)와 프런트엔드 next dev(:5173)를 함께 실행한다.
# 프런트엔드 /api는 백엔드로 리버스 프록시하며, Ctrl-C로 함께 종료한다.
#
# 단일 바이너리(프런트엔드 임베드) 방식은 README의 '방법 4: 소스에서 단일 바이너리 컴파일' 절을 참고하며, 이 스크립트를 사용하지 않는다.
set -euo pipefail
cd "$(dirname "$0")"

# 종료 시 이 프로세스 그룹의 모든 자식 프로세스(백엔드 + 프런트엔드)를 종료한다.
cleanup() { kill 0 2>/dev/null || true; }
trap cleanup EXIT INT TERM

# 백엔드(일반 go run이며 프런트엔드를 임베드하지 않음). 동시 실행할 work agent 수는 '시스템 설정'에서 설정한다.
go run ./cmd/artex -addr :8787 -proxy 127.0.0.1:8788 &

# 프런트엔드 Hot Reload(Vite/Next dev server, /api는 :8787로 리버스 프록시).
( cd web && npm run dev ) &

echo "[dev] 백엔드 :8787 / 프록시 :8788 / 프런트엔드 http://localhost:5173  (Ctrl-C로 종료)"
wait
