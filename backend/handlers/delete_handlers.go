package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Delete(c *gin.Context) {
	testIdStr := c.Param("id")

	testID, err := strconv.ParseInt(testIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid test id"})
		return
	}

	// Check if testID exists in db
	exists, err := modelCheckTestExists(testID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Test ID not found"})
		return
	}

	if err := modelDeleteTest(testID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "test deleted"})
}
