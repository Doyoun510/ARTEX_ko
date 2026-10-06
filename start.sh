#!/bin/sh
# ARTEX 데몬 시작 스크립트(Linux / macOS / Docker ENTRYPOINT)
#
# 사용법:
#   ./start.sh                       포그라운드 실행(Ctrl-C로 중지)
#   nohup ./start.sh >artex.log 2>&1 &   백그라운드 상시 실행
#   ./start.sh -addr :9000           추가 파라미터를 artex에 그대로 전달
#
# 하는 일은 하나뿐이다: artex를 실행하고, 프로세스 종료 후 종료 코드에 따라 재시작 여부를 결정한다.
#
#   0      사용자가 정상 중지        → 루프 종료
#   75     프로그램이 재시작 요청        → 즉시 재실행(화면에서 "원클릭 업데이트" 또는 "롤백" 클릭)
#   기타   비정상 종료                → 백오프 후 재실행(1→2→4…최대 60초)
#
# 다운로드·SHA256 검증·버전 교체는 의도적으로 여기서 하지 않는다: 이 로직은 sh와 bat에 각각 작성해야 하는데,
# 바로 이 부분에서 오류가 나면 안 되기 때문이다. 실행할 수 없는 바이너리로 교체되면 이 스크립트는 계속
# 그 바이너리를 재시작하며, 사용자가 해당 시스템에서 직접 복구해야 한다. 따라서 검증/교체는 전부 Go(selfupdate 패키지)에 맡기고,
# artex가 시작 시 직접 수행하며, 스크립트는 단순하게 유지한다.
set -u

cd "$(dirname "$0")" || exit 1

BIN=./artex
[ -x "$BIN" ] || { echo "[artex] 실행 파일을 찾을 수 없음: $BIN" >&2; exit 1; }

RESTART_CODE=75
MAX_DELAY=60

child=0
stopping=0

# 중지 신호를 artex에 전달한다.
#
# Docker에서는 필수다: docker stop은 PID 1(이 스크립트)에만 SIGTERM을 보내고,
# 자식 프로세스에는 보내지 않는다. 전달하지 않으면 artex가 신호를 받지 못해 정상 종료 절차를 수행할 수 없고 10초 후 SIGKILL로
# 강제 종료되어 실행 중인 작업이 중간에 끊긴다.
forward() {
	stopping=1
	if [ "$child" -ne 0 ]; then
		kill -TERM "$child" 2>/dev/null || true
	fi
}
trap forward INT TERM

delay=1
while :; do
	"$BIN" "$@" &
	child=$!

	# 신호는 wait를 중단시켜 >128을 반환하게 한다. 이때 자식 프로세스는 아직 정상 종료 절차를 수행 중이므로,
	# 실제 종료 코드를 얻으려면 반드시 wait를 한 번 더 호출해야 한다.
	wait "$child"
	code=$?
	if [ "$code" -gt 128 ]; then
		wait "$child"
		code=$?
	fi
	child=0

	if [ "$stopping" -eq 1 ]; then
		echo "[artex] 중지됨"
		exit 0
	fi

	case "$code" in
		0)
			echo "[artex] 정상 종료"
			exit 0
			;;
		"$RESTART_CODE")
			# 업데이트/롤백 준비 완료: 재실행 후 artex가 시작 시 버전 교체를 완료한다(selfupdate.Bootstrap 참고).
			echo "[artex] 재시작 요청(새 버전 적용)…"
			delay=1
			;;
		*)
			echo "[artex] 비정상 종료 (code=$code), ${delay}s 후 재시작" >&2
			sleep "$delay"
			delay=$((delay * 2))
			[ "$delay" -gt "$MAX_DELAY" ] && delay=$MAX_DELAY
			;;
	esac
done
