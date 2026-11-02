package processor

import "math"

type Metrics struct {
	Commits, PullRequests, Reviews, Issues             int
	Additions, Deletions                               int
	Features, Bugfixes, Refactors, Tests, Docs, Chores int
}

// ActivityScore is a bounded activity proxy with diminishing returns.
func ActivityScore(m Metrics) float64 {
	raw := float64(m.Commits)*1.0 + float64(m.PullRequests)*2.0 + float64(m.Reviews)*1.5 + float64(m.Issues)*0.5
	return round2(100 * (1 - math.Exp(-raw/20)))
}
func round2(v float64) float64 { return math.Round(v*100) / 100 }
