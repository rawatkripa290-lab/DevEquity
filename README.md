# ⚖️ DevEquity — Fair Contribution Analysis Platform

<p align="center">
  <strong>Understand contributions. Recognize collaboration. Build fairer teams.</strong>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Project-DevEquity-blueviolet?style=for-the-badge" />
  <img src="https://img.shields.io/badge/Status-MVP-orange?style=for-the-badge" />
  <img src="https://img.shields.io/badge/Backend-Go-00ADD8?style=for-the-badge&logo=go" />
  <img src="https://img.shields.io/badge/AI-Python-yellow?style=for-the-badge&logo=python" />
</p>

<p align="center">
  <a href="#-overview">Overview</a> •
  <a href="#-features">Features</a> •
  <a href="#-architecture">Architecture</a> •
  <a href="#-installation">Installation</a> •
  <a href="#-api-documentation">API</a> •
  <a href="#-contributing">Contributing</a>
</p>

---

## 📌 Overview

**DevEquity** is a GitHub repository contribution analysis platform designed to help development teams understand how work is distributed across contributors.

Simply submit a GitHub repository URL and DevEquity collects repository activity, processes commit information, analyzes contributor metrics, and generates a structured contribution report.

The goal is to make contribution patterns more transparent and encourage fairer recognition of collaborative software development.

> Our philosophy: Contribution is more than the number of commits someone makes.

## ✨ Features

| Feature | Description |
|---|---|
| 🔗 Repository Analysis | Analyze a GitHub repository using its URL |
| 📊 Contributor Metrics | Calculate contribution-related metrics |
| 🤖 AI Classification | Classify commits using a Python AI service |
| ⚙️ Background Processing | Process repository data asynchronously |
| 📈 Progress Tracking | Monitor analysis stages |
| 📄 Report Generation | Generate a structured analysis report |
| 🐳 Docker Support | Run backend services using Docker Compose |

### 🔄 How DevEquity works

1. Enter a GitHub repository URL.
2. The Go API validates the repository.
3. An analysis task is added to Redis using Asynq.
4. The Go worker fetches repository activity.
5. The Python AI service classifies commit messages.
6. Go processing code calculates contributor metrics.
7. A report is generated and stored in PostgreSQL.
8. The frontend retrieves and displays the report.

---

## 🏗️ Architecture

DevEquity uses a distributed service architecture combining Go, Python, React, PostgreSQL, and Redis.

```mermaid
flowchart TD
    A["React + Vite Frontend"] --> B["Go REST API - Gin"]
    B --> C[("PostgreSQL")]
    B --> D["Redis + Asynq Queue"]
    D --> E["Go Ingestion Worker"]
    E --> F["GitHub REST API"]
    E --> G["Python FastAPI Classifier"]
    E --> H["Go Metrics & Scoring"]
    H --> I["Report Generator"]
    I --> C
    C --> J["Frontend Report Display"]
```

### 🧰 Technology stack

| Layer | Technology |
|---|---|
| Frontend | React, Vite, JavaScript |
| Backend API | Go, Gin |
| Background Worker | Go, Asynq |
| AI Service | Python, FastAPI, Pydantic |
| Database | PostgreSQL |
| Message Queue | Redis |
| External API | GitHub REST API |
| Containerization | Docker, Docker Compose |

---

## 📁 Project Structure

```text
DevEquity/
│
├── frontend/
│   ├── src/
│   ├── public/
│   ├── package.json
│   └── vite.config.js
│
├── fcv-analysis-service/
│   │
│   ├── cmd/
│   │   ├── api/
│   │   └── worker/
│   │
│   ├── internal/
│   │   ├── ai/
│   │   ├── config/
│   │   ├── db/
│   │   ├── github/
│   │   ├── handlers/
│   │   ├── ingestion/
│   │   ├── processor/
│   │   ├── queue/
│   │   ├── reports/
│   │   ├── scoring/
│   │   ├── server/
│   │   └── worker/
│   │
│   ├── ai-service/
│   │   └── app/
│   │
│   ├── migrations/
│   ├── Dockerfile
│   ├── docker-compose.yml
│   └── .env.example
│
└── README.md
```

---

## 🚀 Installation

### Prerequisites

Make sure you have installed:

