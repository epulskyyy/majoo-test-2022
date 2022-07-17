package model

import (
	"time"
)

type User struct {
	Id        uint      `json:"id" gorm:"primaryKey;not null"`
	Name      string    `json:"name" gorm:"size:45"`
	UserName  string    `json:"user_name" gorm:"size:45"`
	Password  string    `json:"-"`
	Merchant  Merchant  `json:"merchant,omitempty" swaggerignore:"true"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"type:timestamp;not null;default:current_timestamp;"`
	CreatedBy uint      `json:"created_by,omitempty" gorm:"not null"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"type:timestamp;not null;default:current_timestamp ON update current_timestamp;"`
	UpdatedBy uint      `json:"updated_by,omitempty" gorm:"not null"`
}
