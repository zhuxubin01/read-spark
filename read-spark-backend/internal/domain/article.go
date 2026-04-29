package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Article struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	Title       string     `gorm:"not null;size:500" json:"title"`
	Summary     *string    `json:"summary,omitempty"`
	Content     string     `gorm:"not null;type:text" json:"content"`
	Translation *string    `gorm:"type:text" json:"translation,omitempty"`
	Category    string     `gorm:"not null;size:50" json:"category"`
	Difficulty  string     `gorm:"not null;size:10" json:"difficulty"`
	WordCount   int        `gorm:"default:0" json:"word_count"`
	AudioURL    *string    `gorm:"size:500" json:"audio_url,omitempty"`
	CoverImage  *string    `gorm:"size:500" json:"cover_image,omitempty"`
	IsPremium   bool       `gorm:"default:true" json:"is_premium"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type ArticleSummary struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Summary     *string    `json:"summary,omitempty"`
	Category    string     `json:"category"`
	Difficulty  string     `json:"difficulty"`
	WordCount   int        `json:"word_count"`
	CoverImage  *string    `json:"cover_image,omitempty"`
	IsPremium   bool       `json:"is_premium"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
}

type SearchQuery struct {
	Keyword    string
	Category   *string
	Difficulty *string
	Page       int
	PageSize   int
}

type SearchResult struct {
	Articles []ArticleSummary
	Total    int64
}

// ArticleSearcher interface - allows swapping PG full-text for Elasticsearch later
type ArticleSearcher interface {
	Search(ctx context.Context, query SearchQuery) (SearchResult, error)
	Index(ctx context.Context, article Article) error
	Delete(ctx context.Context, articleID uuid.UUID) error
}
