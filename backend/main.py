from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from routes.analysis import router as analysis_router
app = FastAPI(
  title="Fair Contribution View",
  description="Analyze fair contributions in GitHub repositories",
  version="1.0/0"
)
app.add_middleware(
    CORSMiddleware,
    allow_origins=[
    "http://localhost:5173",
    "http://127.0.0.1:5173"
],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)
@app.get("/")
def home():
  return {
    "message":"Fair Contribution View Backend is running"
  }
app.include_router(analysis_router,prefix="/api")