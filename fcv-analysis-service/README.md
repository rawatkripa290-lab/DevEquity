# Fair Contribution View — backend

Go + Gin API, PostgreSQL, Redis/Asynq ingestion worker, and a Python FastAPI classification service.

## Quick start (PowerShell)

1. Copy `.env.example` to `.env` and set `GITHUB_TOKEN` (recommended for higher GitHub API limits).
2. Keep the `DATABASE_URL` hostname as `postgres` and `REDIS_ADDR=redis:6379` when using Compose.
3. Build and start services:

```powershell
Copy-Item .env.example .env
docker compose up --build -d
docker compose ps
```

For a brand-new database, PostgreSQL runs SQL files in `migrations/` in filename order on first initialization. **For an existing `pgdata` volume**, init scripts do not run again; apply the new migration explicitly:

```powershell
Get-Content .\migrations\003_metrics_and_identity.sql | docker compose exec -T postgres psql -U fairview -d fairview
```

4. Check services and submit a repository:

```powershell
Invoke-RestMethod http://localhost:8080/health
$body = @{ repository_url = 'https://github.com/golang/example' } | ConvertTo-Json
$result = Invoke-RestMethod -Method Post -Uri http://localhost:8080/api/v1/analyses -ContentType 'application/json' -Body $body
$result
Invoke-RestMethod "http://localhost:8080/api/v1/analyses/$($result.analysis_id)"
Invoke-RestMethod "http://localhost:8080/api/v1/analyses/$($result.analysis_id)/report"
```

AI service health: `http://localhost:8000/health`. Classification endpoint: `POST http://localhost:8000/internal/classify`.

## Local checks

```powershell
go test ./...
```

```powershell
python -m pip install -r ai-service/requirements.txt
python -m pytest ai-service/app/test_classifier.py
```

## Scoring and limitations

The current `fcv-1.0` score is an **activity proxy**, not a validated measurement of effort, skill, quality, or human worth. It uses observable commit/PR/review/issue activity, with logarithmic normalization and a confidence indicator. Private activity, deleted events, squash merges, shared work, and context not exposed by GitHub can affect results. The Python classifier is deterministic keyword/path heuristics; it is not a trained ML model or an LLM. The current Go worker does not yet send commit batches to the Python service, so its classification output is not yet included in the report score.

Do not expose GitHub tokens or `.env` files. Rotate any token that was ever committed or shared publicly.
