package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/readspark/backend/internal/config"
)

var (
	ErrReceiptRequired = errors.New("receipt is required")
	ErrReceiptInvalid  = errors.New("receipt is invalid")
)

type ReceiptVerifier interface {
	Verify(ctx context.Context, channel, receipt string) error
}

type MockReceiptVerifier struct{}

func (v *MockReceiptVerifier) Verify(ctx context.Context, channel, receipt string) error {
	_ = ctx
	_ = channel
	if strings.TrimSpace(receipt) == "" {
		return ErrReceiptRequired
	}
	return nil
}

type AppleReceiptVerifier struct {
	client       *http.Client
	verifyURL    string
	sharedSecret string
}

func NewAppleReceiptVerifier(verifyURL, sharedSecret string, client *http.Client) *AppleReceiptVerifier {
	if strings.TrimSpace(verifyURL) == "" {
		verifyURL = "https://sandbox.itunes.apple.com/verifyReceipt"
	}
	if client == nil {
		client = &http.Client{Timeout: 8 * time.Second}
	}
	return &AppleReceiptVerifier{client: client, verifyURL: verifyURL, sharedSecret: sharedSecret}
}

func (v *AppleReceiptVerifier) Verify(ctx context.Context, channel, receipt string) error {
	if strings.TrimSpace(receipt) == "" {
		return ErrReceiptRequired
	}
	if strings.TrimSpace(channel) != "apple" {
		return fmt.Errorf("apple verifier does not support channel: %s", channel)
	}

	payload := map[string]string{"receipt-data": receipt}
	if strings.TrimSpace(v.sharedSecret) != "" {
		payload["password"] = v.sharedSecret
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, v.verifyURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := v.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("apple verify upstream error: %d", resp.StatusCode)
	}

	var result struct {
		Status int `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	if result.Status != 0 {
		return fmt.Errorf("%w: apple status=%d", ErrReceiptInvalid, result.Status)
	}
	return nil
}

func NewReceiptVerifier(cfg config.BillingConfig) ReceiptVerifier {
	switch strings.ToLower(strings.TrimSpace(cfg.ReceiptProvider)) {
	case "apple":
		verifyURL := cfg.AppleVerifyURL
		if strings.TrimSpace(verifyURL) == "" {
			if strings.ToLower(strings.TrimSpace(cfg.AppleEnvironment)) == "production" {
				verifyURL = "https://buy.itunes.apple.com/verifyReceipt"
			} else {
				verifyURL = "https://sandbox.itunes.apple.com/verifyReceipt"
			}
		}
		return NewAppleReceiptVerifier(verifyURL, cfg.AppleSharedSecret, nil)
	default:
		return &MockReceiptVerifier{}
	}
}
