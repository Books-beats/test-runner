package utils

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Returns a byte slice of the secret key used for signing JWTs
var GetSecretKey = func() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	return []byte(secret)
}

func GenerateToken(userID int64, email string) (string, error) {
	// Create token with claims & signing method
	// NewWithClaims takes a signing method and a set of claims (payload), returns a token object
	// Claims are the data stored in the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Expires in 24 hours
	})

	// Jwt token contains 3 parts: header, payload, signature
	// Here, we encode the header and payload,
	// then sign it with our secret key to create the signature
	// Eg: HMACSHA256(
	//    base64(header) + "." + base64(payload),
	//    secret_key
	// )
	tokenString, err := token.SignedString(GetSecretKey())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Verify signature method
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}
		return GetSecretKey(), nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	// Payload is converted to a map[string]interface{} by jwt.MapClaims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	// Extract user_id. jwt.MapClaims parses numbers as float64
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("user_id not found in token")
	}

	return int64(userIDFloat), nil
}
