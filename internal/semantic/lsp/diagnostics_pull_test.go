package lsp

import (
	"encoding/json"
	"testing"
)

func TestPullDiagnosticsAsksAServerThatAdvertisesIt(t *testing.T) {
	server := newFakeLSPServer()
	var asked string
	server.handle("textDocument/diagnostic", func(params json.RawMessage) (any, *jsonRPCError) {
		var req struct {
			TextDocument TextDocumentIdentifier `json:"textDocument"`
		}
		if err := json.Unmarshal(params, &req); err != nil {
			t.Error(err)
		}
		asked = req.TextDocument.URI
		return map[string]any{"kind": "full", "items": []Diagnostic{{Message: "TS2322"}}}, nil
	})
	p, cleanup := providerWithFakeServer(t, server, []string{"typescript"})
	defer cleanup()
	p.caps.DiagnosticProvider = map[string]any{"interFileDependencies": true}

	got, ok := p.PullDiagnostics("/w/a.ts")

	if !ok || len(got) != 1 || got[0].Message != "TS2322" {
		t.Fatalf("got %v %v", got, ok)
	}
	if asked != pathToURI("/w/a.ts") {
		t.Fatalf("asked about %q", asked)
	}
}

func TestPullDiagnosticsReportsACleanFileAsAnAnswer(t *testing.T) {
	server := newFakeLSPServer()
	server.handle("textDocument/diagnostic", func(json.RawMessage) (any, *jsonRPCError) {
		return map[string]any{"kind": "full", "items": []Diagnostic{}}, nil
	})
	p, cleanup := providerWithFakeServer(t, server, []string{"typescript"})
	defer cleanup()
	p.caps.DiagnosticProvider = true

	if got, ok := p.PullDiagnostics("/w/a.ts"); !ok || len(got) != 0 {
		t.Fatalf("got %v %v, want an empty answer", got, ok)
	}
}

func TestPullDiagnosticsLeavesPushServersAlone(t *testing.T) {
	server := newFakeLSPServer()
	server.handle("textDocument/diagnostic", func(json.RawMessage) (any, *jsonRPCError) {
		t.Error("a server that did not advertise pull must not be asked")
		return nil, nil
	})
	p, cleanup := providerWithFakeServer(t, server, []string{"typescript"})
	defer cleanup()

	if _, ok := p.PullDiagnostics("/w/a.ts"); ok {
		t.Fatal("no pull support must read as no answer")
	}
}

func TestPullDiagnosticsFallsBackWhenTheRequestFails(t *testing.T) {
	server := newFakeLSPServer()
	server.handle("textDocument/diagnostic", func(json.RawMessage) (any, *jsonRPCError) {
		return nil, &jsonRPCError{Code: -32801, Message: "content modified"}
	})
	p, cleanup := providerWithFakeServer(t, server, []string{"typescript"})
	defer cleanup()
	p.caps.DiagnosticProvider = true

	if _, ok := p.PullDiagnostics("/w/a.ts"); ok {
		t.Fatal("a failed pull must read as no answer so the caller uses the push path")
	}
}
