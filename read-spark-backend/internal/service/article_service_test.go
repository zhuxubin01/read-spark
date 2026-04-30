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

type fakeSearcher struct {
	result domain.SearchResult
	err    error
}

func (f *fakeSearcher) Search(ctx context.Context, query domain.SearchQuery) (domain.SearchResult, error) {
	return f.result, f.err
}
func (f *fakeSearcher) Index(ctx context.Context, article domain.Article) error { return nil }
func (f *fakeSearcher) Delete(ctx context.Context, articleID uuid.UUID) error   { return nil }

func newArticleServiceForTest(t *testing.T) *ArticleService {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+uuid.NewString()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&domain.Article{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	repo := repository.NewArticleRepository(db)
	return NewArticleService(repo, &fakeSearcher{})
}

func TestArticleServiceGetArticle_PremiumPreview(t *testing.T) {
	svc := newArticleServiceForTest(t)
	ctx := context.Background()

	content := ""
	for i := 0; i < 700; i++ {
		content += "a"
	}
	translation := "zh"
	now := time.Now()
	article := &domain.Article{
		ID:          uuid.New(),
		Title:       "premium",
		Content:     content,
		Translation: &translation,
		Category:    "news",
		Difficulty:  "B2",
		IsPremium:   true,
		PublishedAt: &now,
	}
	if err := svc.articleRepo.Create(ctx, article); err != nil {
		t.Fatalf("create article: %v", err)
	}

	got, err := svc.GetArticle(ctx, article.ID, uuid.New(), false)
	if err != nil {
		t.Fatalf("get article: %v", err)
	}
	if got.Translation != nil {
		t.Fatalf("expected translation hidden for unsubscribed user")
	}
	if len(got.Content) != 503 {
		t.Fatalf("expected preview length 503, got %d", len(got.Content))
	}
}

func TestArticleServiceListArticles_NormalizesPagination(t *testing.T) {
	svc := newArticleServiceForTest(t)
	ctx := context.Background()
	now := time.Now()
	for i := 0; i < 3; i++ {
		a := &domain.Article{ID: uuid.New(), Title: "t", Content: "c", Category: "news", Difficulty: "B1", PublishedAt: &now}
		if err := svc.articleRepo.Create(ctx, a); err != nil {
			t.Fatalf("create article: %v", err)
		}
	}

	rows, total, err := svc.ListArticles(ctx, nil, nil, 0, 1000)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 3 {
		t.Fatalf("expected total 3, got %d", total)
	}
	if len(rows) != 3 {
		t.Fatalf("expected 3 rows with normalized page size, got %d", len(rows))
	}
}
