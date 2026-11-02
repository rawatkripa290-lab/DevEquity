from .schemas import CommitInput, Classification
from .prompts import CLASSIFICATION_POLICY


def classify(commit: CommitInput) -> Classification:
    text = (commit.message + " " + commit.diff_summary).lower()
    paths = [p.lower() for p in commit.files]
    category, confidence, reason = "unknown", 0.0, "Insufficient evidence for a reliable category."
    rules = [
        (("fix", "bug", "patch", "hotfix"), "bugfix", 0.78, "Commit wording indicates a bug fix."),
        (("feat", "feature", "implement", "add "), "feature", 0.75, "Commit wording indicates new functionality."),
        (("refactor", "cleanup", "restructure"), "refactor", 0.76, "Commit wording indicates restructuring."),
        (("test", "spec"), "test", 0.78, "Commit wording indicates test work."),
        (("docs", "documentation", "readme"), "docs", 0.78, "Commit wording indicates documentation."),
        (("chore", "dependency", "bump "), "chore", 0.70, "Commit wording indicates maintenance."),
    ]
    for keywords, label, conf, explanation in rules:
        if any(k in text for k in keywords):
            category, confidence, reason = label, conf, explanation
            break
    if category == "unknown" and paths:
        if all(p.endswith(".md") or "/docs/" in p for p in paths):
            category, confidence, reason = "docs", 0.65, "All changed paths appear documentation-related."
        elif all("test" in p or "spec" in p for p in paths):
            category, confidence, reason = "test", 0.65, "All changed paths appear test-related."
    return Classification(sha=commit.sha, category=category, confidence=confidence,
                          explanation=f"{reason} {CLASSIFICATION_POLICY}")
