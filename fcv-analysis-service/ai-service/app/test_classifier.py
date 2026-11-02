from .classifier import classify
from .schemas import CommitInput


def test_bugfix_classification():
    result = classify(CommitInput(sha="abc", message="fix: prevent nil pointer crash"))
    assert result.category == "bugfix"
    assert result.confidence > 0


def test_unknown_when_no_evidence():
    result = classify(CommitInput(sha="def", message="misc changes"))
    assert result.category == "unknown"
    assert result.confidence == 0
