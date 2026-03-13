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
	// Stores the original model
	orig := modelCreateTest
	// restore it after the test
	defer func() { modelCreateTest = orig }()

	// valid results. model returns testid with no error
	modelCreateTest = func(test models.TestRequest, userID int64) (int64, error) { return 99, nil }

	// Fakes a http response (stores resp body, statuscode & headers)
	w := httptest.NewRecorder()
	// test context since http request needs a gin context
	c, _ := gin.CreateTestContext(w)
	// set user_id as done by RequireAuth fn in the middleware
	c.Set("user_id", int64(1))
	body := `{"name":"my test","url":"https://example.com","method":"GET", "status_code":200}`
	c.Request, _ = http.NewRequest(http.MethodPost, "/tests", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	CreateTest(c)

	// 200 status code
	if w.Code != http.StatusOK {
		// logf followed by failnow (reports a failure & stops the execution)
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	// unmarshal the resp body & assign it to resp variable
	json.Unmarshal(w.Body.Bytes(), &resp)
	// Match the id with our return value of mock model
	if resp["id"] != float64(99) {
		t.Errorf("expected id 99, got %v", resp["id"])
	}
}

func TestCreateTest_Unauthorized(t *testing.T) {
	// fakes http response (store resp body, headers & statuscode)
	w := httptest.NewRecorder()
	// test context since http request needs a gin context
	c, _ := gin.CreateTestContext(w)
	// no user_id set in context hence unauthorized
	c.Request, _ = http.NewRequest(http.MethodPost, "/tests",
		strings.NewReader(`{"name":"t","url":"https://x.com","method":"GET","status_code":200}`))
	c.Request.Header.Set("Content-Type", "application/json")

	CreateTest(c)

	// 401 statuscode
	if w.Code != http.StatusUnauthorized {
		// logf followed by fail (reports the failure without stopping the execution)
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestCreateTest_InvalidBody(t *testing.T) {
	// fakes a http response (stores resp body, statuscode & headers)
	w := httptest.NewRecorder()
	// test context since http request needs a gin context
	c, _ := gin.CreateTestContext(w)
	// Set user_id as done by the RequireAuth fn in the middleware
	c.Set("user_id", int64(1))
	// Invalid body in the request. shouldbindjson will throw error
	c.Request, _ = http.NewRequest(http.MethodPost, "/tests",
		strings.NewReader(`not json`))
	c.Request.Header.Set("Content-Type", "application/json")

	CreateTest(c)

	// 400 statuscode
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateTest_ModelError(t *testing.T) {
	// stores the original model
	orig := modelCreateTest
	// restores the model after the test
	defer func() { modelCreateTest = orig }()

	// model returns 0 as test id & an error
	modelCreateTest = func(test models.TestRequest, userID int64) (int64, error) {
		return 0, errors.New("db error")
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", int64(1))
	body := `{"name":"my test","url":"https://example.com","method":"GET","status_code":200}`
	c.Request, _ = http.NewRequest(http.MethodPost, "/tests", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	CreateTest(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}
