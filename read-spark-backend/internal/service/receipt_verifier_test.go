package service

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
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
	err := v.Verify(context.Background(), "apple", "receipt-data")
	if err == nil || !errors.Is(err, ErrReceiptInvalid) {
		t.Fatalf("expected invalid receipt error, got %v", err)
	}
}

func TestAppleReceiptVerifier_ProductionFallback21007(t *testing.T) {
	prod := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":21007}`))
	}))
	defer prod.Close()

	sandbox := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":0}`))
	}))
	defer sandbox.Close()

	v := NewAppleReceiptVerifier(prod.URL, "", prod.Client())
	v.productionMode = true
	v.sandboxURL = sandbox.URL

	if err := v.Verify(context.Background(), "apple", "receipt-data"); err != nil {
		t.Fatalf("expected fallback success, got %v", err)
	}
}

func TestAppleReceiptVerifier_UpstreamError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer ts.Close()

	v := NewAppleReceiptVerifier(ts.URL, "", ts.Client())
	err := v.Verify(context.Background(), "apple", "receipt-data")
	if err == nil || !errors.Is(err, ErrReceiptUpstream) {
		t.Fatalf("expected upstream error, got %v", err)
	}
}

func TestAppleReceiptVerifier_ChannelMismatch(t *testing.T) {
	v := NewAppleReceiptVerifier("", "", nil)
	err := v.Verify(context.Background(), "google", "receipt-data")
	if err == nil || !strings.Contains(err.Error(), "does not support") {
		t.Fatalf("expected channel mismatch error, got %v", err)
	}
}
