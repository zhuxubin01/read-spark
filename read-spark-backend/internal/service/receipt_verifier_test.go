package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/readspark/backend/internal/config"
)

func TestNewReceiptVerifier_MockDefault(t *testing.T) {
	v := NewReceiptVerifier(config.BillingConfig{ReceiptProvider: ""})
	if _, ok := v.(*MockReceiptVerifier); !ok {
		t.Fatalf("expected mock verifier")
	}
}

func TestNewReceiptVerifier_Apple(t *testing.T) {
	v := NewReceiptVerifier(config.BillingConfig{ReceiptProvider: "apple", AppleEnvironment: "sandbox"})
	if _, ok := v.(*AppleReceiptVerifier); !ok {
		t.Fatalf("expected apple verifier")
	}
}

func TestAppleReceiptVerifier_VerifySuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":0}`))
	}))
	defer ts.Close()

	v := NewAppleReceiptVerifier(ts.URL, "", ts.Client())
	if err := v.Verify(context.Background(), "apple", "receipt-data"); err != nil {
		t.Fatalf("expected success, got %v", err)
	}
}

func TestAppleReceiptVerifier_VerifyInvalidStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":21010}`))
	}))
	defer ts.Close()

	v := NewAppleReceiptVerifier(ts.URL, "", ts.Client())
	if err := v.Verify(context.Background(), "apple", "receipt-data"); err == nil {
		t.Fatalf("expected invalid receipt error")
	}
}
