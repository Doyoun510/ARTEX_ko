package notify

import (
	"errors"
	"strings"
	"testing"
)

func TestMaskedValueHidesBodyButKeepsTailHint(t *testing.T) {
	const secret = "https://oapi.dingtalk.com/robot/send?access_token=abcdef123456"
	got := MaskedValue(secret)
	if strings.Contains(got, "abcdef123456") {
		t.Fatalf("마스킹된 값에 전체 자격 증명 노출: %q", got)
	}
	if strings.Contains(got, "oapi.dingtalk.com") {
		t.Fatalf("마스킹된 값에 주소 본체가 노출되어서는 안 됩니다: %q", got)
	}
	// 어떤 봇인지 사용자가 알아볼 수 있도록 마지막 6자리를 유지해야 합니다.
	if !strings.HasSuffix(got, "123456") {
		t.Fatalf("식별 힌트로 마지막 6자리를 유지해야 합니다: %q", got)
	}
	if !IsMasked(got) {
		t.Fatalf("마스킹된 값은 반드시 IsMasked가 식별할 수 있어야 합니다: %q", got)
	}
}

func TestMaskedValueShortSecretGivesNoHint(t *testing.T) {
	// 짧은 자격 증명도 마지막 6자리를 노출하면 자격 증명 전체를 노출하는 셈입니다.
	for _, s := range []string{"abc", "abcdef", ""} {
		got := MaskedValue(s)
		if got != MaskedPrefix {
			t.Fatalf("길이 %d인 자격 증명에는 끝부분 힌트를 주어서는 안 됩니다. 실제: %q", len(s), got)
		}
		if s != "" && strings.Contains(got, s) {
			t.Fatalf("마스킹된 값에 원래 값 포함: %q", got)
		}
	}
}

func TestMaskConfigMasksOnlySecrets(t *testing.T) {
	cfg := map[string]any{
		"webhook": "https://example.com/hook?token=SECRETVALUE",
		"secret":  "SECtest123456",
		"port":    float64(587),
		"host":    "smtp.example.com",
	}
	masked := MaskConfig(KindDingTalk, cfg)
	for _, k := range []string{"webhook", "secret"} {
		s, _ := masked[k].(string)
		if !IsMasked(s) {
			t.Errorf("%s: 마스킹되어야 합니다. 실제: %q", k, s)
		}
	}
	// 자격 증명이 아닌 필드는 반드시 그대로 유지해야 합니다. 그렇지 않으면 UI에서 표시할 수 없습니다.
	if masked["port"] != float64(587) {
		t.Errorf("자격 증명이 아닌 port 필드를 변경해서는 안 됩니다: %v", masked["port"])
	}
}

func TestMaskConfigUnknownKindReturnsEmpty(t *testing.T) {
	// 채널 유형을 식별할 수 없으면 UI에 빈 설정을 표시하더라도, 자격 증명을 포함할 수 있는 원래 내용을 반환해서는 안 됩니다.
	got := MaskConfig("nope", map[string]any{"webhook": "https://x/y?token=LEAK"})
	if len(got) != 0 {
		t.Fatalf("알 수 없는 채널 유형은 빈 설정을 반환해야 합니다. 실제: %v", got)
	}
}

func TestMaskConfigDoesNotMutateInput(t *testing.T) {
	// 마스킹은 표시 계층의 동작이며, DB의 실제 값을 바꿔서는 안 됩니다.
	cfg := map[string]any{"webhook": "https://example.com/hook", "secret": "SECtest123456"}
	_ = MaskConfig(KindDingTalk, cfg)
	if IsMasked(cfg["secret"].(string)) {
		t.Fatal("MaskConfig가 입력 파라미터를 수정하여 실제 자격 증명이 마스킹된 값으로 덮어써지게 됩니다")
	}
}

func TestMergeConfigKeepsStoredOnMaskedIncoming(t *testing.T) {
	stored := map[string]any{"webhook": "https://real/hook", "secret": "REALSECRET", "method": "POST"}
	// 사용자는 method만 바꿨으며, 브라우저는 마스킹된 값+새 method를 제출합니다.
	incoming := map[string]any{
		"webhook": MaskedValue("https://real/hook"),
		"secret":  MaskedValue("REALSECRET"),
		"method":  "PUT",
	}
	got := MergeConfig(stored, incoming)
	if got["webhook"] != "https://real/hook" || got["secret"] != "REALSECRET" {
		t.Fatalf("마스킹된 필드는 DB의 기존 값을 유지해야 합니다. 실제: %v", got)
	}
	if got["method"] != "PUT" {
		t.Fatalf("수정한 필드가 적용되어야 합니다. 실제: %v", got["method"])
	}
}

