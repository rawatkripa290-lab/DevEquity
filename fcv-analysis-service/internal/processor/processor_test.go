package processor

import "testing"

func TestNoiseFilter(t *testing.T) {
	if !IsNoiseFile("frontend/node_modules/pkg/index.js") {
		t.Fatal("expected generated dependency file to be ignored")
	}
	if IsNoiseFile("src/payment_service.go") {
		t.Fatal("source file must not be ignored")
	}
}

func TestClassifyCommit(t *testing.T) {
	if got := ClassifyCommit("fix: handle nil pointer", nil); got != Bugfix {
		t.Fatalf("got %s", got)
	}
	if got := ClassifyCommit("update", []string{"docs/guide.md"}); got != Docs {
		t.Fatalf("got %s", got)
	}
}
