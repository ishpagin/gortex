package lsp

import (
	"encoding/json"
	"testing"
)

func TestEffectiveInitOptionsConsultsSpecFunc(t *testing.T) {
	spec := &ServerSpec{
		InitializationOptions: json.RawMessage(`{"static":true}`),
		InitOptionsFunc: func(root string) json.RawMessage {
			return json.RawMessage(`{"root":"` + root + `"}`)
		},
	}
	p := &Provider{spec: spec}
	if got := string(p.effectiveInitOptions("/w")); got != `{"root":"/w"}` {
		t.Fatalf("options = %s", got)
	}
	spec.InitOptionsFunc = func(string) json.RawMessage { return nil }
	if got := string(p.effectiveInitOptions("/w")); got != `{"static":true}` {
		t.Fatalf("empty func result must fall back, got %s", got)
	}
}
