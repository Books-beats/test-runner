package handlers

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"main.go/models"
	"main.go/services"
)

func ValidateTestRunRequest(testID int64, concurrency int, c *gin.Context) (bool, error) {
	// Check if testID exists in db
	exists, err := models.CheckTestIdExists(testID)
	if err != nil {
		log.Println("er", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return false, err
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Test ID not found"})
		return false, nil
	}

	e1 := godotenv.Load()
	if e1 != nil {
		log.Println("exits", exists, "e1", e1)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load environment variables"})
		return false, e1
	}

	// Check if concurrency is within allowed limits
	maxallowedconcurrency := os.Getenv("MAX_ALLOWED_CONCURRENCY")
	maxallowedconcurrencyInt, _ := strconv.Atoi(maxallowedconcurrency)
	log.Println("maxon", maxallowedconcurrency)
	if concurrency > maxallowedconcurrencyInt {
		c.JSON(http.StatusBadRequest, gin.H{"error": `Concurrency exceeds the maximum allowed limit of` + maxallowedconcurrency})
		return false, nil
	}
	log.Println("tst", testID, concurrency)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid test id"})
		return
	}

	var request models.TestRunRequest

	e := c.ShouldBindJSON(&request)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
		return
	}

	isValid, _ := ValidateTestRunRequest(testId, request.Concurrency, c)
	log.Println("isvaild", isValid, testIdStr, testId)
	if !isValid {
		return
	}

	testRunId, status, err := services.StartTestRun(testId, request.Concurrency)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error: "})
		return
	}

	c.JSON(http.StatusOK, gin.H{"testRunId": testRunId, "status": status})
}
