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

const (
	appleSandboxURL    = "https://sandbox.itunes.apple.com/verifyReceipt"
	appleProductionURL = "https://buy.itunes.apple.com/verifyReceipt"
)

var (
	ErrReceiptRequired = errors.New("receipt is required")
	ErrReceiptInvalid  = errors.New("receipt is invalid")
	ErrReceiptUpstream = errors.New("receipt verifier upstream error")
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
	client         *http.Client
	verifyURL      string
	sandboxURL     string
	sharedSecret   string
	productionMode bool
}

func NewAppleReceiptVerifier(verifyURL, sharedSecret string, client *http.Client) *AppleReceiptVerifier {
	if strings.TrimSpace(verifyURL) == "" {
		verifyURL = appleSandboxURL
	}
	if client == nil {
		client = &http.Client{Timeout: 8 * time.Second}
	}
	return &AppleReceiptVerifier{client: client, verifyURL: verifyURL, sandboxURL: appleSandboxURL, sharedSecret: sharedSecret, productionMode: verifyURL == appleProductionURL}
}

func (v *AppleReceiptVerifier) Verify(ctx context.Context, channel, receipt string) error {
	if strings.TrimSpace(receipt) == "" {
		return ErrReceiptRequired
	}
	if strings.TrimSpace(channel) != "apple" {
		return fmt.Errorf("apple verifier does not support channel: %s", channel)
	}

	status, err := v.verifyAgainst(ctx, v.verifyURL, receipt)
	if err != nil {
		return err
	}
	if status == 0 {
		return nil
	}

	if status == 21007 && v.productionMode {
		sandboxStatus, fallbackErr := v.verifyAgainst(ctx, v.sandboxURL, receipt)
		if fallbackErr != nil {
			return fallbackErr
		}
		if sandboxStatus == 0 {
			return nil
		}
		return fmt.Errorf("%w: apple status=%d", ErrReceiptInvalid, sandboxStatus)
	}

	return fmt.Errorf("%w: apple status=%d", ErrReceiptInvalid, status)
}

func (v *AppleReceiptVerifier) verifyAgainst(ctx context.Context, endpoint, receipt string) (int, error) {
	payload := map[string]string{"receipt-data": receipt}
	if strings.TrimSpace(v.sharedSecret) != "" {
		payload["password"] = v.sharedSecret
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := v.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrReceiptUpstream, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, fmt.Errorf("%w: status=%d", ErrReceiptUpstream, resp.StatusCode)
	}

	var result struct {
		Status int `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}
	return result.Status, nil
}

func NewReceiptVerifier(cfg config.BillingConfig) ReceiptVerifier {
	switch strings.ToLower(strings.TrimSpace(cfg.ReceiptProvider)) {
	case "apple":
		verifyURL := cfg.AppleVerifyURL
		if strings.TrimSpace(verifyURL) == "" {
			if strings.ToLower(strings.TrimSpace(cfg.AppleEnvironment)) == "production" {
				verifyURL = appleProductionURL
			} else {
				verifyURL = appleSandboxURL
			}
		}
		return NewAppleReceiptVerifier(verifyURL, cfg.AppleSharedSecret, nil)
	default:
		return &MockReceiptVerifier{}
	}
}
