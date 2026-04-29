package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/readspark/backend/internal/domain"
	"github.com/readspark/backend/internal/repository"
)

type ProgressService struct {
	progressRepo *repository.ProgressRepository
}

func NewProgressService(progressRepo *repository.ProgressRepository) *ProgressService {
	return &ProgressService{progressRepo: progressRepo}
}

func (s *ProgressService) SyncProgress(ctx context.Context, userID uuid.UUID, req domain.SyncProgressRequest) (*domain.ReadingProgress, error) {
	return s.progressRepo.Upsert(ctx, userID, req)
}

func (s *ProgressService) ListProgress(ctx context.Context, userID uuid.UUID) ([]domain.ReadingProgress, error) {
	return s.progressRepo.ListByUser(ctx, userID)
}
