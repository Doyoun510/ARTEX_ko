# syntax=docker/dockerfile:1
#
# 실행 이미지(이미지 안에서 컴파일하지 않음): 자주 쓰는 도구만 설치하고 **미리 컴파일한 Linux 단일 바이너리**를 배치한다.
# 바이너리는 CI의 binaries job에서 교차 컴파일(순수 Go, QEMU 없음)하며 대상 아키텍처에 따라
# 빌드 컨텍스트의 dist/<TARGETARCH>/artex에 배치한다. 다중 아키텍처 빌드에서는 arm64가 apt 계층만 에뮬레이션하면 되어,
# Next/Go 컴파일을 에뮬레이션하지 않아 훨씬 빠르다.
#
# 로컬에서 이미지를 직접 빌드할 때는 먼저 바이너리를 준비한다:
#   cd web && npm run build:static && cd ..
#   cp -r web/out server/webui/dist
#   CGO_ENABLED=0 GOARCH=amd64 go build -tags embedui -o dist/amd64/artex ./cmd/artex
#   docker build -t artex:local .
FROM python:3.12-slim-bookworm
ARG TARGETARCH
# 자주 쓰는 도구: ripgrep / curl / vim에 recon 기본 도구를 추가(필요에 따라 추가·삭제).
# Node는 NodeSource에서 20.x를 설치: bookworm에 포함된 apt nodejs는 18이며, Playwright는 >=20이 필요하다.
RUN apt-get update && apt-get install -y --no-install-recommends \
      ca-certificates ripgrep curl wget vim git jq unzip \
      dnsutils iputils-ping netcat-openbsd inetutils-telnet whois nmap \
    && curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
    && apt-get install -y --no-install-recommends nodejs \
    && rm -rf /var/lib/apt/lists/*
# Playwright MCP와 CLI를 미리 설치(전역)해 런타임에 npx로 네트워크 다운로드를 하지 않는다.
# @playwright/mcp: browser MCP는 바로 `npx @playwright/mcp` 실행(이미 전역 설치되어 -y/@latest 불필요).
# @playwright/cli: playwright-cli를 제공하며, 설치 후 --help로 실행 가능 여부도 확인한다.
# playwright도 설치(브라우저 관리 제공)하고, 설치 후 --with-deps로 chromium과 시스템 의존성을 미리 설치해,
# 컨테이너 안에서 MCP/CLI를 처음 시작할 때 바로 사용할 수 있으며 브라우저를 네트워크로 다운로드하지 않는다.
RUN npm install -g @playwright/mcp@latest @playwright/cli@latest playwright@latest \
    && playwright-cli --help \
    && playwright install --with-deps chromium \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /app
# 해당 아키텍처용으로 미리 컴파일한 바이너리(dist/amd64/artex 또는 dist/arm64/artex)
COPY dist/${TARGETARCH}/artex /app/artex
# 데몬 시작 스크립트: 프로세스 종료 후 종료 코드에 따라 재시작 여부를 결정하며, 화면의 원클릭 업데이트는 이 스크립트로 버전 교체를 완료한다.
# SIGTERM도 artex에 전달한다. docker stop은 PID 1에만 신호를 보내므로,
# 전달하지 않으면 artex가 신호를 받지 못해 정상 종료 절차를 수행할 수 없고, 10초 후 SIGKILL로 강제 종료된다.
COPY start.sh /app/start.sh
RUN chmod +x /app/artex /app/start.sh
COPY skills/ /app/skills/
# data/(SQLite + jwt.key) 영구 저장 위치
VOLUME ["/app/data"]
EXPOSE 8787 8788
ENTRYPOINT ["/app/start.sh"]
CMD ["-addr", ":8787", "-proxy", ":8788"]
