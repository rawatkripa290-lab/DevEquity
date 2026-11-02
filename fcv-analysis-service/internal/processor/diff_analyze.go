package processor

import "strings"

type ChangeCategory string

const (
	Feature  ChangeCategory = "feature"
	Bugfix   ChangeCategory = "bugfix"
	Refactor ChangeCategory = "refactor"
	Test     ChangeCategory = "test"
	Docs     ChangeCategory = "docs"
	Chore    ChangeCategory = "chore"
	Unknown  ChangeCategory = "unknown"
)

// ClassifyCommit is a deterministic baseline classifier. It is not an ML verdict.
func ClassifyCommit(message string, files []string) ChangeCategory {
	m := strings.ToLower(strings.TrimSpace(message))
	if strings.HasPrefix(m, "test") || strings.Contains(m, "add test") {
		return Test
	}
	if strings.HasPrefix(m, "docs") || strings.Contains(m, "documentation") {
		return Docs
	}
	if strings.HasPrefix(m, "fix") || strings.Contains(m, "bug") || strings.Contains(m, "patch") {
		return Bugfix
	}
	if strings.HasPrefix(m, "feat") || strings.Contains(m, "implement") || strings.Contains(m, "add ") {
		return Feature
	}
	if strings.HasPrefix(m, "refactor") || strings.Contains(m, "cleanup") {
		return Refactor
	}
	if strings.HasPrefix(m, "chore") || strings.Contains(m, "bump ") || strings.Contains(m, "update deps") {
		return Chore
	}
	if len(files) > 0 {
		allDocs, allTests, allNoise := true, true, true
		for _, f := range files {
			p := strings.ToLower(f)
			allDocs = allDocs && (strings.HasSuffix(p, ".md") || strings.Contains(p, "docs/"))
			allTests = allTests && (strings.Contains(p, "test") || strings.Contains(p, "spec"))
			allNoise = allNoise && IsNoiseFile(f)
		}
		if allDocs {
			return Docs
		}
		if allTests {
			return Test
		}
		if allNoise {
			return Chore
		}
	}
	return Unknown
}
