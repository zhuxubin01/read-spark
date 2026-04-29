package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/readspark/backend/internal/service"
)

func TestSubscriptionCreate_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewSubscriptionHandler(nil)
	r.POST("/subscriptions", h.Create)

	req := httptest.NewRequest(http.MethodPost, "/subscriptions", strings.NewReader(`{"plan_type":"monthly"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestSubscriptionStatus_MissingUserContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewSubscriptionHandler(nil)
	r.GET("/subscriptions/status", h.Status)

	req := httptest.NewRequest(http.MethodGet, "/subscriptions/status", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

type errVerifier struct{ err error }

func (e *errVerifier) Verify(ctx context.Context, channel, receipt string) error { return e.err }

func TestSubscriptionCreate_MapsReceiptErrors(t *testing.T) {
	cases := []struct {
		name string
		err  error
		code int
	}{
		{name: "invalid", err: service.ErrReceiptInvalid, code: http.StatusBadRequest},
		{name: "required", err: service.ErrReceiptRequired, code: http.StatusBadRequest},
		{name: "upstream", err: service.ErrReceiptUpstream, code: http.StatusBadGateway},
		{name: "internal", err: errors.New("db down"), code: http.StatusInternalServerError},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			r := gin.New()
			svc := service.NewSubscriptionService(nil, &errVerifier{err: tc.err})
			h := NewSubscriptionHandler(svc)
			r.POST("/subscriptions", func(c *gin.Context) {
				c.Set("userID", uuid.New())
				h.Create(c)
			})

			req := httptest.NewRequest(http.MethodPost, "/subscriptions", strings.NewReader(`{"plan_type":"monthly","receipt":"r","payment_channel":"apple"}`))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.code {
				t.Fatalf("expected %d, got %d, body=%s", tc.code, w.Code, w.Body.String())
			}
		})
	}
}
