package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/readspark/backend/internal/domain"
	"github.com/readspark/backend/internal/repository"
)

func TestDBProgressRepository_UpsertAndList(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:"+uuid.NewString()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&domain.ReadingProgress{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	repo := repository.NewProgressRepository(db)
	ctx := context.Background()
	userID := uuid.New()
	articleID := uuid.New()

	first, err := repo.Upsert(ctx, userID, domain.SyncProgressRequest{ArticleID: articleID, Position: 11, Percentage: 10})
	if err != nil {
		t.Fatalf("upsert create: %v", err)
	}
	if first.Position != 11 {
		t.Fatalf("unexpected first position: %d", first.Position)
	}

	second, err := repo.Upsert(ctx, userID, domain.SyncProgressRequest{ArticleID: articleID, Position: 88, Percentage: 77.7})
	if err != nil {
		t.Fatalf("upsert update: %v", err)
	}
	if second.Position != 88 {
		t.Fatalf("unexpected second position: %d", second.Position)
	}

	rows, err := repo.ListByUser(ctx, userID)
	if err != nil {
		t.Fatalf("list by user: %v", err)
	}
	if len(rows) != 1 || rows[0].Position != 88 {
		t.Fatalf("unexpected rows: %+v", rows)
	}
}
