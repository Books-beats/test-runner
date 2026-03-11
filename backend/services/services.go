package services

import (
	"bytes"
	"io"
	"net/http"
	"sync"
	"time"

	"main.go/models"
)

func UpdateTestRun(testRunId int64) {
	results, err := models.GetJobResultsByTestRunID(testRunId)
	if err != nil {
		return
	}
	total := len(results)
	var passed, failed, avgDurationMs, minDurationMs, maxDurationMs int

	for _, r := range results {
		if r.Passed {
			passed++
		} else {
			failed++
		}

		if r.DurationMs != nil {
			duration := *r.DurationMs
			avgDurationMs += duration

			if minDurationMs == 0 || duration < minDurationMs {
				minDurationMs = duration
			}

			if duration > maxDurationMs {
				maxDurationMs = duration
			}
		}
	}

	if total > 0 {
		avgDurationMs = avgDurationMs / total
	}

	models.UpdateTestRun(testRunId, "completed", total, passed, failed, avgDurationMs, minDurationMs, maxDurationMs)
}

func executeJob(testID int64, testRunId int64, jobID int, resultsChan chan<- models.JobResult) {
	var jobresult models.JobResult

	jobresult.TestRunID = testRunId
	jobresult.JobNumber = jobID

	test, e := models.GetTestByID(testID)
	if e != nil {
		// Assigning error string to a variable bcoz e1.Error() returns a value that is temporary
		// Go doesn't allow taking address of temporary value,
		// so we need to assign it to a variable before taking its address
		errStr := e.Error()
		jobresult.Error = &errStr
		resultsChan <- jobresult
		return
	}

	// Start timer to measure duration of the job execution.
	start := time.Now()
	jobresult.CreatedAt = start

	// Create request
	// converting test.Body (string) -> bytes.NewBuffer([]byte(test.Body))
	// to create a buffer of bytes (io.Reader) that can be sent in the HTTP request body.
	request, e1 := http.NewRequest(test.Method, test.URL, bytes.NewBuffer([]byte(test.Body)))
	if e1 != nil {
		errStr := e1.Error()
		jobresult.Error = &errStr
		resultsChan <- jobresult
		return
	}

	// Set headers
	for key, value := range test.Headers {
		request.Header.Set(key, value)
	}

	// Send request
	client := &http.Client{Timeout: 10 * time.Second}
	response, e2 := client.Do(request)
	if e2 != nil {
		errStr := e2.Error()
		jobresult.Error = &errStr
		return
	}
	defer response.Body.Close() // must close the response body after reading to free up resources

	endedAt := time.Now()
	duration := time.Since(start)

	completedAtStr := endedAt
	jobresult.CompletedAt = &completedAtStr

	durationMsStr := int(duration.Milliseconds())
	jobresult.DurationMs = &durationMsStr

	statusCodeStr := int(response.StatusCode)
	jobresult.StatusCode = &statusCodeStr

	// Read response body
	respBody, e3 := io.ReadAll(response.Body)
	if e3 != nil {
		errStr := e3.Error()
		jobresult.Error = &errStr
		resultsChan <- jobresult
		return
	}

	responsezeStr := int(len(respBody))
	jobresult.ResponseSize = &responsezeStr

	// Compare response with expected response
	if string(respBody) == test.ExpectedResponse {
		jobresult.Passed = true
		jobresult.Status = "completed"
	} else {
		jobresult.Passed = false
		jobresult.Status = "completed"
	}

	// Send job result to resultsChan
	resultsChan <- jobresult
}

func runJobs(testID int64, concurrency int, testRunID int64) {
	for i := 0; i < concurrency; i++ {
		_, e3 := models.CreateJob(testRunID, i)
		if e3 != nil {
			return
		}
	}

	jobChan := make(chan int, concurrency)
	resultsChan := make(chan models.JobResult, concurrency)

	var wg sync.WaitGroup

	// Creating worker pool to execute jobs concurrently
	for w := 0; w < concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for jobID := range jobChan {
				executeJob(testID, testRunID, jobID, resultsChan)
			}
		}()
	}

	// Send job IDs to jobChan for workers to execute
	for i := 0; i < concurrency; i++ {
		jobChan <- i
	}
	close(jobChan) // Close the job channel after sending all job IDs

	wg.Wait()
	close(resultsChan) // Close the results channel after all workers are done

	// Collect results from resultsChan and store them in the database
	for result := range resultsChan {
		models.CreateJobResult(result)
	}

	// After all jobs are done, update the test run fields
	UpdateTestRun(testRunID)

}

func StartTestRun(testID int64, concurrency int) (int64, string, error) {
	// Create a test run entry in the database with status "pending"
	testRunID, status, e2 := models.CreateTestRun(testID, concurrency)
	if e2 != nil {
		return 0, "stopped", e2
	}
	go runJobs(testID, concurrency, testRunID)
	return testRunID, status, nil
}
