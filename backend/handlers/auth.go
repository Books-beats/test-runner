package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func RegisterUser(c *gin.Context) {
	var input RegisterInput
	// Validate user input (email & password)
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Invalid user input: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := modelRegisterUser(input.Email, input.Password)

	if err != nil {
		log.Println("Failed to register user: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	token, e := generateToken(userID, input.Email)

	if e != nil {
		log.Println("Failed to generate token: ", e)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Registration successful",
		"token":   token,
	})
}

func LoginUser(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Failed to login: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := modelAuthenticateUser(input.Email, input.Password)

	if err != nil {
		log.Println("Failed to authenticate: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, e := generateToken(user.ID, user.Email)
	if e != nil {
		log.Println("Failed to generate token: ", e)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}
