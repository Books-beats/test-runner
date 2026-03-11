package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"main.go/models"
)

func Update(c *gin.Context) {
	testIdStr := c.Param("id")

	testID, err := strconv.ParseInt(testIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid test id"})
		return
	}

	var test models.Test

	e := c.ShouldBindJSON(&test)

	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
		return
	}

	models.UpdateTest(test, testID)
}
