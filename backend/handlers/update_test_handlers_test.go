package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"main.go/models"
)

func TestUpdate_Success(t *testing.T) {
	orig := modelUpdateTest
	defer func() { modelUpdateTest = orig }()

	called := false
	modelUpdateTest = func(test models.Test, testId int64) error {
		called = true
		return nil
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	body := `{"name":"updated","url":"https://example.com","method":"GET"}`
	c.Request, _ = http.NewRequest(http.MethodPut, "/tests/1/edit", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	Update(c)

	// handler does not write a response on success
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !called {
		t.Error("expected modelUpdateTest to be called")
	}
}

func TestUpdate_InvalidID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}
	c.Request, _ = http.NewRequest(http.MethodPut, "/tests/abc/edit",
		strings.NewReader(`{"name":"x","url":"https://x.com","method":"GET"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	Update(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUpdate_InvalidBody(t *testing.T) {
	orig := modelUpdateTest
	defer func() { modelUpdateTest = orig }()

	modelUpdateTest = func(test models.Test, testId int64) error { return errors.New("should not be called") }

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request, _ = http.NewRequest(http.MethodPut, "/tests/1/edit",
		strings.NewReader(`not json`))
	c.Request.Header.Set("Content-Type", "application/json")

	Update(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUpdate_ModelError(t *testing.T) {
	orig := modelUpdateTest
	defer func() { modelUpdateTest = orig }()

	modelUpdateTest = func(test models.Test, testId int64) error { return errors.New("db error") }

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	body := `{"name":"updated","url":"https://example.com","method":"GET"}`
	c.Request, _ = http.NewRequest(http.MethodPut, "/tests/1/edit", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	Update(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}
