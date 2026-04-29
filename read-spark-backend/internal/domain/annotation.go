package domain

import (
	"time"

	"github.com/google/uuid"
)

type Annotation struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	UserID     uuid.UUID `gorm:"not null;index" json:"user_id"`
	ArticleID  uuid.UUID `gorm:"not null;index" json:"article_id"`
	Type       string    `gorm:"not null;size:20" json:"type"`
	RangeStart int       `gorm:"not null" json:"range_start"`
	RangeEnd   int       `gorm:"not null" json:"range_end"`
	Content    *string   `gorm:"type:text" json:"content,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateAnnotationRequest struct {
	ArticleID  uuid.UUID `json:"article_id" binding:"required"`
	Type       string    `json:"type" binding:"required,oneof=highlight note vocabulary"`
	RangeStart int       `json:"range_start" binding:"min=0"`
	RangeEnd   int       `json:"range_end" binding:"min=0"`
	Content    *string   `json:"content,omitempty"`
}
