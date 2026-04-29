package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/readspark/backend/internal/domain"
)

type AnnotationRepository struct {
	db *gorm.DB
}

func NewAnnotationRepository(db *gorm.DB) *AnnotationRepository {
	return &AnnotationRepository{db: db}
}

func (r *AnnotationRepository) Create(ctx context.Context, annotation *domain.Annotation) error {
	return r.db.WithContext(ctx).Create(annotation).Error
}

func (r *AnnotationRepository) ListByUser(ctx context.Context, userID uuid.UUID, articleID *uuid.UUID) ([]domain.Annotation, error) {
	query := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC")
	if articleID != nil {
		query = query.Where("article_id = ?", *articleID)
	}

	var rows []domain.Annotation
	if err := query.Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
