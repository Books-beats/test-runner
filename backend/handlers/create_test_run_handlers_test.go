package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreateTestRun_Success(t *testing.T) {
	orig1, orig2 := modelCheckTestExists, serviceStartTestRun
	defer func() { modelCheckTestExists, serviceStartTestRun = orig1, orig2 }()

	modelCheckTestExists = func(testID int64) (bool, error) { return true, nil }
	serviceStartTestRun = func(testID int64, concurrency int) (int64, string, error) {
		return 55, "pending", nil
	}

	t.Setenv("MAX_ALLOWED_CONCURRENCY", "10")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request, _ = http.NewRequest(http.MethodPost, "/tests/1/run",
		strings.NewReader(`{"concurrency":3}`))
	c.Request.Header.Set("Content-Type", "application/json")

	CreateTestRun(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["testRunId"] != float64(55) {
		t.Errorf("expected testRunId 55, got %v", resp["testRunId"])
	}
	if resp["status"] != "pending" {
		t.Errorf("expected status 'pending', got %v", resp["status"])
	}
}

func TestCreateTestRun_InvalidID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "not-a-number"}}
	c.Request, _ = http.NewRequest(http.MethodPost, "/tests/not-a-number/run",
		strings.NewReader(`{"concurrency":1}`))
	c.Request.Header.Set("Content-Type", "application/json")

	CreateTestRun(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateTestRun_InvalidBody(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request, _ = http.NewRequest(http.MethodPost, "/tests/1/run",
		strings.NewReader(`not json`))
	c.Request.Header.Set("Content-Type", "application/json")

	CreateTestRun(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateTestRun_NotFound(t *testing.T) {
	orig := modelCheckTestExists
	defer func() { modelCheckTestExists = orig }()

	// return (false, nil) → 404
	modelCheckTestExists = func(testID int64) (bool, error) { return false, nil }

	t.Setenv("MAX_ALLOWED_CONCURRENCY", "10")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "999"}}
	c.Request, _ = http.NewRequest(http.MethodPost, "/tests/999/run",
		strings.NewReader(`{"concurrency":1}`))
	c.Request.Header.Set("Content-Type", "application/json")

	CreateTestRun(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestCreateTestRun_ConcurrencyExceeded(t *testing.T) {
	orig := modelCheckTestExists
	defer func() { modelCheckTestExists = orig }()

	modelCheckTestExists = func(testID int64) (bool, error) { return true, nil }

	t.Setenv("MAX_ALLOWED_CONCURRENCY", "5")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request, _ = http.NewRequest(http.MethodPost, "/tests/1/run",
		strings.NewReader(`{"concurrency":10}`))
	c.Request.Header.Set("Content-Type", "application/json")

	CreateTestRun(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateTestRun_ServiceError(t *testing.T) {
	orig1, orig2 := modelCheckTestExists, serviceStartTestRun
	defer func() { modelCheckTestExists, serviceStartTestRun = orig1, orig2 }()

	modelCheckTestExists = func(testID int64) (bool, error) { return true, nil }
	serviceStartTestRun = func(testID int64, concurrency int) (int64, string, error) {
		return 0, "stopped", errors.New("failed to insert")
	}

	t.Setenv("MAX_ALLOWED_CONCURRENCY", "10")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request, _ = http.NewRequest(http.MethodPost, "/tests/1/run",
		strings.NewReader(`{"concurrency":2}`))
	c.Request.Header.Set("Content-Type", "application/json")

	CreateTestRun(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}
