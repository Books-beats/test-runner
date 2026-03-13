package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"main.go/models"
)

func CreateTest(c *gin.Context) {
	// Extract user_id from Gin Context (set by RequireAuth middleware)
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user"})
		return
	}
	userID, ok := userIDVal.(int64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user"})
		return
	}

	var newTest models.TestRequest

	// Reads the JSON body, maps it to request Go variable, checks for errors & returns a JSON reposnse
	err := c.ShouldBindJSON(&newTest)

	if err != nil {
		// http.StatusBadRequest is a constant that represents HTTP status code 400.
		// gin.H type is a shortcut for map[string]interface{} & is used in Gin to create JSON responses.
		log.Println("Error in json binding: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	testID, err := modelCreateTest(newTest, userID)
	if err != nil {
		log.Println("Failed to create test: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": testID})
}