func TestMergeConfigEmptyStringClears(t *testing.T) {
	stored := map[string]any{"webhook": "https://real/hook", "secret": "REALSECRET"}
	got := MergeConfig(stored, map[string]any{"secret": ""})
	if _, ok := got["secret"]; ok {
		t.Fatalf("빈 문자열은 해당 필드를 비워야 합니다. 실제: %v", got)
	}
	// 언급하지 않은 필드는 유지합니다(부분 업데이트 의미).
	if got["webhook"] != "https://real/hook" {
		t.Fatalf("언급하지 않은 필드는 유지해야 합니다. 실제: %v", got)
	}
}

func TestMergeConfigKeepsUnmentionedStoredKeys(t *testing.T) {
	stored := map[string]any{"host": "smtp.example.com", "port": float64(587), "password": "pw"}
	got := MergeConfig(stored, map[string]any{"port": float64(465)})
	if got["host"] != "smtp.example.com" || got["password"] != "pw" {
		t.Fatalf("언급하지 않은 필드는 유지해야 합니다. 실제: %v", got)
	}
	if got["port"] != float64(465) {
		t.Fatalf("언급한 필드는 업데이트해야 합니다. 실제: %v", got["port"])
	}
}

// TestPrepareConfigUpdateBlocksDestinationSwap은 이 패키지의 가장 중요한 보안 불변 조건입니다.
// **대상 주소를 바꿀 때 기존 자격 증명을 보내서는 안 됩니다**.
//
// 테스트는 바로 공격 형태의 입력(주소만 바꾸고 자격 증명은 언급하지 않음)을 사용하며,
// "방어 로직의 올바른 입력"을 사용하지 않습니다. 후자만 검사하면 방어가 적용되지 않아도 모두 통과합니다.
func TestPrepareConfigUpdateBlocksDestinationSwap(t *testing.T) {
	cases := []struct {
		name     string
		kind     string
		stored   map[string]any
		incoming map[string]any
		// wantMissing은 명시되어야 하는 자격 증명 키입니다.
		wantMissing string
	}{
		{
			name: "通用 Webhook 改地址想沿用 Authorization 头",
			kind: KindWebhook,
			stored: map[string]any{
				"url":     "https://legit.example.com/hook",
				"headers": map[string]any{"Authorization": "Bearer REAL-TOKEN"},
			},
			incoming:    map[string]any{"url": "https://attacker.tld/c"},
			wantMissing: "headers",
		},
		{
			name:        "Telegram 改 base_url 想把 Bot Token 发到自己的端点",
			kind:        KindTelegram,
			stored:      map[string]any{"bot_token": "123456:REAL", "chat_id": "1", "base_url": "https://api.telegram.org"},
			incoming:    map[string]any{"base_url": "https://attacker.tld"},
			wantMissing: "bot_token",
		},
		{
			name:        "邮件改 SMTP 主机想交出密码",
			kind:        KindEmail,
			stored:      map[string]any{"host": "smtp.corp.com", "port": 587, "password": "REALPW", "from": "a@b.c", "to": []any{"d@e.f"}},
			incoming:    map[string]any{"host": "smtp.attacker.tld"},
			wantMissing: "password",
		},
		{
			name:        "邮件关掉 TLS 也必须重新表态密码",
			kind:        KindEmail,
			stored:      map[string]any{"host": "smtp.corp.com", "port": 587, "tls": false, "password": "REALPW", "from": "a@b.c", "to": []any{"d@e.f"}},
			incoming:    map[string]any{"tls": true},
			wantMissing: "password",
		},
		{
			// 마스킹된 값 = "기존 자격 증명 사용"이며, 주소 변경 문맥에서는 마찬가지로 반드시 거부해야 합니다.
			name:        "回传掩码凭据 + 新地址",
			kind:        KindTelegram,
			stored:      map[string]any{"bot_token": "123456:REAL", "chat_id": "1", "base_url": "https://api.telegram.org"},
			incoming:    map[string]any{"base_url": "https://attacker.tld", "bot_token": MaskedValue("123456:REAL")},
			wantMissing: "bot_token",
		},
		{
			name:        "钉钉改 Webhook 想沿用加签密钥",
			kind:        KindDingTalk,
			stored:      map[string]any{"webhook": "https://oapi.dingtalk.com/robot/send?access_token=OLD", "secret": "REALSEC"},
			incoming:    map[string]any{"webhook": "https://attacker.tld/hook"},
			wantMissing: "secret",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			merged, err := PrepareConfigUpdate(tc.kind, tc.stored, tc.incoming)
			if err == nil {
				t.Fatalf("주소를 바꾸면서 자격 증명을 다시 명시하지 않았으므로 거부해야 합니다. 실제 설정: %v", merged)
			}
			var target *ErrDestinationChangedWithoutCredentials
			if !errors.As(err, &target) {
				t.Fatalf("API가 동작 가능한 안내를 제공하도록 전용 오류 타입을 반환해야 합니다. 실제: %T: %v", err, err)
			}
			found := false
			for _, m := range target.Missing {
				if m == tc.wantMissing {
					found = true
				}
			}
			if !found {
				t.Fatalf("누락된 자격 증명 키 %q 명시 필요. 실제: %v", tc.wantMissing, target.Missing)
			}
			// 오류 메시지는 동작을 수행하는 사람이 수정하는 방법을 알 수 있도록 해야 합니다.
			if !strings.Contains(err.Error(), tc.wantMissing) {
				t.Errorf("오류 메시지에 %q 포함 필요: %v", tc.wantMissing, err)
			}
		})
	}
}

