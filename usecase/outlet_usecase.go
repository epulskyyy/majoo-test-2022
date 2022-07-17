package usecase

import (
	"encoding/json"
	"github.com/epulskyyy/majoo-test-2022/httputil"
	"github.com/epulskyyy/majoo-test-2022/model"
	"github.com/epulskyyy/majoo-test-2022/repository"
	"github.com/go-redis/redis"
	"net/http"
)

type IOutletUseCase interface {
	GetById(id string) *httputil.ResponseMessage
}

type OutletUseCase struct {
	repo   repository.IOutletRepository
	client *redis.Client
	res    httputil.ResponseMessage
}

func (s *OutletUseCase) GetById(id string) *httputil.ResponseMessage {
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
	outlet, err := s.repo.GetOneById(id)
	if err != nil {
		errorList["message"] = err.Error()
		s.res.Errors(http.StatusBadRequest, "000", errorList)
		return &s.res
	}

	if user.Id != outlet.Merchant.UserId {
		errorList["message"] = "Cannot access"
		s.res.Errors(http.StatusBadRequest, "000", errorList)
		return &s.res
	}

	s.res.Success(http.StatusOK, "000", "", outlet)
	return &s.res

}

func NewOutletUseCase(userRepository repository.IOutletRepository, redisClient *redis.Client) IOutletUseCase {
	return &OutletUseCase{repo:userRepository, client: redisClient}
}
