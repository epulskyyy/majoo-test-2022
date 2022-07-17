package model

import "time"

type Outlet struct {
	Id          uint          `json:"id" gorm:"primaryKey;not null"`
	MerchantId  uint          `json:"merchant_id,omitempty" gorm:"not null;"`
	Merchant  	*Merchant  	`json:"merchant,omitempty" swaggerignore:"true"`
	Transactions []Transaction `json:"transactions,omitempty" swaggerignore:"true"`
	OutletName  string        `json:"outlet_name" gorm:"size:45"`
	CreatedAt   *time.Time    `json:"created_at,omitempty" gorm:"type:timestamp;not null;default:current_timestamp;"`
	CreatedBy   uint          `json:"created_by,omitempty" gorm:"not null"`
	UpdatedAt   *time.Time    `json:"updated_at,omitempty" gorm:"type:timestamp;not null;default:current_timestamp ON update current_timestamp;"`
	UpdatedBy   uint          `json:"updated_by,omitempty" gorm:"not null"`
}
