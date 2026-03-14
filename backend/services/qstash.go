package services

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	qstash "github.com/upstash/qstash-go"
)

// ExecutePayload is the message body published to QStash and received
// by the /internal/execute webhook handler.
type ExecutePayload struct {
	TestID      int64 `json:"testId"`
	TestRunID   int64 `json:"testRunId"`
	Concurrency int   `json:"concurrency"`
}

func publishToQStash(payload ExecutePayload) error {
	token := os.Getenv("QSTASH_TOKEN")
	if token == "" {
		return fmt.Errorf("QSTASH_TOKEN is not set")
	}
	targetURL := os.Getenv("QSTASH_TARGET_URL")
	if targetURL == "" {
		return fmt.Errorf("QSTASH_TARGET_URL is not set")
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal QStash payload: %w", err)
	}

	client := qstash.NewClient(token)
	res, err := client.Publish(qstash.PublishOptions{
		Url:         targetURL + "/internal/execute",
		Body:        string(body),
		ContentType: "application/json",
	})
	if err != nil {
		return fmt.Errorf("QStash publish failed: %w", err)
	}

	log.Printf("QStash message published, messageId: %s", res.MessageId)
	return nil
}
