package usecase

import (
	"encoding/json"
	"github.com/epulskyyy/majoo-test-2022/httputil"
	"github.com/epulskyyy/majoo-test-2022/model"
	"github.com/epulskyyy/majoo-test-2022/repository"
	"github.com/go-redis/redis"
	"log"
	"net/http"
)

type IMerchantUseCase interface {
	GetById(id string) *httputil.ResponseMessage
}

type MerchantUseCase struct {
	repo   repository.IMerchantRepository
	client *redis.Client
	res    httputil.ResponseMessage
}

func (s *MerchantUseCase) GetById(id string) *httputil.ResponseMessage {
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
	merchant, err := s.repo.GetOneById(id)
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

	s.res.Success(http.StatusOK, "000", "", merchant)
	return &s.res

}

func NewMerchantUseCase(userRepository repository.IMerchantRepository, redisClient *redis.Client) IMerchantUseCase {
	return &MerchantUseCase{repo:userRepository, client: redisClient}
}