- [Git](https://git-scm.com/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [Node.js](https://nodejs.org/)
- npm
- GitHub personal access token (recommended)

Check your installations:

```bash
git --version
docker --version
docker compose version
node --version
npm --version
```

### 1️⃣ Clone the repository

```bash
git clone https://github.com/rawatkripa290-lab/DevEquity.git

cd DevEquity
```

### 2️⃣ Configure the backend

```powershell
cd fcv-analysis-service

Copy-Item .env.example .env

notepad .env
```

Configure your environment variables:

```env
APP_ENV=development
PORT=8080

DATABASE_URL=postgres://fairview:fairview_dev_password@postgres:5432/fairview?sslmode=disable

REDIS_ADDR=redis:6379
REDIS_PASSWORD=

GITHUB_TOKEN=YOUR_GITHUB_TOKEN

COMMIT_LIMIT=300

AI_SERVICE_URL=http://ai-service:8000
```

Replace `YOUR_GITHUB_TOKEN` with your own GitHub token.

**Never commit your `.env` file or expose credentials publicly.**

### 3️⃣ Start backend services

Run from `fcv-analysis-service`:

```powershell
docker compose config --quiet

docker compose up -d --build

docker compose ps
```

Check the logs:

```powershell
docker compose logs --tail 100 api

docker compose logs --tail 150 worker

docker compose logs --tail 100 ai-service
```

Check PostgreSQL:

```powershell
docker compose exec postgres pg_isready -U fairview -d fairview
```

### 4️⃣ Start the frontend

Open a new terminal:

```powershell
cd frontend

Copy-Item .env.example .env

npm ci

npm run dev
```

Open the URL printed by Vite, usually:

```text
http://localhost:5173
```

Frontend API configuration:

```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

### 5️⃣ Build the frontend

```powershell
npm run build
```

---

## 🔌 API Documentation

The frontend communicates with the Go backend through REST APIs.

| Method | Endpoint | Description |
|---|---|---|
| POST | `/api/v1/analyses` | Create repository analysis |
| GET | `/api/v1/analyses/{id}` | Get analysis progress |
| GET | `/api/v1/analyses/{id}/report` | Retrieve generated report |

### Example: Create an analysis

```powershell
$body = @{
    repository_url = "https://github.com/golang/example"
} | ConvertTo-Json

$response = Invoke-RestMethod `
    -Uri "http://localhost:8080/api/v1/analyses" `
    -Method POST `
    -ContentType "application/json" `
    -Body $body

$response | ConvertTo-Json -Depth 10
```

Copy the returned `analysis_id`.

### Example: Check analysis status

```powershell
$id = "YOUR_ANALYSIS_ID"

Invoke-RestMethod `
    -Uri "http://localhost:8080/api/v1/analyses/$id" |
    ConvertTo-Json -Depth 10
```

### Example: Retrieve report

Run after the analysis status becomes `completed`:

```powershell
Invoke-RestMethod `
    -Uri "http://localhost:8080/api/v1/analyses/$id/report" |
    ConvertTo-Json -Depth 20
```

---

## 🤖 AI Classification Service

The Python FastAPI service handles commit classification.

The Go client sends commit data to:

```text
POST /internal/classify
```

Current Python schema constraints:

| Field | Maximum |
|---|---:|
| SHA | 64 characters |
| Commit message | 2,000 characters |
| Files | 200 entries |
| Diff summary | 8,000 characters |
| Commits per request | 100 |

The Go client prepares commit data to satisfy these limits and batches requests containing more than 100 commits.

---

## 🧪 Testing

### Go backend

Run from `fcv-analysis-service`:

```bash
go test ./...
```

### Frontend

Run from `frontend`:

```bash
npm run build
```

Run Python tests or lint checks if configured in the project.

For an end-to-end check, submit a repository URL, wait for the analysis to complete, and confirm that the report endpoint returns the generated report.

---

## 🛠️ Troubleshooting

<details>
<summary><strong>❌ npm installation fails with EPERM</strong></summary>

Stop Vite and close processes using the frontend dependency folder.

```powershell
taskkill /F /IM node.exe 2>$null

Remove-Item node_modules -Recurse -Force -ErrorAction SilentlyContinue

npm cache verify

npm ci

npm run dev
```
</details>

<details>
<summary><strong>❌ Python AI service returns HTTP 422</strong></summary>

Check worker logs:

```powershell
docker compose logs --tail 200 worker
```

Possible causes include:

- Commit message exceeds 2,000 characters.
- More than 200 files are sent for one commit.
- Diff summary exceeds 8,000 characters.
- More than 100 commits are submitted in one request.

Check that the Go client applies the Python schema constraints.
</details>

<details>
<summary><strong>❌ Analysis remains stuck at fetching_commits</strong></summary>

Fetching all commits does not mean the entire analysis is complete.

Check:

```powershell
docker compose logs --tail 200 worker
```

Inspect the next processing stage, errors, and Asynq retry exhaustion.
</details>

<details>
<summary><strong>❌ Report not found</strong></summary>

The report becomes available only after report generation and database persistence succeed.

Check:

1. Analysis status.
2. Worker logs.
3. PostgreSQL connectivity.
4. Report generation errors.
</details>

<details>
<summary><strong>❌ PostgreSQL connection error</strong></summary>

```powershell
docker compose ps

docker compose exec postgres pg_isready -U fairview -d fairview
```

Inside Docker Compose, use the service hostname `postgres`, not `localhost`.
</details>

---

## 🔐 Security

- Never commit `.env` files.
- Never expose GitHub tokens or database passwords.
- Use least-privilege credentials.
- Keep PostgreSQL and Redis private in production.
- Use strong production secrets.
- Validate and rate-limit API requests.
- Avoid logging tokens or sensitive repository data.

---

## ⚖️ Responsible Use & Limitations

Contribution metrics are not perfect measures of individual performance.

Commit counts and commit classifications cannot fully capture:

- Code reviews
- Mentoring
- Architecture and design
- Planning and coordination
- Incident response
- Work outside the repository

DevEquity should be used as a **decision-support tool**, not as the sole basis for employee ranking, hiring, promotion, compensation, or disciplinary decisions.

Metrics should always be interpreted with context and reviewed alongside the underlying evidence.

---

## 🤝 Contributing

Contributions are welcome!

1. Fork the repository.
2. Create a feature branch.
3. Make your changes.
4. Run the relevant tests.
5. Submit a pull request with a clear description.

Example:

```bash
git checkout -b feature/your-feature-name

git add .

git commit -m "feat: describe your change"

git push origin feature/your-feature-name
```

---

## 🌟 Support the Project

If you find DevEquity interesting:

- ⭐ Star the repository.
- 🍴 Fork it and experiment.
- 🐛 Report bugs through GitHub Issues.
- 💡 Suggest improvements.
- 🤝 Contribute code and documentation.

<p align="center">
  <strong>DevEquity</strong><br/>
  Making repository contribution patterns easier to understand.
</p>

<p align="center">
  Built with ❤️ using Go, Python, React, PostgreSQL and Redis.
</p>

---

<p align="center">
  <sub>Development / MVP — features and deployment readiness are subject to verification.</sub>
</p>
