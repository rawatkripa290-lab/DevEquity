package processor

import (
	"path/filepath"
	"strings"
)

var ignoredExact = map[string]bool{
	"package-lock.json": true, "yarn.lock": true, "pnpm-lock.yaml": true,
	"go.sum": true, "cargo.lock": true, "poetry.lock": true,
	"vendor": true, "node_modules": true, "dist": true, "build": true,
}

// IsNoiseFile excludes generated, vendored, minified and lock files from code-change signals.
func IsNoiseFile(path string) bool {
	p := strings.ToLower(filepath.ToSlash(path))
	base := filepath.Base(p)
	if ignoredExact[base] {
		return true
	}
	if strings.Contains(p, "/vendor/") || strings.Contains(p, "/node_modules/") || strings.Contains(p, "/dist/") || strings.Contains(p, "/build/") {
		return true
	}
	if strings.HasSuffix(p, ".min.js") || strings.HasSuffix(p, ".min.css") || strings.HasSuffix(p, ".map") || strings.HasSuffix(p, ".lock") {
		return true
	}
	if strings.HasPrefix(base, "generated_") || strings.HasSuffix(base, ".generated.go") {
		return true
	}
	return false
}
