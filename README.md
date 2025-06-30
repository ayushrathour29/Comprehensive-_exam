# Go Job Queue System

A high-performance asynchronous job queue system built with Go, Gin, PostgreSQL, and Logrus. Supports job submission, retrieval, listing with pagination, and background processing with a worker pool. Designed for easy deployment on Render or Docker.


# deployed link: https://comprehensive-exam.onrender.com/jobs


## Features
- RESTful API (Gin)
- PostgreSQL persistence
- Structured logging (Logrus)
- Asynchronous job processing with a worker pool (5+ workers)
- Job status tracking (pending, processing, done, failed)
- Pagination for job listing
- Docker and Render deployment ready

## API Endpoints

### Submit a Job
```
POST /jobs
Content-Type: application/json
{
  "payload": "Process this data"
}
```
**Response:**
```
{
  "id": "...",
  "payload": "...",
  "status": "pending",
  "result": "",
  "created_at": "...",
  "updated_at": "..."
}
```

### Get Job by ID
```
GET /jobs/{id}
```
**Response:**
```
{
  "id": "...",
  "payload": "...",
  "status": "done",
  "result": "Job processed successfully",
  ...
}
```

### List Jobs (with Pagination)
```
GET /jobs?limit=10&offset=0
```
**Response:**
```
[
  { ... },
  { ... }
]
```

## Environment Variables
- `DATABASE_URL` (required): PostgreSQL connection string
- `PORT` (optional): Port to run the server (default: 8080)

## Local Development
1. **Start PostgreSQL** (Docker example):
   ```sh
   docker run --name jobqueue-postgres -e POSTGRES_PASSWORD=admin123 -e POSTGRES_USER=admin -e POSTGRES_DB=myjobdb -p 5432:5432 -d postgres:15
   ```
2. **Create `.env` file** in project root:
   ```
   DATABASE_URL=postgres://admin:admin123@localhost:5432/myjobdb?sslmode=disable
   ```
3. **Run migrations** (in psql):
   ```sql
   CREATE EXTENSION IF NOT EXISTS "pgcrypto";
   CREATE TABLE IF NOT EXISTS jobs (
     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
     payload TEXT NOT NULL,
     status VARCHAR(20) NOT NULL,
     result TEXT,
     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
     updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
   );
   ```
4. **Run the app:**
   ```sh
   go run cmd/main.go
   ```

## Docker Compose
1. Build and run both app and DB:
   ```sh
   docker-compose up --build
   ```

## Deploying to Render
- Create a new **Web Service** on Render, connect your repo.
- Set `DATABASE_URL` in Render's environment settings (use Render's managed PostgreSQL or your own).
- Use the provided `Dockerfile` or let Render auto-detect Go.


## License
MIT
