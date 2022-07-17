package repository

import (
	"github.com/epulskyyy/majoo-test-2022/model"
	"gorm.io/gorm"
)

type IOutletRepository interface{
	GetAll(userId string) ([]model.Outlet, error)
	GetOneById(id string) (*model.Outlet, error)
}
type OutletRepository struct {
	db *gorm.DB
}

func (u OutletRepository) GetAll(userId string) ([]model.Outlet, error) {
	var outlets []model.Outlet
	err := u.db.Find(&outlets,"user_id = ?", userId).Error
	if err != nil {
		return nil, err
	}
	return outlets, nil
}

func (u OutletRepository) GetOneById(id string) (*model.Outlet, error) {
	var outlet model.Outlet
	err:= u.db.Preload("Merchant").First(&outlet,"id = ?",id).Error
	if err != nil {
		return nil, err
	}
	return &outlet, nil
}

func NewOutletRepository(resource *gorm.DB) IOutletRepository {
	return &OutletRepository{db: resource}
}
