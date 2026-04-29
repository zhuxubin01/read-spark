package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/readspark/backend/internal/domain"
)

type ArticleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

func (r *ArticleRepository) Create(ctx context.Context, article *domain.Article) error {
	return r.db.WithContext(ctx).Create(article).Error
}

func (r *ArticleRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Article, error) {
	var article domain.Article
	if err := r.db.WithContext(ctx).First(&article, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrArticleNotFound
		}
		return nil, err
	}
	return &article, nil
}

func (r *ArticleRepository) FindDaily(ctx context.Context, limit int) ([]domain.Article, error) {
	var articles []domain.Article
	today := time.Now().Truncate(24 * time.Hour)
	err := r.db.WithContext(ctx).
		Where("published_at >= ?", today).
		Order("published_at DESC").
		Limit(limit).
		Find(&articles).Error
	return articles, err
}

func (r *ArticleRepository) FindAll(ctx context.Context, category, difficulty *string, page, pageSize int) ([]domain.Article, int64, error) {
	var articles []domain.Article
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Article{})
	if category != nil && *category != "" {
		query = query.Where("category = ?", *category)
	}
	if difficulty != nil && *difficulty != "" {
		query = query.Where("difficulty = ?", *difficulty)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("published_at DESC").Offset(offset).Limit(pageSize).Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

func (r *ArticleRepository) SearchFullText(ctx context.Context, keyword string, page, pageSize int) ([]domain.Article, int64, error) {
	var articles []domain.Article
	var total int64

	q := r.db.WithContext(ctx).
		Where("search_vector @@ plainto_tsquery('english', ?)", keyword).
		Order(gorm.Expr("ts_rank(search_vector, plainto_tsquery('english', ?)) DESC", keyword))

	if err := q.Model(&domain.Article{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := q.Offset(offset).Limit(pageSize).Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}
