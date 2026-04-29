package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/readspark/backend/internal/domain"
)

type ProgressRepository struct {
	db *gorm.DB
}

func NewProgressRepository(db *gorm.DB) *ProgressRepository {
	return &ProgressRepository{db: db}
}

func (r *ProgressRepository) Upsert(ctx context.Context, userID uuid.UUID, req domain.SyncProgressRequest) (*domain.ReadingProgress, error) {
	var progress domain.ReadingProgress
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND article_id = ?", userID, req.ArticleID).
		First(&progress).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			progress = domain.ReadingProgress{
				ID:         uuid.New(),
				UserID:     userID,
				ArticleID:  req.ArticleID,
				Position:   req.Position,
				Percentage: req.Percentage,
				LastReadAt: time.Now(),
			}
			if createErr := r.db.WithContext(ctx).Create(&progress).Error; createErr != nil {
				return nil, createErr
			}
			return &progress, nil
		}
		return nil, err
	}

	progress.Position = req.Position
	progress.Percentage = req.Percentage
	progress.LastReadAt = time.Now()
	if err := r.db.WithContext(ctx).Save(&progress).Error; err != nil {
		return nil, err
	}

	return &progress, nil
}

func (r *ProgressRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.ReadingProgress, error) {
	var rows []domain.ReadingProgress
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		Find(&rows).Error
	return rows, err
}
