package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/readspark/backend/internal/domain"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=readspark port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()
	summary := "Sleep is one of the most important aspects of our life..."
	content := `Sleep is one of the most important but least understood aspects of our life. Until very recently, science had no answer to the question of why we sleep, or what good it served, or why we suffer such devastating health consequences when it is absent.`
	translation := `睡眠是我们生活中最重要但了解最少的方面之一。直到最近，科学对为什么我们睡觉、它有什么好处、或者为什么缺少睡眠时我们会遭受如此严重的健康后果这些问题都没有答案。`

	articles := []domain.Article{
		{
			ID:          uuid.New(),
			Title:       "Why We Sleep",
			Summary:     &summary,
			Content:     content,
			Translation: &translation,
			Category:    "news",
			Difficulty:  "B2",
			WordCount:   450,
			IsPremium:   true,
			PublishedAt: &now,
		},
		{
			ID:          uuid.New(),
			Title:       "The Art of Saying No",
			Summary:     &[]string{"Learning to say no is a crucial skill..."}[0],
			Content:     `Many of us find it difficult to say no. We worry about disappointing others or being seen as uncooperative. However, learning to say no is essential for maintaining healthy boundaries and protecting our time and energy.`,
			Translation: &[]string{`我们中的许多人发现很难说不。我们担心让别人失望或被视为不合作。然而，学会说不对于维持健康的界限和保护我们的时间和精力至关重要。`}[0],
			Category:    "news",
			Difficulty:  "B1",
			WordCount:   320,
			IsPremium:   false,
			PublishedAt: &now,
		},
	}

	for _, a := range articles {
		if err := db.WithContext(context.Background()).Create(&a).Error; err != nil {
			log.Printf("failed to create article %s: %v", a.Title, err)
		} else {
			fmt.Printf("Created article: %s\n", a.Title)
		}
	}
}
