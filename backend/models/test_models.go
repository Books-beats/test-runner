package models

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"main.go/db"
)

type TestRequest struct {
	Name             string            `json:"name" binding:"required"`
	URL              string            `json:"url" binding:"required"`
	Method           string            `json:"method" binding:"required"`
	Headers          map[string]string `json:"headers"`
	Body             string            `json:"body"`
	ExpectedResponse string            `json:"expected_response" binding:"required_without=StatusCode"`
	StatusCode       *int              `json:"status_code" binding:"required_without=ExpectedResponse"`
}

type Test struct {
	ID               int64             `json:"id"`
	UserID           int64             `json:"user_id"`
	Name             string            `json:"name"`
	URL              string            `json:"url"`
	Method           string            `json:"method"`
	Headers          map[string]string `json:"headers"`
	Body             string            `json:"body"`
	ExpectedResponse string            `json:"expected_response"`
	StatusCode       *int              `json:"status_code"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	LatestRunID      *int64            `json:"latest_run_id,omitempty"`
	LatestRunStatus  *string           `json:"latest_run_status,omitempty"`
}

func CreateTest(test TestRequest, userID int64) (int64, error) {
	// Marshaling the Headers map into a JSON string to store it in the database.
	// Converting Go data structrues to JSON or bytes
	// Headers field is converted to map[string]string in the Test struct,
	// but in the database, we want to store it as a JSON string.
	headersJSON, err := json.Marshal(test.Headers)
	if err != nil {
		return 0, err
	}

	query := `INSERT INTO tests (user_id, name, url, method, headers, body, expected_response, status_code) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	// Creating a context to manage the lifecycle of the database operation. This allows for better control over timeouts and cancellations.
	ctx := context.Background()
	var testID int64
	err = db.Pool.QueryRow(
		ctx,
		query,
		userID,
		test.Name,
		test.URL,
		test.Method,
		headersJSON, // Store as JSON string in PostgreSQL
		test.Body,
		test.ExpectedResponse,
		test.StatusCode,
	).Scan(&testID)

	if err != nil {
		return 0, err
	}

	return testID, nil
}

func GetTestByID(testID int64) (*Test, error) {
	var test Test
	var headersJson []byte

	query := `SELECT id, url, method, headers, body, expected_response, status_code FROM tests WHERE id = $1`

	ctx := context.Background()
	row := db.Pool.QueryRow(ctx, query, testID)

	err := row.Scan(&test.ID, &test.URL, &test.Method, &headersJson, &test.Body, &test.ExpectedResponse, &test.StatusCode)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(headersJson, &test.Headers) // convert JSON back to Go data type
	return &test, nil
}

func GetAllTests(userID int64) ([]Test, error) {
	// For each test, find the latest run id and status using LATERAL JOIN
	// We run the subquery for each row in the tests table to get the latest test run for that test,
	// and join it with the tests table to get the latest run id and status along with the test details.
	query := `
		SELECT t.id, t.user_id, t.name, t.url, t.method, t.headers, t.body, t.expected_response, t.status_code, t.created_at, t.updated_at,
               r.id as latest_run_id, r.status as latest_run_status
        FROM tests t
        LEFT JOIN LATERAL (
            SELECT id, status 
            FROM test_runs 
            WHERE test_id = t.id 
            ORDER BY created_at DESC 
            LIMIT 1
        ) r ON true
        WHERE t.user_id = $1 
        ORDER BY t.created_at DESC`

	ctx := context.Background()
	rows, err := db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tests []Test
	for rows.Next() {
		var test Test
		var headersJSON []byte

		err := rows.Scan(
			&test.ID,
			&test.UserID,
			&test.Name,
			&test.URL,
			&test.Method,
			&headersJSON,
			&test.Body,
			&test.ExpectedResponse,
			&test.StatusCode,
			&test.CreatedAt,
			&test.UpdatedAt,
			&test.LatestRunID,
			&test.LatestRunStatus,
		)
		if err != nil {
			return nil, err
		}

		json.Unmarshal(headersJSON, &test.Headers)
		tests = append(tests, test)
	}

	return tests, nil
}

func DeleteTest(testID int64) error {
	query := `DELETE FROM tests where id = $1`

	ctx := context.Background()
	_, err := db.Pool.Exec(ctx, query, testID)

	if err != nil {
		log.Println("Failed to delete row with id ", testID)
		return err
	}
	return nil
}

func UpdateTest(test Test, testId int64) error {
	headersJSON, err := json.Marshal(test.Headers)
	if err != nil {
		return err
	}

	query := `
		UPDATE tests
		SET name = $1,
			url = $2,
			method = $3,
			headers = $4,
			body = $5,
			expected_response = $6,
			status_code = $7,
			updated_at = NOW()
		WHERE id = $8
	`

	ctx := context.Background()

	_, err = db.Pool.Exec(
		ctx,
		query,
		test.Name,
		test.URL,
		test.Method,
		headersJSON,
		test.Body,
		test.ExpectedResponse,
		test.StatusCode,
		testId,
	)

	if err != nil {
		return err
	}

	return nil
}
