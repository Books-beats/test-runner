package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"main.go/models"
)

func TestCreateTest_Success(t *testing.T) {
	orig := modelCreateTest
	defer func() { modelCreateTest = orig }()

	modelCreateTest = func(test models.Test, userID int64) (int64, error) { return 99, nil }

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", int64(1))
	body := `{"name":"my test","url":"https://example.com","method":"GET"}`
	c.Request, _ = http.NewRequest(http.MethodPost, "/tests", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	CreateTest(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["id"] != float64(99) {
		t.Errorf("expected id 99, got %v", resp["id"])
	}
}

func TestCreateTest_Unauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// no user_id set in context
	c.Request, _ = http.NewRequest(http.MethodPost, "/tests",
		strings.NewReader(`{"name":"t","url":"https://x.com","method":"GET"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	CreateTest(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestCreateTest_InvalidBody(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", int64(1))
	c.Request, _ = http.NewRequest(http.MethodPost, "/tests",
		strings.NewReader(`not json`))
	c.Request.Header.Set("Content-Type", "application/json")

	CreateTest(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateTest_ModelError(t *testing.T) {
	orig := modelCreateTest
	defer func() { modelCreateTest = orig }()

	modelCreateTest = func(test models.Test, userID int64) (int64, error) {
		return 0, errors.New("db error")
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", int64(1))
	body := `{"name":"my test","url":"https://example.com","method":"GET"}`
	c.Request, _ = http.NewRequest(http.MethodPost, "/tests", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	CreateTest(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}
