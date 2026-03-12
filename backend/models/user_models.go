package models

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"main.go/db"
)

type User struct {
	ID           int64  `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

func RegisterUser(email, password string) (int64, error) {
	// Hash password
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return 0, err
	}

	query := `INSERT INTO users (email, password_hash, created_at) VALUES ($1, $2, NOW()) RETURNING id`

	ctx := context.Background()
	var userID int64
	err = db.Pool.QueryRow(ctx, query, email, string(hashedBytes)).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func AuthenticateUser(email, password string) (*User, error) {
	query := `SELECT id, email, password_hash FROM users WHERE email = $1`

	ctx := context.Background()
	var user User
	err := db.Pool.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Compare hashed passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}
