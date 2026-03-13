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

func init() {
	gin.SetMode(gin.TestMode)
}

func TestRegisterUser_Success(t *testing.T) {
	// Save the original functions (modelRegisterUser & generateToken)
	orig1, orig2 := modelRegisterUser, generateToken
	// After the test run, restore them back
	defer func() { modelRegisterUser, generateToken = orig1, orig2 }()

	// Mock the modelRegisterUser fn that takes input & returns fake output
	modelRegisterUser = func(email, password string) (int64, error) { return 42, nil }
	// Mock generateToken fn: returns "test-token"
	generateToken = func(userID int64, email string) (string, error) { return "test-token", nil }

	// Create a recorder (it stores resp body, statuscode & headers). Fakes a http response
	w := httptest.NewRecorder()
	// Create a test context (http requests need gin context, CreateTestContext creates test context)
	c, _ := gin.CreateTestContext(w)
	// Create a http request with params (method, url, request body)
	c.Request, _ = http.NewRequest(http.MethodPost, "/register",
		strings.NewReader(`{"email":"test@example.com","password":"pass123"}`))
	// Set the headers
	c.Request.Header.Set("Content-Type", "application/json")

	// Call the handler function
	RegisterUser(c)

	// Check if the response statuscode we got (stored in w) is okay, otherwise throw error
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	// unmarshal the response body (converting json to go map)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	// Check if value doesn't match
	if resp["token"] != "test-token" {
		t.Errorf("expected token 'test-token', got %v", resp["token"])
	}
}

func TestRegisterUser_InvalidInput(t *testing.T) {
	/* Note that in this test we haven't modified the original modelregisteruser
	*  because this invalid input error is returned before entering to
	*  modelregisteruser
	 */
	// Recorder (stores responses, statuscode & headers). Fakes a http response
	w := httptest.NewRecorder()
	// Test context since http request need a gin context
	c, _ := gin.CreateTestContext(w)
	// missing email field, password too short
	c.Request, _ = http.NewRequest(http.MethodPost, "/register",
		strings.NewReader(`{"password":"123"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	RegisterUser(c)

	if w.Code != http.StatusBadRequest {
		// Logf followed by fail (used to report a failure without stopping the execution)
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRegisterUser_ModelError(t *testing.T) {
	// Store original modelregisteruser
	orig := modelRegisterUser
	// restore it after the test finishes
	defer func() { modelRegisterUser = orig }()

	// create mock model fn which return 0 as userid & an error (model error)
	modelRegisterUser = func(email, password string) (int64, error) {
		// New formats the error based on the given test
		return 0, errors.New("duplicate email")
	}

	// Fake http response (stores resp body, headers & statuscode)
	w := httptest.NewRecorder()
	// Test context since http request needs a gin context
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/register",
		strings.NewReader(`{"email":"test@example.com","password":"pass123"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	RegisterUser(c)

	// Internal server error since the model is giving the error
	if w.Code != http.StatusInternalServerError {
		// logf followed by fail (used to report a failure without stopping the execution)
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestRegisterUser_TokenError(t *testing.T) {
	// store original models
	orig1, orig2 := modelRegisterUser, generateToken
	// restore them after the test
	defer func() { modelRegisterUser, generateToken = orig1, orig2 }()

	// modelregisteruser works fine and gives required results (an userid with no error)
	modelRegisterUser = func(email, password string) (int64, error) { return 1, nil }
	// generatetoken returns empty string & an error
	generateToken = func(userID int64, email string) (string, error) {
		return "", errors.New("signing error")
	}

	// Fakes a http response (stores response, statuscode & headers)
	w := httptest.NewRecorder()
	// test context since http requests needs a gin context
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/register",
		strings.NewReader(`{"email":"test@example.com","password":"pass123"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	RegisterUser(c)

	// 500 error that is to be thrown incase of token error
	if w.Code != http.StatusInternalServerError {
		// logf followed by fail (reports a failure without stopping an execution)
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestLoginUser_Success(t *testing.T) {
	// stores the original models
	orig1, orig2 := modelAuthenticateUser, generateToken
	// restores the models after the test
	defer func() { modelAuthenticateUser, generateToken = orig1, orig2 }()

	// model returns a valid object of type User (containing userid & email) with n errors
	modelAuthenticateUser = func(email, password string) (*models.User, error) {
		return &models.User{ID: 1, Email: email}, nil
	}
	// model returns a token (string) with no errors
	generateToken = func(userID int64, email string) (string, error) { return "test-token", nil }

	// Fakes a http response (stores resp body, headers & statuscode)
	w := httptest.NewRecorder()
	// test context since http requests needs a gin context
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/login",
		strings.NewReader(`{"email":"test@example.com","password":"pass123"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	LoginUser(c)

	// Status code should 200 if not, report as failed
	if w.Code != http.StatusOK {
		// logf followed by failnow (reports a failure & stops further execution since we don't
		// need to test for token in this case. There must be something wrong with the mock model)
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["token"] != "test-token" {
		// logf followed by fail (reports a failure without stopping the execution)
		t.Errorf("expected token 'test-token', got %v", resp["token"])
	}
}

func TestLoginUser_InvalidInput(t *testing.T) {
	// Fakes a http response (stores resp body, headers & statuscode)
	w := httptest.NewRecorder()
	// test context since http request needs a gin context
	c, _ := gin.CreateTestContext(w)
	// missing password field
	c.Request, _ = http.NewRequest(http.MethodPost, "/login",
		strings.NewReader(`{"email":"test@example.com"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	LoginUser(c)

	// status code 400
	if w.Code != http.StatusBadRequest {
		// logf followed by fail (reports failure without stopping execution)
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestLoginUser_AuthError(t *testing.T) {
	// Store the original model
	orig := modelAuthenticateUser
	// restore the model after the test
	defer func() { modelAuthenticateUser = orig }()

	// mock model returns no object & an error
	modelAuthenticateUser = func(email, password string) (*models.User, error) {
		return nil, errors.New("invalid credentials")
	}

	// Fakes a http response (store resp body, headers & statuscode)
	w := httptest.NewRecorder()
	// test context since http request need a gin context
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodPost, "/login",
		strings.NewReader(`{"email":"test@example.com","password":"wrongpass"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	LoginUser(c)

	// 401 statuscode
	if w.Code != http.StatusUnauthorized {
		// logf followed by fail (reports a failure without stopping execution)
		t.Errorf("expected 401, got %d", w.Code)
	}
}
