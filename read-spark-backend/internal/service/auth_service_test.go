package service

import (
	"context"
	"testing"

	"github.com/readspark/backend/internal/config"
)

func TestAuthServiceVerifyCode_Configurable(t *testing.T) {
	svc := NewAuthService(nil, config.JWTConfig{Secret: "x"}, config.AuthConfig{VerificationCode: "654321"})
	if err := svc.VerifyCode(context.Background(), "13800138000", "654321"); err != nil {
		t.Fatalf("expected code pass, got err=%v", err)
	}
	if err := svc.VerifyCode(context.Background(), "13800138000", "123456"); err == nil {
		t.Fatalf("expected invalid code error")
	}
}

func TestAuthServiceVerifyCode_Default(t *testing.T) {
	svc := NewAuthService(nil, config.JWTConfig{Secret: "x"}, config.AuthConfig{})
	if err := svc.VerifyCode(context.Background(), "13800138000", "123456"); err != nil {
		t.Fatalf("expected default code pass, got err=%v", err)
	}
}
