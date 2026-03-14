package middlewares

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	qstash "github.com/upstash/qstash-go"
)

// VerifyQStash validates the HMAC signature QStash sends on every webhook call.
// Stores the raw body bytes in context so the handler can decode them
// (body reader is consumed here and cannot be read again downstream).
func VerifyQStash() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentKey := os.Getenv("QSTASH_CURRENT_SIGNING_KEY")
		nextKey := os.Getenv("QSTASH_NEXT_SIGNING_KEY")

		if currentKey == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "QStash signing key not configured"})
			c.Abort()
			return
		}

		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
			c.Abort()
			return
		}

		receiver := qstash.NewReceiver(currentKey, nextKey)
		err = receiver.Verify(qstash.VerifyOptions{
			Signature: c.GetHeader("Upstash-Signature"),
			Body:      string(body),
			Url:       os.Getenv("QSTASH_TARGET_URL") + "/internal/execute",
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid QStash signature"})
			c.Abort()
			return
		}

		// Store body bytes in context — handler reads from here since reader is consumed.
		c.Set("rawBody", body)
		c.Next()
	}
}