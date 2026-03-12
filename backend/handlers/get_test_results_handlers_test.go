package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"main.go/models"
)

func TestGetTestResult_Success(t *testing.T) {
	orig := modelGetTestResult
	defer func() { modelGetTestResult = orig }()

	total, passed, failed := 5, 4, 1
	modelGetTestResult = func(testRunID int64) (models.TestRun, error) {
		return models.TestRun{
			ID:          testRunID,
			TestID:      10,
			Concurrency: 5,
			Status:      "completed",
			Total:       &total,
			Passed:      &passed,
			Failed:      &failed,
		}, nil
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "42"}}
	c.Request, _ = http.NewRequest(http.MethodGet, "/tests/42", nil)

	GetTestResult(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	result, ok := resp["result"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected result object, got %v", resp["result"])
	}
	if result["status"] != "completed" {
		t.Errorf("expected status 'completed', got %v", result["status"])
	}
}

func TestGetTestResult_InvalidID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "not-a-number"}}
	c.Request, _ = http.NewRequest(http.MethodGet, "/tests/not-a-number", nil)

	GetTestResult(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetTestResult_ModelError(t *testing.T) {
	orig := modelGetTestResult
	defer func() { modelGetTestResult = orig }()

	modelGetTestResult = func(testRunID int64) (models.TestRun, error) {
		return models.TestRun{}, errors.New("not found")
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "99"}}
	c.Request, _ = http.NewRequest(http.MethodGet, "/tests/99", nil)

	GetTestResult(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}
