package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"main.go/models"
)

func GetAllTests(c *gin.Context) {
	// Check if user exists
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user"})
		return
	}
	userID := userIDVal.(int64)

	tests, err := models.GetAllTests(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tests"})
		return
	}

	// Returning empty array instead of null if no tests exist
	if tests == nil {
		tests = []models.Test{}
	}

	c.JSON(http.StatusOK, gin.H{"tests": tests})
}
