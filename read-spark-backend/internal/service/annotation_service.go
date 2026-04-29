package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/readspark/backend/internal/domain"
	"github.com/readspark/backend/internal/repository"
)

type AnnotationService struct {
	annotationRepo *repository.AnnotationRepository
}

func NewAnnotationService(annotationRepo *repository.AnnotationRepository) *AnnotationService {
	return &AnnotationService{annotationRepo: annotationRepo}
}

func (s *AnnotationService) Create(ctx context.Context, userID uuid.UUID, req domain.CreateAnnotationRequest) (*domain.Annotation, error) {
	if req.RangeEnd < req.RangeStart {
		return nil, errors.New("range_end must be greater than or equal to range_start")
	}

	row := &domain.Annotation{
		ID:         uuid.New(),
		UserID:     userID,
		ArticleID:  req.ArticleID,
		Type:       req.Type,
		RangeStart: req.RangeStart,
		RangeEnd:   req.RangeEnd,
		Content:    req.Content,
	}
	if err := s.annotationRepo.Create(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *AnnotationService) List(ctx context.Context, userID uuid.UUID, articleID *uuid.UUID) ([]domain.Annotation, error) {
	return s.annotationRepo.ListByUser(ctx, userID, articleID)
}
