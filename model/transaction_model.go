package model

import "time"

type Transaction struct {
	Id         uint       `json:"id" gorm:"primaryKey;not null"`
	MerchantId uint       `json:"merchant_id,omitempty" gorm:"not null;"`
	Merchant   *Merchant  `json:"merchant,omitempty" swaggerignore:"true"`
	Outlet     *Outlet    `json:"outlet,omitempty" swaggerignore:"true"`
	OutletId   uint       `json:"outlet_id,omitempty" gorm:"not null;"`
	BillTotal  float64    `json:"bill_total" gorm:"not null"`
	CreatedAt  *time.Time `json:"created_at,omitempty" gorm:"type:timestamp;not null;default:current_timestamp;"`
	CreatedBy  uint       `json:"created_by,omitempty" gorm:"not null"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty" gorm:"type:timestamp;not null;default:current_timestamp ON update current_timestamp;"`
	UpdatedBy  uint       `json:"updated_by,omitempty" gorm:"not null"`
}

type TransactionReport struct{
	Merchant string `json:"merchant"`
	Outlet string `json:"outlet"`
	Omzet float64 `json:"omzet"`
	Date  *time.Time `json:"date,omitempty" gorm:"type:timestamp"`
}