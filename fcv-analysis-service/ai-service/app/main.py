from fastapi import FastAPI
from .schemas import Classification, ClassificationRequest
from .classifier import classify

app = FastAPI(title="Fair Contribution View AI Service", version="1.0.0")

@app.get("/health")
async def health():
    return {"status": "healthy", "service": "ai-service"}

@app.post("/internal/classify", response_model=list[Classification])
async def classify_commits(request: ClassificationRequest):
    return [classify(commit) for commit in request.commits]
