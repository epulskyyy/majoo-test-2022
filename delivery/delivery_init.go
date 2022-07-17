package delivery

import (
	"github.com/epulskyyy/majoo-test-2022/manager"
	"github.com/epulskyyy/majoo-test-2022/middleware"
	"github.com/gin-gonic/gin"
)

func NewServer(engine *gin.Engine, useCaseManager manager.UseCaseManager) error {
	publicRoute := engine.Group("/api")
	_, err := NewAuthApi(publicRoute, useCaseManager.AuthUseCase())
	publicRoute.Use(middleware.JWTMiddleware(useCaseManager.AuthUseCase()))

	_, err = NewUserApi(publicRoute, useCaseManager.UserUseCase())
	_, err = NewMerchantApi(publicRoute, useCaseManager.MerchantUseCase())
	_, err = NewOutletApi(publicRoute, useCaseManager.OutletUseCase())
	_, err = NewTransactionApi(publicRoute, useCaseManager.TransactionUseCase())
	return err
}
