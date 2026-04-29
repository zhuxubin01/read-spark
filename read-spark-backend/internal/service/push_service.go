package service

import (
	"context"

	"github.com/google/uuid"
)

type RegisterPushTokenRequest struct {
	DeviceToken string `json:"device_token" binding:"required"`
	Platform    string `json:"platform" binding:"required,oneof=ios android"`
}

type PushService struct{}

func NewPushService() *PushService {
	return &PushService{}
}

func (s *PushService) RegisterToken(_ context.Context, _ uuid.UUID, req RegisterPushTokenRequest) map[string]any {
	return map[string]any{
		"registered":   true,
		"platform":     req.Platform,
		"device_token": req.DeviceToken,
	}
}
