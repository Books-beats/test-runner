package handlers

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"main.go/models"
)

func ValidateTestRunRequest(testID int64, concurrency int, c *gin.Context) (bool, error) {
	// Check if testID exists in db
	exists, err := modelCheckTestExists(testID)
	if err != nil {
		log.Println("Failed to check test id: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return false, err
	}

	if !exists {
		log.Println("Test id doesn't exists")
		c.JSON(http.StatusNotFound, gin.H{"error": "Test ID not found"})
		return false, nil
	}

	// Check if concurrency is within allowed limits
	maxAllowedConcurrency, err := strconv.Atoi(os.Getenv("MAX_ALLOWED_CONCURRENCY"))
	if err != nil || maxAllowedConcurrency <= 0 {
		maxAllowedConcurrency = 10
	}

	if concurrency > maxAllowedConcurrency {
		log.Println("Concurrency exceeds")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Concurrency exceeds the maximum allowed limit of " + strconv.Itoa(maxAllowedConcurrency)})
		return false, nil
	}

	return true, nil
}

func CreateTestRun(c *gin.Context) {
	// Handler logic for creating a test run
	// user clicks on run button, backend request comes, we read test id from url params (c.Param("id"))
	// concurrency is sent in req body, we read it using c.ShouldBindJSON
	// Based on the concurrency, we create that many goroutines to run the test concurrently
	// while test is running, show status as pending, once test is done, update status to completed
	// After all the goroutines are done, we update the test run status to completed
	// We collect the results from all goroutines, aggregate them
	testIdStr := c.Param("id")

	testId, err := strconv.ParseInt(testIdStr, 10, 64)
	if err != nil {
		log.Println("Invalid test id: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid test id"})
		return
	}

	var request models.TestRunRequest

	e := c.ShouldBindJSON(&request)
	if e != nil {
		log.Println("Failed to bind", e)
		c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
		return
	}

	isValid, _ := ValidateTestRunRequest(testId, request.Concurrency, c)

	if !isValid {
		return
	}

	testRunId, status, err := serviceStartTestRun(testId, request.Concurrency)

	if err != nil {
		log.Println("Failed to start test run: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error: "})
		return
	}

	c.JSON(http.StatusOK, gin.H{"testRunId": testRunId, "status": status})
}
