package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/readspark/backend/internal/domain"
	"github.com/readspark/backend/internal/repository"
)

type SubscriptionService struct {
	subRepo         *repository.SubscriptionRepository
	receiptVerifier ReceiptVerifier
}

func NewSubscriptionService(subRepo *repository.SubscriptionRepository, receiptVerifier ReceiptVerifier) *SubscriptionService {
	if receiptVerifier == nil {
		receiptVerifier = &MockReceiptVerifier{}
	}
	return &SubscriptionService{subRepo: subRepo, receiptVerifier: receiptVerifier}
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, userID uuid.UUID, req domain.CreateSubscriptionRequest) (*domain.Subscription, error) {
	if err := s.receiptVerifier.Verify(ctx, req.PaymentChannel, req.Receipt); err != nil {
		return nil, err
	}

	now := time.Now()
	endDate := now.AddDate(0, 1, 0)
	if req.PlanType == "yearly" {
		endDate = now.AddDate(1, 0, 0)
	}

	paymentChannel := req.PaymentChannel
	transactionID := "verified-" + uuid.New().String()
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
