package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"main.go/services"
)

func ExecuteTestRun(c *gin.Context) {
	rawBody, exists := c.Get("rawBody")
	if !exists {
		log.Println("rawBody not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	var payload services.ExecutePayload
	if err := json.Unmarshal(rawBody.([]byte), &payload); err != nil {
		log.Println("Failed to parse QStash payload: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	if payload.TestID == 0 || payload.TestRunID == 0 || payload.Concurrency <= 0 {
		log.Println("Invalid QStash payload values")
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required fields"})
		return
	}

	log.Printf("QStash executing: testID=%d testRunID=%d concurrency=%d",
		payload.TestID, payload.TestRunID, payload.Concurrency)

	// RunJobs is called synchronously — this IS the background work.
	// QStash retries on non-2xx, so 200 is only sent after full completion.
	serviceRunJobs(payload.TestID, payload.Concurrency, payload.TestRunID)

	c.JSON(http.StatusOK, gin.H{"status": "completed"})
}
