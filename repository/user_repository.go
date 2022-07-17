package repository

import (
	"github.com/epulskyyy/majoo-test-2022/model"
	"gorm.io/gorm"
)

type IUserRepository interface{
	GetAll() ([]model.User, error)
	GetOneById(id string) (*model.User, error)
	GetOneByUserName(username string) (*model.User, error)
	CreateOne(user model.User) (*model.User, error)
}
type UserRepository struct {
	db *gorm.DB
}

func (u UserRepository) GetOneByUserName(username string) (*model.User, error) {
	var user model.User
	err:= u.db.Where("user_name = ?",username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u UserRepository) GetAll() ([]model.User, error) {
	var users []model.User
	err := u.db.Model(model.User{}).Scan(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u UserRepository) GetOneById(id string) (*model.User, error) {
	var user model.User
	err:= u.db.Preload("Merchant").First(&user,"id = ?",id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u UserRepository) CreateOne(user model.User) (*model.User, error) {
	err := u.db.Select("ID","Name", "Age", "JoinDate").Create(&user)
	if err.RowsAffected == 0 {
		return nil, err.Error
	}
	return &user, nil
}

func NewUserRepository(resource *gorm.DB) IUserRepository {
	return &UserRepository{db: resource}
}
