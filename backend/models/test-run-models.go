package models

import (
	"context"
	"log"
	"time"

	"main.go/db"
)

type TestRunRequest struct {
	Concurrency int `json:"concurrency" binding:"required"`
}

type TestRun struct {
	ID            int64      `json:"id"`
	TestID        int64      `json:"test_id"`
	Concurrency   int        `json:"concurrency"`
	Status        string     `json:"status"`
	Total         *int       `json:"total"`
	Passed        *int       `json:"passed"`
	Failed        *int       `json:"failed"`
	AvgDurationMs *int       `json:"avg_duration_ms"`
	MinDurationMs *int       `json:"min_duration_ms"`
	MaxDurationMs *int       `json:"max_duration_ms"`
	CreatedAt     time.Time  `json:"created_at"`
	CompletedAt   *time.Time `json:"completed_at"`
}

func CheckTestIdExists(testID int64) (bool, error) {
	query := `SELECT id FROM tests WHERE id = $1`

	ctx := context.Background()
	var id int64
	err := db.Pool.QueryRow(ctx, query, testID).Scan(&id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func CreateTestRun(testID int64, concurrency int) (int64, string, error) {
	query := `INSERT INTO test_runs (test_id, concurrency, status, created_at) VALUES ($1, $2, 'pending', NOW()) RETURNING id, status`

	ctx := context.Background()
	var testRunID int64
	var status string
	err := db.Pool.QueryRow(ctx, query, testID, concurrency).Scan(&testRunID, &status)

	if err != nil {
		log.Println("Model error:", err)
		return 0, "stopped", err
	}

	return testRunID, status, nil
}

func CreateJob(testRunId int64, jobNumber int) (int64, error) {
	query := `INSERT INTO job_results (test_run_id, job_number, status) VALUES ($1, $2, 'pending') RETURNING id`

	ctx := context.Background()
	var jobID int64
	err := db.Pool.QueryRow(ctx, query, testRunId, jobNumber).Scan(&jobID)

	if err != nil {
		return 0, err
	}

	return jobID, nil
}

func UpdateTestRun(testRunID int64, status string, total int, passed int, failed int, avgDurationMs int, minDurationMs int, maxDurationMs int) error {
	query := `UPDATE test_runs 
	SET status = $1, total = $2, passed = $3, failed = $4, avg_duration_ms = $5, min_duration_ms = $6, max_duration_ms = $7, completed_at = NOW()
	WHERE id = $8`

	ctx := context.Background()
	_, err := db.Pool.Exec(
		ctx,
		query,
		status,
		total,
		passed,
		failed,
		avgDurationMs,
		minDurationMs,
		maxDurationMs,
		testRunID,
	)

	return err
}

func GetTestRunResult(testRunID int64) (TestRun, error) {
	query := `SELECT id, test_id, concurrency, status, total, passed, failed, avg_duration_ms, min_duration_ms, max_duration_ms, created_at, completed_at FROM test_runs WHERE id = $1`

	ctx := context.Background()
	row := db.Pool.QueryRow(ctx, query, testRunID)

	var testRunResult TestRun

	// In case of null values what Scan does is
	// assigns nil to the pointer fields
	// For eg: *total = nil if total is null in db
	err := row.Scan(
		&testRunResult.ID,
		&testRunResult.TestID,
		&testRunResult.Concurrency,
		&testRunResult.Status,
		&testRunResult.Total,
		&testRunResult.Passed,
		&testRunResult.Failed,
		&testRunResult.AvgDurationMs,
		&testRunResult.MinDurationMs,
		&testRunResult.MaxDurationMs,
		&testRunResult.CreatedAt,
		&testRunResult.CompletedAt,
	)

	if err != nil {
		return TestRun{}, err
	}

	return testRunResult, nil
}
