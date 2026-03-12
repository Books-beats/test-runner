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

func TestGetAllTests_Success(t *testing.T) {
	orig := modelGetAllTests
	defer func() { modelGetAllTests = orig }()

	modelGetAllTests = func(userID int64) ([]models.Test, error) {
		return []models.Test{
			{ID: 1, Name: "test-one", URL: "https://example.com", Method: "GET"},
			{ID: 2, Name: "test-two", URL: "https://example.org", Method: "POST"},
		}, nil
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", int64(1))
	c.Request, _ = http.NewRequest(http.MethodGet, "/tests", nil)

	GetAllTests(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	tests, ok := resp["tests"].([]interface{})
	if !ok || len(tests) != 2 {
		t.Errorf("expected 2 tests, got %v", resp["tests"])
	}
}

func TestGetAllTests_ReturnsEmptyArrayWhenNil(t *testing.T) {
	orig := modelGetAllTests
	defer func() { modelGetAllTests = orig }()

	modelGetAllTests = func(userID int64) ([]models.Test, error) { return nil, nil }

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", int64(1))
	c.Request, _ = http.NewRequest(http.MethodGet, "/tests", nil)

	GetAllTests(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	tests, ok := resp["tests"].([]interface{})
	if !ok {
		t.Fatalf("expected tests to be an array, got %v", resp["tests"])
	}
	if len(tests) != 0 {
		t.Errorf("expected empty array, got length %d", len(tests))
	}
}

func TestGetAllTests_Unauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// no user_id
	c.Request, _ = http.NewRequest(http.MethodGet, "/tests", nil)

	GetAllTests(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestGetAllTests_ModelError(t *testing.T) {
	orig := modelGetAllTests
	defer func() { modelGetAllTests = orig }()

	modelGetAllTests = func(userID int64) ([]models.Test, error) {
		return nil, errors.New("db error")
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", int64(1))
	c.Request, _ = http.NewRequest(http.MethodGet, "/tests", nil)

	GetAllTests(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}
