@echo off
rem 콘솔을 UTF-8로 전환한다. 그렇지 않으면 이 파일의 중국어가 GBK 터미널에서 깨진다.
chcp 65001 >nul 2>&1
rem ARTEX 데몬 시작 스크립트(Windows)
rem
rem 사용법:
rem   start.bat                  포그라운드 실행(Ctrl-C로 중지)
rem   start.bat -addr :9000      추가 파라미터를 artex에 그대로 전달
rem
rem 하는 일은 하나뿐이다: artex.exe를 실행하고, 프로세스 종료 후 종료 코드에 따라 재시작 여부를 결정한다.
rem
rem   0      사용자가 정상 중지     -> 루프 종료
rem   75     프로그램이 재시작 요청     -> 즉시 재실행(화면에서 "원클릭 업데이트" 또는 "롤백" 클릭)
rem   기타   비정상 종료             -> 백오프 후 재실행(1->2->4…최대 60초)
rem
rem 다운로드·SHA256 검증·버전 교체는 여기서 하지 않고 전부 artex가 시작 시 직접 수행한다
rem (selfupdate 패키지). 스크립트는 단순하게 유지하며, 자세한 내용은 start.sh 상단 설명을 참고한다.

setlocal enabledelayedexpansion
cd /d "%~dp0"

set "BIN=artex.exe"
if not exist "%BIN%" (
	echo [artex] 실행 파일을 찾을 수 없음: %BIN% 1>&2
	exit /b 1
)

set "RESTART_CODE=75"
set "MAX_DELAY=60"
set /a delay=1

:loop
"%BIN%" %*
set "code=!ERRORLEVEL!"

if "!code!"=="0" (
	echo [artex] 정상 종료
	exit /b 0
)

if "!code!"=="%RESTART_CODE%" (
	rem 업데이트/롤백 준비 완료: 재실행 후 artex가 시작 시 버전 교체를 완료한다.
	echo [artex] 재시작 요청: 새 버전 적용…
	set /a delay=1
	goto loop
)

echo [artex] 비정상 종료 ^(code=!code!^), !delay!s 후 재시작 1>&2
rem timeout은 리디렉션된 콘솔에서 실패하므로 ping을 대체 대기 수단으로 사용한다(N초 지연에는 N+1회 필요).
set /a pings=!delay!+1
ping -n !pings! 127.0.0.1 >nul 2>&1
set /a delay=!delay!*2
if !delay! gtr %MAX_DELAY% set /a delay=%MAX_DELAY%
goto loop
