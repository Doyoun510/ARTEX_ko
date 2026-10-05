package agent

import (
	"log"
	"path/filepath"

	"github.com/Autumn-27/norma/agentcore"
	"github.com/Autumn-27/norma/noaadapter"
)

// noaWarn returns a diagnostics sink tagging non-fatal noa messages with the
// session, routed through the package logger (agents have no per-instance one).
func noaWarn(session string) func(string) {
	return func(msg string) { log.Printf("[noa] %s: %s", session, msg) }
}

// noa는 norma v0.4.0이 도입한 '모델 기반 컨텍스트 압축' 메커니즘으로, 플랫폼 실험 기능으로서 사용자가
// 시스템 설정에서 켜고 끈다. 내장 compaction과 상호 배타적: noaadapter.Enable이 유일한 진입점으로, 한 번에
// 컨텍스트 인계기(Compactor)·Compress 도구·3단 상주 프롬프트를 달며, Enable을 호출하지 않으면 꺼진 상태
// (내장 compaction이 그대로 동작). 스위치는 각 agent가 주입한 noaEnabledFn이 해석하며 run마다 한 번
// 읽으므로, 전환은 이후 시작되는 run에만 영향을 주고 agent를 재생성할 필요가 없다.

// enableNoa는 해석기가 켜짐을 보고하면 noa를 opts에 연결한다. archiveRoot는 압축 원문의 영속화 기준 디렉터리
// (전역 workDir 사용, 각 agent가 <workDir>/noa 아래에 통일돼 작업/의도 디렉터리로 흩어지지 않음), sessionID가
// 그 아래 아카이브 서브디렉터리를 명명한다(전역 유일하므로 같은 기준 디렉터리 안에서 충돌 없음).
//
// noa는 실험 기능: 연결 실패가 실제 작업을 중단시켜선 안 된다. 오류 발생 시 onWarn으로 보고하고 내장 압축으로 회귀한다.
// 활성화 성공 시 opts.Compaction을 비워, agentcore가 '컨텍스트 관리자 두 개 동시 설정'으로 경고하지 않게 한다.
func enableNoa(opts *agentcore.Options, enabled func() bool, archiveRoot, sessionID string, onWarn func(string)) {
	if enabled == nil || !enabled() {
		return
	}
	if opts.OnWarn == nil {
		opts.OnWarn = onWarn
	}
	if err := noaadapter.Enable(opts, noaadapter.Options{
		ArchiveBaseDir: filepath.Join(archiveRoot, "noa"),
		SessionID:      sessionID,
		OnWarn:         onWarn,
	}); err != nil {
		if onWarn != nil {
			onWarn("noa 압축 활성화 실패, 내장 압축으로 회귀: " + err.Error())
		}
		return
	}
	// Compactor가 Compaction을 덮어쓰지만, 둘이 공존하면 agentcore가 매번 경고한다; 명시적으로 비운다.
	opts.Compaction = nil
}
