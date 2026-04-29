package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestProgressSync_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewProgressHandler(nil)
	r.POST("/progress", h.Sync)

	req := httptest.NewRequest(http.MethodPost, "/progress", strings.NewReader(`{"position":1}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestProgressList_MissingUserContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewProgressHandler(nil)
	r.GET("/progress", h.List)

	req := httptest.NewRequest(http.MethodGet, "/progress", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}
