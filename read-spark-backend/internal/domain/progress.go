package domain

import (
	"time"

	"github.com/google/uuid"
)

type ReadingProgress struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	UserID     uuid.UUID `gorm:"not null;uniqueIndex:idx_user_article" json:"user_id"`
	ArticleID  uuid.UUID `gorm:"not null;uniqueIndex:idx_user_article" json:"article_id"`
	Position   int       `gorm:"default:0" json:"position"`
	Percentage float64   `gorm:"default:0" json:"percentage"`
	LastReadAt time.Time `json:"last_read_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (ReadingProgress) TableName() string {
	return "reading_progress"
}

type SyncProgressRequest struct {
	ArticleID  uuid.UUID `json:"article_id" binding:"required"`
	Position   int       `json:"position" binding:"min=0"`
	Percentage float64   `json:"percentage" binding:"min=0,max=100"`
}
