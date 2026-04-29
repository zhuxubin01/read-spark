package domain

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID         uuid.UUID  `gorm:"not null;index" json:"user_id"`
	PlanType       string     `gorm:"not null;size:20" json:"plan_type"`
	Status         string     `gorm:"not null;size:20" json:"status"`
	StartDate      time.Time  `gorm:"not null" json:"start_date"`
	EndDate        time.Time  `gorm:"not null" json:"end_date"`
	AutoRenew      bool       `gorm:"default:true" json:"auto_renew"`
	PaymentChannel *string    `gorm:"size:50" json:"payment_channel,omitempty"`
	TransactionID  *string    `gorm:"size:255" json:"transaction_id,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type CreateSubscriptionRequest struct {
	PlanType       string `json:"plan_type" binding:"required,oneof=monthly yearly"`
	Receipt        string `json:"receipt" binding:"required"`
	PaymentChannel string `json:"payment_channel" binding:"required,oneof=apple google wechat alipay"`
}

type SubscriptionStatus struct {
	IsSubscribed bool       `json:"is_subscribed"`
	PlanType     *string    `json:"plan_type,omitempty"`
	EndDate      *time.Time `json:"end_date,omitempty"`
	AutoRenew    *bool      `json:"auto_renew,omitempty"`
}
