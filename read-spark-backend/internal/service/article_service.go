package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/readspark/backend/internal/domain"
	"github.com/readspark/backend/internal/repository"
)

type ArticleService struct {
	articleRepo *repository.ArticleRepository
	searcher    domain.ArticleSearcher
}

func NewArticleService(articleRepo *repository.ArticleRepository, searcher domain.ArticleSearcher) *ArticleService {
	return &ArticleService{
		articleRepo: articleRepo,
		searcher:    searcher,
	}
}

func (s *ArticleService) GetDailyArticles(ctx context.Context) ([]domain.ArticleSummary, error) {
	articles, err := s.articleRepo.FindDaily(ctx, 5)
	if err != nil {
		return nil, err
	}

	summaries := make([]domain.ArticleSummary, len(articles))
	for i, a := range articles {
		summaries[i] = toSummary(a)
	}
	return summaries, nil
}

func (s *ArticleService) GetArticle(ctx context.Context, id uuid.UUID, userID uuid.UUID, isSubscribed bool) (*domain.Article, error) {
	article, err := s.articleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Premium check
	if article.IsPremium && !isSubscribed {
		// Return preview only
		preview := *article
		preview.Content = preview.Content[:min(len(preview.Content), 500)] + "..."
		preview.Translation = nil
		return &preview, nil
	}

	return article, nil
}

func (s *ArticleService) ListArticles(ctx context.Context, category, difficulty *string, page, pageSize int) ([]domain.ArticleSummary, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	articles, total, err := s.articleRepo.FindAll(ctx, category, difficulty, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	summaries := make([]domain.ArticleSummary, len(articles))
	for i, a := range articles {
		summaries[i] = toSummary(a)
	}
	return summaries, total, nil
}

func (s *ArticleService) SearchArticles(ctx context.Context, keyword string, page, pageSize int) (domain.SearchResult, error) {
	return s.searcher.Search(ctx, domain.SearchQuery{
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
