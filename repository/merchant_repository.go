package repository

import (
	"github.com/epulskyyy/majoo-test-2022/model"
	"gorm.io/gorm"
)

type IMerchantRepository interface{
	GetAll(userId string) ([]model.Merchant, error)
	GetOneById(id string) (*model.Merchant, error)
}
type MerchantRepository struct {
	db *gorm.DB
}

func (u MerchantRepository) GetAll(userId string) ([]model.Merchant, error) {
	var merchants []model.Merchant
	err := u.db.Find(&merchants,"user_id = ?", userId).Error
	if err != nil {
		return nil, err
	}
	return merchants, nil
}

func (u MerchantRepository) GetOneById(id string) (*model.Merchant, error) {
	var merchant model.Merchant
	err:= u.db.Preload("Outlets").Where("id = ?",id).First(&merchant).Error
	if err != nil {
		return nil, err
	}
	return &merchant, nil
}

func NewMerchantRepository(resource *gorm.DB) IMerchantRepository {
	return &MerchantRepository{db: resource}
}
