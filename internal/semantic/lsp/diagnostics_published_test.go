package lsp

import "testing"

// get_diagnostics reports the second return as `published`: an empty report
// from the server and no report at all must stay distinguishable.
func TestLastDiagnosticsTellsCleanFromNoAnswer(t *testing.T) {
	const path = "/w/a.ts"
	p := &Provider{lastDiag: map[string][]Diagnostic{}}

	if _, ok := p.LastDiagnostics(path); ok {
		t.Fatal("no publish yet must read as not published")
	}
	p.lastDiag[path] = []Diagnostic{}
	got, ok := p.LastDiagnostics(path)
	if !ok || len(got) != 0 {
		t.Fatalf("an empty publish must read as published and clean, got %v %v", got, ok)
	}
}
