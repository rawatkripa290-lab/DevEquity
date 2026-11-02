# AI service starter

This service currently exposes a health endpoint and a safe `unknown` classification
placeholder. It deliberately does not call an LLM until you configure a provider,
structured-output schema, timeout, retry policy, prompt versioning and evaluation set.

Run locally:
`pip install -r requirements.txt`
`uvicorn app.main:app --reload --port 8000`

Endpoints:
- `GET /health`
- `POST /internal/classify`

Do not expose `/internal/classify` publicly. Add service-to-service authentication
before connecting it to the Go worker.
