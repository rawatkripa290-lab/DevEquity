from typing import Literal
from pydantic import BaseModel, Field

CommitType = Literal["feature", "bugfix", "refactor", "docs", "test", "chore", "unknown"]

class CommitInput(BaseModel):
    sha: str = Field(min_length=1, max_length=64)
    message: str = Field(default="", max_length=2000)
    files: list[str] = Field(default_factory=list, max_length=200)
    diff_summary: str = Field(default="", max_length=8000)

class Classification(BaseModel):
    sha: str
    category: CommitType
    confidence: float = Field(ge=0, le=1)
    explanation: str

class ClassificationRequest(BaseModel):
    commits: list[CommitInput] = Field(max_length=100)
