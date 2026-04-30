package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/readspark/backend/internal/domain"
	"github.com/readspark/backend/internal/repository"
)

func newSubscriptionServiceForTest(t *testing.T, verifier ReceiptVerifier) *SubscriptionService {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+uuid.NewString()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&domain.Subscription{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	repo := repository.NewSubscriptionRepository(db)
	return NewSubscriptionService(repo, verifier)
}

type failVerifier struct{}

func (f *failVerifier) Verify(ctx context.Context, channel, receipt string) error {
	return ErrReceiptInvalid
}

func TestSubscriptionServiceCreate_VerifierError(t *testing.T) {
	svc := newSubscriptionServiceForTest(t, &failVerifier{})
	_, err := svc.CreateSubscription(context.Background(), uuid.New(), domain.CreateSubscriptionRequest{PlanType: "monthly", Receipt: "x", PaymentChannel: "apple"})
	if err == nil {
		t.Fatalf("expected verifier error")
	}
}

func TestSubscriptionServiceCreateAndStatus(t *testing.T) {
	svc := newSubscriptionServiceForTest(t, &MockReceiptVerifier{})
	userID := uuid.New()

	sub, err := svc.CreateSubscription(context.Background(), userID, domain.CreateSubscriptionRequest{PlanType: "yearly", Receipt: "ok", PaymentChannel: "apple"})
	if err != nil {
		t.Fatalf("create subscription failed: %v", err)
	}
	if sub.EndDate.Before(time.Now().AddDate(0, 11, 0)) {
		t.Fatalf("expected yearly end date")
	}

	status, err := svc.GetStatus(context.Background(), userID)
	if err != nil {
		t.Fatalf("status failed: %v", err)
	}
	if !status.IsSubscribed || status.PlanType == nil || *status.PlanType != "yearly" {
		t.Fatalf("unexpected status: %+v", status)
	}
}
