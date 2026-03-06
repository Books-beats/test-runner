# Test Runner

A fullstack application for managing and running API tests with concurrency support.

## Tech Stack

- **Frontend**: React + Vite + TypeScript
- **Backend**: Go + Gin + PostgreSQL

## Prerequisites

- [Node.js](https://nodejs.org/) & npm
- [Go](https://golang.org/)
- [PostgreSQL](https://www.postgresql.org/)

## Getting Started

### Backend

1. Navigate to the backend directory:

   ```bash
   cd backend
   ```

2. Make sure PostgreSQL is running and update your credentials in the .env file

3. Run the Go server:
   ```bash
   go run main.go
   ```
   > The backend exposes API routes under `http://localhost:8080/`.

### Frontend

1. Navigate to the frontend directory:

   ```bash
   cd frontend
   ```

2. Install frontend dependencies:

   ```bash
   npm install
   ```

3. Start the Vite development server:
   ```bash
   npm run dev
   ```
   > The frontend will be available at `http://localhost:5173/`.

## Features

- Create API tests
- Configure individual test requests (Method, URL, Expected Response)
- Run tests directly from the dashboard
- Set concurrency configuration per test run
- Poll test run results
- View comprehensive results displaying total requests, average response time, passed assertions, and failed requests.
