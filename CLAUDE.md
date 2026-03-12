# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**test-runner** is a full-stack API load testing tool. Users define HTTP tests, then run them with configurable concurrency to measure performance and correctness. Results include pass/fail counts, avg/min/max durations per job.

## Commands

### Backend (Go + Gin)
```bash
cd backend
go run main.go          # Start server on :8080
go build ./...          # Compile check
go vet ./...            # Static analysis
```

### Frontend (React + Vite + TypeScript)
```bash
cd frontend
npm run dev             # Dev server on :5173 (proxies API calls to :8080)
npm run build           # tsc + vite build
npm run lint            # ESLint (zero warnings enforced)
```

## Architecture

### Backend (`backend/`)

**Request flow:** `main.go` → `routes/routes.go` (Gin router + JWT middleware) → `handlers/` → `models/` (DB queries) or `services/services.go` (business logic)

**Core execution engine** is `services/services.go`:
- `StartTestRun()` — creates a test run record, spawns a background goroutine
- `runJobs()` — worker pool with `sync.WaitGroup`; executes N concurrent HTTP requests
- `executeJob()` — makes the HTTP request, validates response body (JSON deep-equal or trimmed string), checks optional status code
- `UpdateTestRun()` — aggregates job results into pass/fail/duration stats on the test run record

**Key data model relationships:**
- `tests` → `test_runs` (one-to-many; test stores `latest_run_id` as denormalized cache)
- `test_runs` → `job_results` (one-to-many; one per concurrent worker)
- Test validation: optional `expected_response` (body match) and optional `status_code` fields

**Deployment:** `vercel.json` routes all traffic to `api/index.go` as a serverless Go function. Background goroutines may not complete reliably in serverless environments (hence the logging commits).

### Frontend (`frontend/src/`)

**Component tree:** `App.tsx` (theme + auth layout) → `TestForm.tsx` (create/edit) + `TestList.tsx` → `TestItem.tsx` → `RunTestModal.tsx` / `TestResultModal.tsx`

**Auth:** JWT stored in localStorage via `context/AuthContext.tsx`. Token sent as `Authorization: Bearer <token>` header.

**API proxy:** Vite dev server proxies `/tests`, `/register`, `/login` to `localhost:8080`. In production, set `VITE_API_BASE_URL`.

### API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/register` | Create account → `{token}` |
| POST | `/login` | Authenticate → `{token}` |
| GET | `/tests` | List all tests with latest run status |
| POST | `/tests` | Create test |
| PUT | `/tests/:id/edit` | Update test |
| DELETE | `/tests/:id/delete` | Delete test |
| POST | `/tests/:id/run` | Start test run `{concurrency}` → `{testRunId}` |
| GET | `/tests/:id` | Poll test run result by run ID |

### Environment Variables (`backend/.env`)

```
DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME
DB_SSLMODE, DB_CHANNEL_BINDING   # Neon-specific SSL settings
APP_ENV                          # local | production
MAX_ALLOWED_CONCURRENCY          # default 10
JWT_SECRET
```
