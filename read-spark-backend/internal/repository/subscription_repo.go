package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/readspark/backend/internal/domain"
)

type SubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(ctx context.Context, sub *domain.Subscription) error {
	return r.db.WithContext(ctx).Create(sub).Error
}

func (r *SubscriptionRepository) FindLatestByUser(ctx context.Context, userID uuid.UUID) (*domain.Subscription, error) {
	var sub domain.Subscription
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("end_date DESC").
		First(&sub).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &sub, nil
}

func (r *SubscriptionRepository) FindActiveByUser(ctx context.Context, userID uuid.UUID, now time.Time) (*domain.Subscription, error) {
	var sub domain.Subscription
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ? AND end_date >= ?", userID, "active", now).
		Order("end_date DESC").
		First(&sub).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &sub, nil
}