// TestPrepareConfigUpdateAllowsLegitimateEdits는 반대 테스트입니다. 정상 편집을 잘못 차단해서는 안 되며,
// 그렇지 않으면 "너무 번거롭다"는 이유로 보호를 우회하거나 삭제할 수 있습니다.
func TestPrepareConfigUpdateAllowsLegitimateEdits(t *testing.T) {
	cases := []struct {
		name     string
		kind     string
		stored   map[string]any
		incoming map[string]any
	}{
		{
			name:     "只改名字（配置原样回传）",
			kind:     KindWebhook,
			stored:   map[string]any{"url": "https://legit.example.com/hook", "headers": map[string]any{"Authorization": "Bearer REAL"}},
			incoming: map[string]any{"url": MaskedValue("https://legit.example.com/hook")},
		},
		{
			name:     "只改请求方法，地址与凭据都不动",
			kind:     KindWebhook,
			stored:   map[string]any{"url": "https://legit.example.com/hook", "method": "POST"},
			incoming: map[string]any{"method": "PUT"},
		},
		{
			name:     "换地址并**同时**给新凭据",
			kind:     KindWebhook,
			stored:   map[string]any{"url": "https://old.example.com/hook", "headers": map[string]any{"Authorization": "Bearer OLD"}},
			incoming: map[string]any{"url": "https://new.example.com/hook", "headers": map[string]any{"Authorization": "Bearer NEW"}},
		},
		{
			name:     "换地址并显式声明不再需要凭据",
			kind:     KindWebhook,
			stored:   map[string]any{"url": "https://old.example.com/hook", "headers": map[string]any{"Authorization": "Bearer OLD"}},
			incoming: map[string]any{"url": "https://new.example.com/hook", "headers": ""},
		},
		{
			name:     "Telegram 改 chat_id（不是目的地）",
			kind:     KindTelegram,
			stored:   map[string]any{"bot_token": "t", "chat_id": "1", "base_url": "https://api.telegram.org"},
			incoming: map[string]any{"chat_id": "-100200"},
		},
		{
			name:     "邮件改收件人（不是目的地）",
			kind:     KindEmail,
			stored:   map[string]any{"host": "smtp.corp.com", "port": 587, "password": "PW", "from": "a@b.c", "to": []any{"x@y.z"}},
			incoming: map[string]any{"to": []any{"new@y.z"}},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			merged, err := PrepareConfigUpdate(tc.kind, tc.stored, tc.incoming)
			if err != nil {
				t.Fatalf("정상 편집을 잘못 차단했습니다: %v", err)
			}
			if merged == nil {
				t.Fatal("병합 결과를 반환해야 합니다")
			}
		})
	}
}

// TestPrepareConfigUpdatePortTypeTolerance는 잘못 판정하기 쉬운 세부 사항을 검사합니다.
// 프런트가 제출하는 포트는 JSON number(float64)이며, DB에서 읽는 값도 float64이지만,
// 두 값의 타입이 다를 수 있습니다(int vs float64 등). ==로 비교하면 "변경 없음"을 "변경됨"으로 판정하여,
// 이름만 바꾼 사용자에게 "비밀번호를 다시 입력하세요"를 표시합니다. 잘못된 경보는 보호를 신뢰하지 못하게 합니다.
func TestPrepareConfigUpdatePortTypeTolerance(t *testing.T) {
	stored := map[string]any{"host": "smtp.corp.com", "port": float64(587), "password": "PW"}
	// 같은 포트를 int 형태로 제출합니다.
	if _, err := PrepareConfigUpdate(KindEmail, stored, map[string]any{"port": 587}); err != nil {
		t.Fatalf("포트 값이 같고 타입만 다른 경우 주소 변경으로 판정해서는 안 됩니다: %v", err)
	}
	// 실제 포트 변경은 반드시 차단해야 합니다.
	if _, err := PrepareConfigUpdate(KindEmail, stored, map[string]any{"port": 25}); err == nil {
		t.Fatal("포트 변경은 차단해야 합니다")
	}
}

