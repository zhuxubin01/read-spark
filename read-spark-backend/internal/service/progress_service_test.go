package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/readspark/backend/internal/domain"
	"github.com/readspark/backend/internal/repository"
)

func newProgressServiceForTest(t *testing.T) *ProgressService {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+uuid.NewString()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&domain.ReadingProgress{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	repo := repository.NewProgressRepository(db)
	return NewProgressService(repo)
}

func TestProgressServiceSyncAndList(t *testing.T) {
	svc := newProgressServiceForTest(t)
	ctx := context.Background()
	userID := uuid.New()
	articleID := uuid.New()

	cases := []struct {
		name       string
		position   int
		percentage float64
	}{
		{name: "create", position: 10, percentage: 5.5},
		{name: "update", position: 100, percentage: 66.6},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := svc.SyncProgress(ctx, userID, domain.SyncProgressRequest{ArticleID: articleID, Position: tc.position, Percentage: tc.percentage})
			if err != nil {
				t.Fatalf("sync failed: %v", err)
			}
			if out.Position != tc.position {
				t.Fatalf("position mismatch: want %d got %d", tc.position, out.Position)
			}
		})
	}

	rows, err := svc.ListProgress(ctx, userID)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(rows) != 1 || rows[0].Position != 100 {
		t.Fatalf("unexpected rows: %+v", rows)
	}
}
