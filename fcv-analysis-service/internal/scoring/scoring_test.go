package scoring

import "testing"

func TestScoreIsBounded(t *testing.T) {
	got := Calculate(Components{Code: 100, Collaboration: 100, Breadth: 100, Consistency: 100}, 1)
	if got.Score != 100 || got.Confidence != 1 {
		t.Fatalf("unexpected result: %+v", got)
	}
}

func TestNormalizeIsBounded(t *testing.T) {
	if got := Normalize(10, 10); got != 100 {
		t.Fatalf("got %v", got)
	}
	if got := Normalize(0, 10); got != 0 {
		t.Fatalf("got %v", got)
	}
}