// TestPrepareConfigUpdateSurvivesRepeatedSaveWithBlankDestination은 "비워 둘 수 있는
// 대상 필드" 경로를 검사합니다. Telegram의 base_url을 비워 두면 공식 API 주소를 사용합니다.
//
// 과거에는 두 번째 저장부터 계속 채널 저장에 실패하는 경우가 있었습니다.
//
//	생성 시 DB에 base_url:"" 저장(생성 경로는 프런트의 config를 바로 저장하며 MergeConfig를 거치지 않음)
//	→ 첫 저장에서 MergeConfig가 빈 문자열을 명시적인 비우기로 보고 키 delete
//	→ 두 번째 저장에서 incoming은 여전히 ""이지만 stored에는 키가 없어 "주소 변경"으로 판정
//	→ bot_token은 마스킹된 반환값 → 400 "대상 주소가 변경되었습니다. 자격 증명 필드도 다시 입력하세요"
//
// 사용자는 아무것도 바꾸지 않았지만, Bot Token을 다시 붙여넣지 않으면 더 이상 저장할 수 없습니다.
func TestPrepareConfigUpdateSurvivesRepeatedSaveWithBlankDestination(t *testing.T) {
	stored := map[string]any{"bot_token": "123:ABC", "chat_id": "-100", "base_url": ""}

	// 프런트의 buildConfig()는 채널의 각 필드 정의에 값을 제출합니다. 자격 증명은 마스킹된 값으로 채우고,
	// 빈 텍스트 입력란은 빈 문자열을 제출합니다. 여기서는 "변경한 키"만 제출하지 않고 출력 전체를 재현합니다.
	submit := func() map[string]any {
		return map[string]any{
			"bot_token": MaskedValue("123:ABC"),
			"chat_id":   "-100",
			"base_url":  "",
		}
	}

	// 첫 번째 저장: 채널 이름만 바꾸고, config를 그대로 반환합니다.
	merged, err := PrepareConfigUpdate(KindTelegram, stored, submit())
	if err != nil {
		t.Fatalf("첫 번째 저장을 잘못 차단했습니다: %v", err)
	}
	if _, ok := merged["base_url"]; ok {
		t.Fatal("전제가 바뀌었습니다. 빈 문자열은 MergeConfig가 삭제해야 하며, 이 테스트는 바로 '키가 사라진 후' 단계를 검사합니다")
	}

	// 두 번째 저장: 제출 내용이 이전과 완전히 같고, 사용자는 아무것도 바꾸지 않았습니다.
	merged2, err := PrepareConfigUpdate(KindTelegram, merged, submit())
	if err != nil {
		t.Fatalf("두 번째 저장을 잘못 차단했습니다(사용자는 아무것도 바꾸지 않음): %v", err)
	}
	// 세 번째 저장: "한 번만 잘못됨"이 아니라 계속 저장할 수 있는지 확인합니다.
	if _, err := PrepareConfigUpdate(KindTelegram, merged2, submit()); err != nil {
		t.Fatalf("세 번째 저장을 잘못 차단했습니다: %v", err)
	}
	// 자격 증명은 빈 문자열 로직으로 함께 지워지지 않고 끝까지 유지되어야 합니다.
	if got := merged2["bot_token"]; got != "123:ABC" {
		t.Fatalf("Bot Token은 기존 값을 사용해야 합니다. 실제: %v", got)
	}
}

