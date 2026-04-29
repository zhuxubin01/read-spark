package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/readspark/backend/internal/domain"
	"github.com/readspark/backend/internal/repository"
)

type SubscriptionService struct {
	subRepo *repository.SubscriptionRepository
}

func NewSubscriptionService(subRepo *repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{subRepo: subRepo}
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, userID uuid.UUID, req domain.CreateSubscriptionRequest) (*domain.Subscription, error) {
	now := time.Now()
	endDate := now.AddDate(0, 1, 0)
	if req.PlanType == "yearly" {
		endDate = now.AddDate(1, 0, 0)
	}

	paymentChannel := req.PaymentChannel
	transactionID := "mock-" + uuid.New().String()
	sub := &domain.Subscription{
		ID:             uuid.New(),
		UserID:         userID,
		PlanType:       req.PlanType,
		Status:         "active",
		StartDate:      now,
		EndDate:        endDate,
		AutoRenew:      true,
		PaymentChannel: &paymentChannel,
		TransactionID:  &transactionID,
	}

	if err := s.subRepo.Create(ctx, sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *SubscriptionService) GetStatus(ctx context.Context, userID uuid.UUID) (*domain.SubscriptionStatus, error) {
	sub, err := s.subRepo.FindActiveByUser(ctx, userID, time.Now())
	if err != nil {
		return nil, err
	}
	if sub == nil {
		return &domain.SubscriptionStatus{IsSubscribed: false}, nil
	}

	return &domain.SubscriptionStatus{
		IsSubscribed: true,
		PlanType:     &sub.PlanType,
		EndDate:      &sub.EndDate,
		AutoRenew:    &sub.AutoRenew,
	}, nil
}
