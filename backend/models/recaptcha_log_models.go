package models

import (
	"context"

	"main.go/db"
)

func CreateRecaptchaLog(userID *int64, action string, success bool) error {
	query := `INSERT INTO recaptcha_logs (user_id, action, success) VALUES ($1, $2, $3)`

	ctx := context.Background()
	_, err := db.Pool.Exec(ctx, query, userID, action, success)
	return err
}