// TestPrepareConfigUpdateStillGuardsBlankDestinationChanges는 이전 테스트와 짝을 이루는
// 검사입니다. 빈 문자열과 "키 없음"을 같다고 보더라도 실제 주소 변경을 함께 허용해서는 **안 됩니다**.
// 두 방향 모두 실제 자격 증명 외부 전송 경로입니다. Telegram의 Bot Token은 URL 경로에 포함되며,
// base_url을 바꾸면 Token을 새 주소에 보내는 셈입니다.
func TestPrepareConfigUpdateStillGuardsBlankDestinationChanges(t *testing.T) {
	// 첫 번째 방향: "빈 값"(공식 주소)에서 자체 주소로 변경합니다.
	official := map[string]any{"bot_token": "123:ABC", "chat_id": "-100"}
	if _, err := PrepareConfigUpdate(KindTelegram, official, map[string]any{
		"bot_token": MaskedValue("123:ABC"),
		"base_url":  "https://tg-proxy.attacker.tld",
	}); err == nil {
		t.Fatal("공식 주소에서 자체 주소로 바꾸면 반드시 Token을 다시 입력하도록 요구해야 합니다")
	}

	// 두 번째 방향: 자체 주소 비우기(= 공식 API로 복귀)도 주소 변경입니다.
	proxied := map[string]any{"bot_token": "123:ABC", "base_url": "https://proxy.internal/bot"}
	if _, err := PrepareConfigUpdate(KindTelegram, proxied, map[string]any{
		"bot_token": MaskedValue("123:ABC"),
		"base_url":  "",
	}); err == nil {
		t.Fatal("자체 주소 비우기(공식 API로 돌아감)도 주소 변경이므로 반드시 Token을 다시 입력하도록 요구해야 합니다")
	}
}

func TestDestinationKeysDeclaredForEveryKind(t *testing.T) {
	// SecretKeys와 마찬가지로, 채널에서 대상 키 선언을 빠뜨리면 PrepareConfigUpdate로 보호할 수 없습니다.
	for kind, ch := range registry {
		if len(ch.DestinationKeys()) == 0 {
			t.Errorf("채널 %s: 대상 키를 선언하지 않아 주소 변경으로 인한 자격 증명 노출 방지가 적용되지 않습니다", kind)
		}
		if len(ch.SecretKeys()) == 0 {
			t.Errorf("채널 %s: 자격 증명 키를 선언하지 않았습니다", kind)
		}
	}
}

func TestSecretKeysDeclaredForEveryKind(t *testing.T) {
	// 컴파일러는 이미 각 채널에 SecretKeys 구현을 요구합니다. 여기서는 "마스킹 처리를 하지 않는
	// 채널이 없음"을 다시 확인합니다. 빈 슬라이스를 반환하는 채널의 자격 증명은 브라우저에 평문으로 반환됩니다.
	expect := map[string]bool{
		KindDingTalk: true, KindFeishu: true, KindWeCom: true,
		KindWebhook: true, KindTelegram: true, KindEmail: true,
	}
	for kind, ch := range registry {
		if !expect[kind] {
			t.Errorf("채널 %s: 테스트에 마스킹 기대 결과를 등록하지 않았습니다", kind)
			continue
		}
		if len(ch.SecretKeys()) == 0 {
			t.Errorf("채널 %s: 자격 증명 필드를 선언하지 않아 설정이 평문으로 반환됩니다", kind)
		}
	}
}

// TestPrepareConfigUpdateRejectsMaskedInContainer는 감사에서 지적된 빈틈을 검사합니다.
// 마스킹 센티널 값을 **문자열이 아닌** 구조(webhook.headers 객체 등)에 넣으면,
// MergeConfig는 '문자열이며 접두사가 있음'만 마스킹으로 인식하므로 리터럴 "__masked__"가
// 실제 헤더 값으로 DB에 저장됩니다. 이후 인증·인가 확인이 오류를 알리지 않고 작동하지 않으며 아무 오류도 발생하지 않습니다.
func TestPrepareConfigUpdateRejectsMaskedInContainer(t *testing.T) {
	stored := map[string]any{
		"url":     "https://legit.example.com/hook",
		"headers": map[string]any{"Authorization": "Bearer REAL"},
	}
	// 객체 내부에 마스킹 센티널 값을 넣습니다.
	incoming := map[string]any{
		"headers": map[string]any{"Authorization": MaskedPrefix},
	}
	if _, err := PrepareConfigUpdate(KindWebhook, stored, incoming); err == nil {
		t.Fatal("구조체 내부에 마스킹 센티널 값을 넣으면 거부해야 합니다(그렇지 않으면 리터럴을 DB에 저장함)")
	}
	// 객체 전체 제출(실제 새 값)은 정상적으로 허용합니다.
	ok := map[string]any{"headers": map[string]any{"Authorization": "Bearer NEW"}}
	if _, err := PrepareConfigUpdate(KindWebhook, stored, ok); err != nil {
		t.Fatalf("새 요청 헤더의 정상 제출을 차단해서는 안 됩니다: %v", err)
	}
}
