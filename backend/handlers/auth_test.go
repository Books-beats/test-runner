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
	orig1, orig2 := modelRegisterUser, generateToken
	defer func() { modelRegisterUser, generateToken = orig1, orig2 }()

	modelRegisterUser = func(email, password string) (int64, error) { return 42, nil }
	generateToken = func(userID int64, email string) (string, error) { return "test-token", nil }

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/register",
		strings.NewReader(`{"email":"test@example.com","password":"pass123"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	RegisterUser(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["token"] != "test-token" {
		t.Errorf("expected token 'test-token', got %v", resp["token"])
	}
}

func TestRegisterUser_InvalidInput(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// missing email field, password too short
	c.Request, _ = http.NewRequest(http.MethodPost, "/register",
		strings.NewReader(`{"password":"123"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	RegisterUser(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRegisterUser_ModelError(t *testing.T) {
	orig := modelRegisterUser
	defer func() { modelRegisterUser = orig }()

	modelRegisterUser = func(email, password string) (int64, error) {
		return 0, errors.New("duplicate email")
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/register",
		strings.NewReader(`{"email":"test@example.com","password":"pass123"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	RegisterUser(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestRegisterUser_TokenError(t *testing.T) {
	orig1, orig2 := modelRegisterUser, generateToken
	defer func() { modelRegisterUser, generateToken = orig1, orig2 }()

	modelRegisterUser = func(email, password string) (int64, error) { return 1, nil }
	generateToken = func(userID int64, email string) (string, error) {
		return "", errors.New("signing error")
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/register",
		strings.NewReader(`{"email":"test@example.com","password":"pass123"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	RegisterUser(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestLoginUser_Success(t *testing.T) {
	orig1, orig2 := modelAuthenticateUser, generateToken
	defer func() { modelAuthenticateUser, generateToken = orig1, orig2 }()

	modelAuthenticateUser = func(email, password string) (*models.User, error) {
		return &models.User{ID: 1, Email: email}, nil
	}
	generateToken = func(userID int64, email string) (string, error) { return "test-token", nil }

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/login",
		strings.NewReader(`{"email":"test@example.com","password":"pass123"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	LoginUser(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["token"] != "test-token" {
		t.Errorf("expected token 'test-token', got %v", resp["token"])
	}
}

func TestLoginUser_InvalidInput(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// missing password field
	c.Request, _ = http.NewRequest(http.MethodPost, "/login",
		strings.NewReader(`{"email":"test@example.com"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	LoginUser(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestLoginUser_AuthError(t *testing.T) {
	orig := modelAuthenticateUser
	defer func() { modelAuthenticateUser = orig }()

	modelAuthenticateUser = func(email, password string) (*models.User, error) {
		return nil, errors.New("invalid credentials")
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/login",
		strings.NewReader(`{"email":"test@example.com","password":"wrongpass"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	LoginUser(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}
