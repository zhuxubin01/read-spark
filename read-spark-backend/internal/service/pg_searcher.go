package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/readspark/backend/internal/domain"
	"github.com/readspark/backend/internal/repository"
)

// PGFullTextSearch implements domain.ArticleSearcher using PostgreSQL full-text search
type PGFullTextSearch struct {
	articleRepo *repository.ArticleRepository
}

func NewPGFullTextSearch(articleRepo *repository.ArticleRepository) *PGFullTextSearch {
	return &PGFullTextSearch{articleRepo: articleRepo}
}

func (s *PGFullTextSearch) Search(ctx context.Context, query domain.SearchQuery) (domain.SearchResult, error) {
	articles, total, err := s.articleRepo.SearchFullText(ctx, query.Keyword, query.Page, query.PageSize)
	if err != nil {
		return domain.SearchResult{}, err
	}

	summaries := make([]domain.ArticleSummary, len(articles))
	for i, a := range articles {
		summaries[i] = toSummary(a)
	}

	return domain.SearchResult{Articles: summaries, Total: total}, nil
}

func (s *PGFullTextSearch) Index(ctx context.Context, article domain.Article) error {
	// PostgreSQL auto-updates tsvector, no manual indexing needed
	return nil
}

func (s *PGFullTextSearch) Delete(ctx context.Context, articleID uuid.UUID) error {
	// Handled by repository
	return nil
}

func toSummary(a domain.Article) domain.ArticleSummary {
	return domain.ArticleSummary{
		ID:          a.ID,
		Title:       a.Title,
		Summary:     a.Summary,
		Category:    a.Category,
		Difficulty:  a.Difficulty,
		WordCount:   a.WordCount,
		CoverImage:  a.CoverImage,
		IsPremium:   a.IsPremium,
		PublishedAt: a.PublishedAt,
	}
}
