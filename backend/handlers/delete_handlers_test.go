package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestDelete_Success(t *testing.T) {
	orig1, orig2 := modelCheckTestExists, modelDeleteTest
	defer func() { modelCheckTestExists, modelDeleteTest = orig1, orig2 }()

	deleted := false
	modelCheckTestExists = func(testID int64) (bool, error) { return true, nil }
	// modify deleted value to check if modelDeleteTest is called or not
	modelDeleteTest = func(testID int64) error { deleted = true; return nil }

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request, _ = http.NewRequest(http.MethodDelete, "/tests/1/delete", nil)

	Delete(c)

	// handler does not write a response on success
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	// If modelDeleteTest is not called
	if !deleted {
		// Log followed by fail (doesn't stop the execution)
		t.Error("expected modelDeleteTest to be called")
	}
}

func TestDelete_InvalidID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "xyz"}}
	// invalid test id (in this case, sending nil)
	c.Request, _ = http.NewRequest(http.MethodDelete, "/tests/xyz/delete", nil)

	Delete(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestDelete_NotFound(t *testing.T) {
	orig := modelCheckTestExists
	defer func() { modelCheckTestExists = orig }()

	// return (false, nil) to exercise the !exists branch → 404
	modelCheckTestExists = func(testID int64) (bool, error) { return false, nil }

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "999"}}
	c.Request, _ = http.NewRequest(http.MethodDelete, "/tests/999/delete", nil)

	Delete(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestDelete_DBError(t *testing.T) {
	orig := modelCheckTestExists
	defer func() { modelCheckTestExists = orig }()

	modelCheckTestExists = func(testID int64) (bool, error) {
		return false, errors.New("connection error")
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request, _ = http.NewRequest(http.MethodDelete, "/tests/1/delete", nil)

	Delete(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestDelete_DeleteError(t *testing.T) {
	orig1, orig2 := modelCheckTestExists, modelDeleteTest
	defer func() { modelCheckTestExists, modelDeleteTest = orig1, orig2 }()

	modelCheckTestExists = func(testID int64) (bool, error) { return true, nil }
	modelDeleteTest = func(testID int64) error { return errors.New("db error") }

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request, _ = http.NewRequest(http.MethodDelete, "/tests/1/delete", nil)

	Delete(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}
