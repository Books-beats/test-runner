package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetTestResult(c *gin.Context) {
	testIdStr := c.Param("id")

	// Convert testIdStr to int64
	testId, err := strconv.ParseInt(testIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid test id"})
		return
	}

	result, err := modelGetTestResult(testId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Test result not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}
