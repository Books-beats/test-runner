package models

import (
	"context"
	"time"

	"main.go/db"
)

type JobResult struct {
	ID           int64      `json:"id"`
	TestRunID    int64      `json:"test_run_id"`
	JobNumber    int        `json:"job_number"`
	Status       string     `json:"status"`
	Passed       bool       `json:"passed"`
	CreatedAt    time.Time  `json:"created_at"`
	StatusCode   *int       `json:"status_code"`   // Using pointer to handle null values
	DurationMs   *int       `json:"duration_ms"`   // If there is null in db, it will be nil in Go
	ResponseSize *int       `json:"response_size"` // That's we are giving them pointer type
	Error        *string    `json:"error"`
	CompletedAt  *time.Time `json:"completed_at"`
}

func CreateJobResult(jobResult JobResult) (int64, error) {
	query := `
	INSERT INTO job_results 
	(test_run_id, job_number, status, status_code, duration_ms, response_size, passed, error, created_at, completed_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id
	`

	ctx := context.Background()
	var jobResultID int64
	err := db.Pool.QueryRow(ctx,
		query,
		jobResult.TestRunID,
		jobResult.JobNumber,
		jobResult.Status,
		jobResult.StatusCode,
		jobResult.DurationMs,
		jobResult.ResponseSize,
		jobResult.Passed,
		jobResult.Error,
		jobResult.CreatedAt,
		jobResult.CompletedAt,
	).Scan(&jobResultID)

	if err != nil {
		return 0, err
	}

	return jobResultID, nil
}

func GetJobResultsByTestRunID(testRunID int64) ([]JobResult, error) {
	query := `SELECT id, test_run_id, job_number, status, status_code, duration_ms, response_size, passed, error, created_at, completed_at FROM job_results WHERE test_run_id = $1`

	ctx := context.Background()
	rows, err := db.Pool.Query(ctx, query, testRunID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobResults []JobResult
	// rows.Next() is used to iterate over the result set returned by the query.
	// It returns true if there is another row to read, and false when there are no more rows.
	for rows.Next() {
		var jobResult JobResult
		// rows.Scan reads the columns of current row that rows.Next is pointing to
		err := rows.Scan(
			&jobResult.ID,
			&jobResult.TestRunID,
			&jobResult.JobNumber,
			&jobResult.Status,
			&jobResult.StatusCode,
			&jobResult.DurationMs,
			&jobResult.ResponseSize,
			&jobResult.Passed,
			&jobResult.Error,
			&jobResult.CreatedAt,
			&jobResult.CompletedAt,
		)
		if err != nil {
			return nil, err
		}

		jobResults = append(jobResults, jobResult)
	}

	return jobResults, nil
}
