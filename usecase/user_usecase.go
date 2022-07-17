package usecase

import (
	"encoding/json"
	"github.com/epulskyyy/majoo-test-2022/httputil"
	"github.com/epulskyyy/majoo-test-2022/model"
	"github.com/epulskyyy/majoo-test-2022/repository"
	"github.com/go-redis/redis"
	"net/http"
)

type IUserUseCase interface {
	UserInfo() *httputil.ResponseMessage
}

type UserUseCase struct {
	repo   repository.IUserRepository
	client *redis.Client
	res    httputil.ResponseMessage
}

func (s *UserUseCase) UserInfo() *httputil.ResponseMessage {
	var user model.User
	errorList := make(map[string]string)
	result, err := s.client.Get("user_info").Result()
	if err != nil {
		errorList["message"] = err.Error()
		s.res.Errors(http.StatusBadRequest, "000", errorList)
		return &s.res
	}
	err = json.Unmarshal([]byte(result), &user)
	s.res.Success(http.StatusOK, "000", "", user)
	return &s.res

}

func NewUserUseCase(userRepository repository.IUserRepository, redisClient *redis.Client) IUserUseCase {
	return &UserUseCase{repo:userRepository, client: redisClient}
}
