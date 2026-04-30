package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/readspark/backend/internal/service"
)

func TestDictionaryLookup_Success(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"word":"hello","phonetic":"həˈləʊ","meanings":[{"definitions":[{"definition":"greeting"}]}]}]`))
	}))
	defer upstream.Close()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := service.NewDictionaryServiceWithClient(upstream.URL, upstream.Client())
	h := NewDictionaryHandler(svc)
	r.GET("/dictionary/:word", h.Lookup)

	req := httptest.NewRequest(http.MethodGet, "/dictionary/hello", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestDictionaryLookup_UpstreamError(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer upstream.Close()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := service.NewDictionaryServiceWithClient(upstream.URL, upstream.Client())
	h := NewDictionaryHandler(svc)
	r.GET("/dictionary/:word", h.Lookup)

	req := httptest.NewRequest(http.MethodGet, "/dictionary/hello", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadGateway {
		t.Fatalf("expected 502, got %d", w.Code)
	}
}
