package service

import (
	"context"
	"errors"
	"strings"
)

type ReceiptVerifier interface {
	Verify(ctx context.Context, channel, receipt string) error
}

type MockReceiptVerifier struct{}

func (v *MockReceiptVerifier) Verify(ctx context.Context, channel, receipt string) error {
	_ = ctx
	_ = channel
	if strings.TrimSpace(receipt) == "" {
		return errors.New("receipt is required")
	}
	return nil
}
