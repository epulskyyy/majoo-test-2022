package usecase

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/epulskyyy/majoo-test-2022/httputil"
	"github.com/epulskyyy/majoo-test-2022/model"
	"github.com/epulskyyy/majoo-test-2022/repository"
	pagination2 "github.com/epulskyyy/majoo-test-2022/utils/pagination"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jinzhu/now"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ITransactionUseCase interface {
	GetById(id string) *httputil.ResponseMessage
	GetAllTrx(limit, page string) *httputil.ResponseMessage
	GetReportTrx(limit, page, month, year string) *httputil.ResponseMessage
	CreateCSV(month, year string, ctx *gin.Context)
}

type TransactionUseCase struct {
	repo         repository.ITransactionRepository
	repoMerchant repository.IMerchantRepository
	client       *redis.Client
	res          httputil.ResponseMessage
}

func (s *TransactionUseCase) CreateCSV(month, year string, ctx *gin.Context) {
	var user model.User
	result, err := s.client.Get("user_info").Result()
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(result), &user)
	if err != nil {
		return
	}
	userId := strconv.Itoa(int(user.Id))

	record := [][]string{{
		//header
		"Nama Merchant", "Nama Outlet","Omzet", "Tanggal",
	}}
	b := &bytes.Buffer{}
	w := csv.NewWriter(b)

	y, m, _ := time.Now().Date()
	if year == "" {
		year = strconv.Itoa(y)
	}
	if month == "" {
		month = strconv.Itoa(int(m))
	}
	var total int64
	transaction, err := s.repo.GetAllReport2(userId,year,month, &total)
	if err != nil {
		return
	}
	log.Println(transaction)
	for _, trx := range transaction {
		newTrx := []string{
			trx.Merchant,trx.Outlet,fmt.Sprintf("%.2f", trx.Omzet),
			strings.TrimSpace(trx.Date.Format("2006-01-02")),
		}
		record = append(record, newTrx)
	}

	w.WriteAll(record) // calls Flush internally
	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
		return
	}

	// Write any buffered data to the underlying writer (standard output).
	w.Flush()
	if err := w.Error(); err != nil {
		log.Println(err)
		return
	}
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Length", strconv.Itoa(b.Len()))
	ctx.Header("Content-Disposition", "attachment; filename=contacts.csv")
	ctx.Data(http.StatusOK, "text/csv", b.Bytes())
}

func (s *TransactionUseCase) GetReportTrx(limit, page, month, year string) *httputil.ResponseMessage {
	errorList := make(map[string]string)
	var user model.User
	result, err := s.client.Get("user_info").Result()
	if err != nil {
		errorList["message"] = err.Error()
		s.res.Errors(http.StatusBadRequest, "000", errorList)
		return &s.res
	}
	err = json.Unmarshal([]byte(result), &user)
	if err != nil {
		errorList["message"] = err.Error()
		s.res.Errors(http.StatusBadRequest, "000", errorList)
		return &s.res
	}
	userId := strconv.Itoa(int(user.Id))
	newLimit, errConv := strconv.Atoi(limit)
	if errConv != nil {
		newLimit = 10
	}
	newPage, errConv := strconv.Atoi(page)
	if errConv != nil {
		newPage = 1
	}
	pagination := pagination2.Pagination{
		Limit: newLimit,
		Page:  newPage,
		Sort: "date ASC",
	}
	y, m, _ := time.Now().Date()
	if year == "" {
		year = strconv.Itoa(y)
	}
	if month == "" {
		month = strconv.Itoa(int(m))
	}

	log.Println(year, month)
	t, err := now.Parse(year + "-" + month)
	if err != nil {
		errorList["message"] = err.Error()
		s.res.Errors(http.StatusBadRequest, "000", errorList)
		return &s.res
	}
	t = now.With(t).EndOfMonth()
	log.Println(t.Day())

	transactions, err := s.repo.GetAllReport(userId, year, month,t.Day(), pagination)
	if err != nil {
		errorList["message"] = err.Error()
		s.res.Errors(http.StatusBadRequest, "000", errorList)
		return &s.res
	}
	s.res.Success(http.StatusOK, "000", "", transactions)
	return &s.res
}

func (s *TransactionUseCase) GetAllTrx(limit, page string) *httputil.ResponseMessage {
	errorList := make(map[string]string)
	var user model.User
	result, err := s.client.Get("user_info").Result()
	if err != nil {
		errorList["message"] = err.Error()
		s.res.Errors(http.StatusBadRequest, "000", errorList)
		return &s.res
	}
	err = json.Unmarshal([]byte(result), &user)
	if err != nil {
		errorList["message"] = err.Error()
		s.res.Errors(http.StatusBadRequest, "000", errorList)
		return &s.res
	}
	userId := strconv.Itoa(int(user.Id))
	newLimit, errConv := strconv.Atoi(limit)
	if errConv != nil {
		newLimit = 10
	}
	newPage, errConv := strconv.Atoi(page)
	if errConv != nil {
		newPage = 1
	}
	pagination := pagination2.Pagination{
		Limit: newLimit,
		Page:  newPage,
	}

	transactions, err := s.repo.GetAllTrx(userId, pagination)
	if err != nil {
		errorList["message"] = err.Error()
		s.res.Errors(http.StatusBadRequest, "000", errorList)
		return &s.res
	}
	s.res.Success(http.StatusOK, "000", "", transactions)
	return &s.res
}

func (s *TransactionUseCase) GetById(id string) *httputil.ResponseMessage {
	errorList := make(map[string]string)
	var user model.User
	result, err := s.client.Get("user_info").Result()
	if err != nil {
		errorList["message"] = err.Error()
		s.res.Errors(http.StatusBadRequest, "000", errorList)
		return &s.res
	}
	err = json.Unmarshal([]byte(result), &user)
	if err != nil {
		errorList["message"] = err.Error()
		s.res.Errors(http.StatusBadRequest, "000", errorList)
		return &s.res
	}
	transaction, err := s.repo.GetOneById(id)
	if err != nil {
		errorList["message"] = err.Error()
		s.res.Errors(http.StatusBadRequest, "000", errorList)
		return &s.res
	}
	merchantId := strconv.Itoa(int(transaction.MerchantId))
	merchant, err := s.repoMerchant.GetOneById(merchantId)
	if err != nil {
		errorList["message"] = err.Error()
		log.Println(err.Error())
		s.res.Errors(http.StatusBadRequest, "000", errorList)
		return &s.res
	}
	if user.Id != merchant.UserId {
		errorList["message"] = "Cannot access"
		s.res.Errors(http.StatusBadRequest, "000", errorList)
		return &s.res
	}
	s.res.Success(http.StatusOK, "000", "", transaction)
	return &s.res

}

func NewTransactionUseCase(userRepository repository.ITransactionRepository, merchantRepository repository.IMerchantRepository, redisClient *redis.Client) ITransactionUseCase {
	return &TransactionUseCase{repo: userRepository, repoMerchant: merchantRepository, client: redisClient}
}
