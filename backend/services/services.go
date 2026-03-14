package services

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"main.go/models"
)

func UpdateTestRun(testRunId int64) {
	results, err := models.GetJobResultsByTestRunID(testRunId)
	if err != nil {
		log.Println("Failed to get job results by test run id: ", err)
		return
	}
	total := len(results)
	var passed, failed, avgDurationMs, maxDurationMs int
	minDurationMs := math.MaxInt

	for _, r := range results {
		if r.Passed {
			passed++
		} else {
			failed++
		}

		if r.DurationMs != nil {
			duration := *r.DurationMs
			avgDurationMs += duration

			if duration < minDurationMs {
				minDurationMs = duration
			}

			if duration > maxDurationMs {
				maxDurationMs = duration
			}
		}
	}

	if minDurationMs == math.MaxInt {
		minDurationMs = 0
	}

	if total > 0 {
		avgDurationMs = avgDurationMs / total
	}

	models.UpdateTestRun(testRunId, "completed", total, passed, failed, avgDurationMs, minDurationMs, maxDurationMs)
}

func responsesMatch(respBody []byte, expected string) bool {
	var respObj interface{}
	var expectedObj interface{}

	respErr := json.Unmarshal(respBody, &respObj)
	expErr := json.Unmarshal([]byte(expected), &expectedObj)

	// If both are valid JSON → compare JSON
	if respErr == nil && expErr == nil {
		return reflect.DeepEqual(respObj, expectedObj)
	}

	// fallback: compare trimmed strings
	return strings.TrimSpace(string(respBody)) == strings.TrimSpace(expected)
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
		log.Println("Failed to get test by id: ", e)
		return
	}

	// Start timer to measure duration of the job execution.
	start := time.Now()
	jobresult.CreatedAt = start

	// Create request
	// converting test.Body (string) -> bytes.NewBuffer([]byte(test.Body))
	// to create a buffer of bytes (io.Reader) that can be sent in the HTTP request body.
	url := strings.TrimSpace(test.URL)
	request, e1 := http.NewRequest(test.Method, url, bytes.NewBuffer([]byte(test.Body)))

	if e1 != nil {
		errStr := e1.Error()
		jobresult.Error = &errStr
		resultsChan <- jobresult
		log.Println("Failed to create a request: ", e1)
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
		resultsChan <- jobresult
		log.Println("Failed to send the request: ", e2)
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
		log.Println("Failed to read the res body: ", e3)
		return
	}

	responsezeStr := int(len(respBody))
	jobresult.ResponseSize = &responsezeStr

	statusMatch := true

	if test.StatusCode != nil {
		statusMatch = *test.StatusCode == response.StatusCode
	}

	// Compare response with expected response & status code match
	if responsesMatch(respBody, test.ExpectedResponse) && statusMatch {
		jobresult.Passed = true
		jobresult.Status = "completed"
	} else {
		jobresult.Passed = false
		jobresult.Status = "completed"
	}

	// Send job result to resultsChan
	resultsChan <- jobresult
}

func RunJobs(testID int64, concurrency int, testRunID int64) {
	for i := 0; i < concurrency; i++ {
		_, e3 := models.CreateJob(testRunID, i)
		if e3 != nil {
			log.Println("Failed to create job: ", e3)
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
		models.UpdateJobResult(result)
	}

	// After all jobs are done, update the test run fields
	UpdateTestRun(testRunID)

}

func StartTestRun(testID int64, concurrency int) (int64, string, error) {
	// Create a test run entry in the database with status "pending"
	testRunID, status, e2 := models.CreateTestRun(testID, concurrency)
	if e2 != nil {
		log.Println("Failed to create test run: ", e2)
		return 0, "stopped", e2
	}
	if os.Getenv("APP_ENV") == "production" {
		log.Println("Production: publishing to QStash")
		payload := ExecutePayload{TestID: testID, TestRunID: testRunID, Concurrency: concurrency}
		if err := publishToQStash(payload); err != nil {
			log.Println("Failed to publish to QStash: ", err)
			return 0, "stopped", err
		}
	} else {
		log.Println("Local: starting background RunJobs goroutine")
		go RunJobs(testID, concurrency, testRunID)
	}
	log.Println("Returning response")
	return testRunID, status, nil
}
