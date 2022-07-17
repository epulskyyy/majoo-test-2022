package model

import "time"

type Merchant struct {
	Id           uint          `json:"id" gorm:"primaryKey;not null"`
	UserId       uint          `json:"user_id,omitempty" gorm:"not null;unique"`
	MerchantName string        `json:"merchant_name" gorm:"size:45"`
	Outlets       []Outlet      `json:"outlets,omitempty" swaggerignore:"true"`
	Transactions  []Transaction `json:"transactions,omitempty" swaggerignore:"true"`
	CreatedAt    *time.Time    `json:"created_at,omitempty" gorm:"type:timestamp;not null;default:current_timestamp;"`
	CreatedBy    uint          `json:"created_by,omitempty" gorm:"not null"`
	UpdatedAt    *time.Time    `json:"updated_at,omitempty" gorm:"type:timestamp;not null;default:current_timestamp ON update current_timestamp;"`
	UpdatedBy    uint          `json:"updated_by,omitempty" gorm:"not null"`
}
