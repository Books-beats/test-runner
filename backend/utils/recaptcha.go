package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type RecaptchaResult struct {
	Success bool `json:"success"`
}

func VerifyRecaptcha(token string) (*RecaptchaResult, error) {
	secret := os.Getenv("RECAPTCHA_SECRET_KEY")
	if secret == "" {
		return nil, fmt.Errorf("RECAPTCHA_SECRET_KEY not set")
	}

	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", url.Values{
		"secret":   {secret},
		"response": {token},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to verify recaptcha: %w", err)
	}
	defer resp.Body.Close()

	var result RecaptchaResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse recaptcha response: %w", err)
	}

	if !result.Success {
		return &result, fmt.Errorf("recaptcha verification failed")
	}

	return &result, nil
}
