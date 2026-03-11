package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"main.go/models"
)

func Delete(c *gin.Context) {
	testIdStr := c.Param("id")

	testID, err := strconv.ParseInt(testIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid test id"})
		return
	}

	// Check if testID exists in db
	exists, err := models.CheckTestIdExists(testID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Test ID not found"})
		return
	}

	models.DeleteTest(testID)
}
