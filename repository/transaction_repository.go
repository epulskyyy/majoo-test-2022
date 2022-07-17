package repository

import (
	"github.com/epulskyyy/majoo-test-2022/model"
	"github.com/epulskyyy/majoo-test-2022/utils/pagination"
	"gorm.io/gorm"
	"log"
	"math"
)

type ITransactionRepository interface {
	GetAllTrx(userId string, pagination pagination.Pagination) (*pagination.Pagination, error)
	GetOneById(id string) (*model.Transaction, error)
	GetAllReport(userId, year, month string, totalDay int, pagination pagination.Pagination) (*pagination.Pagination, error)
	GetAllReport2(userId, year, month string, total *int64) ([]*model.TransactionReport, error)
}
type TransactionRepository struct {
	db *gorm.DB
}

func (u TransactionRepository) GetAllReport2(userId, year, month string, total *int64) ([]*model.TransactionReport, error) {
	var transactions []*model.TransactionReport
	u.db.Table("transactions t").
		Joins("JOIN merchants b ON b.id = t.merchant_id").
		Joins("JOIN outlets o ON o.id = t.outlet_id").
		Select("b.merchant_name as merchant, o.outlet_name as outlet, sum(t.bill_total) as omzet, DATE(t.created_at) as date").
		Where("b.user_id = ? AND DATE_FORMAT(t.created_at,'%Y-%m') = ?", userId, year+"-"+month).Count(total)
	err := u.db.Debug().Table("transactions t").
		Joins("JOIN merchants b ON b.id = t.merchant_id").
		Joins("JOIN outlets o ON o.id = t.outlet_id").
		Select("b.merchant_name as merchant, o.outlet_name as outlet, sum(t.bill_total) as omzet, DATE(t.created_at) as date").
		Where("b.user_id = ? AND DATE_FORMAT(t.created_at,'%Y-%m') = ?", userId, year+"-"+month).
		Group("t.merchant_id, t.outlet_id, date").Scan(&transactions).Error
	if err != nil {
		return nil, err
	}
	log.Println(transactions)
	return transactions, nil
}

func (u TransactionRepository) GetAllReport(userId, year, month string, totalDay int, paginationIn pagination.Pagination) (*pagination.Pagination, error) {
	var transactions []*model.TransactionReport
	var totalRows int64
	u.db.Table("transactions t").
		Joins("JOIN merchants b ON b.id = t.merchant_id").
		Joins("JOIN outlets o ON o.id = t.outlet_id").
		Select("b.merchant_name as merchant, o.outlet_name as outlet, sum(t.bill_total) as omzet, DATE(t.created_at) as date").
		Where("b.user_id = ? AND DATE_FORMAT(t.created_at,'%Y-%m') = ?", userId, year+"-"+month).Count(&totalRows)
	err := u.db.Table("transactions t").
		Joins("JOIN merchants b ON b.id = t.merchant_id").
		Joins("JOIN outlets o ON o.id = t.outlet_id").
		Select("b.merchant_name as merchant, o.outlet_name as outlet, sum(t.bill_total) as omzet, DATE(t.created_at) as date").
		Where("b.user_id = ? AND DATE_FORMAT(t.created_at,'%Y-%m') = ?", userId, year+"-"+month).Offset(paginationIn.GetOffset()).
		Limit(paginationIn.GetLimit()).Order(paginationIn.GetSort()).
		Group("t.merchant_id, t.outlet_id, date").Scan(&transactions).Error
	paginationIn.Rows = transactions
	//if transactions == nil {
	//	totalDay = 0
	//}
	//totalRows := int64(totalDay)
	paginationIn.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(paginationIn.Limit)))
	paginationIn.TotalPages = totalPages
	log.Println(transactions)

	if err != nil {
		return nil, err
	}
	return &paginationIn, nil
}

func (u TransactionRepository) GetAllTrx(userId string, paginationIn pagination.Pagination) (*pagination.Pagination, error) {
	var transactions []*model.Transaction
	var totalRows int64
	u.db.Model(transactions).Joins("JOIN merchants b ON b.id = transactions.merchant_id").Where("b.user_id = ?", userId).Count(&totalRows)
	paginationIn.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(paginationIn.Limit)))
	paginationIn.TotalPages = totalPages
	err := u.db.Preload("Outlet").Preload("Merchant").Joins("JOIN merchants b ON b.id = transactions.merchant_id").Where("b.user_id = ?", userId).Offset(paginationIn.GetOffset()).Limit(paginationIn.GetLimit()).Order(paginationIn.GetSort()).Find(&transactions).Error
	paginationIn.Rows = transactions
	if err != nil {
		return nil, err
	}
	return &paginationIn, nil
}

func (u TransactionRepository) GetOneById(id string) (*model.Transaction, error) {
	var transaction model.Transaction
	err := u.db.Debug().First(&transaction, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func NewTransactionRepository(resource *gorm.DB) ITransactionRepository {
	return &TransactionRepository{db: resource}
}
