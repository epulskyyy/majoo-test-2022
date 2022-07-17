package delivery

import (
	"errors"
	"github.com/epulskyyy/majoo-test-2022/usecase"
	"github.com/gin-gonic/gin"
)

type UserApi struct {
	usecase     usecase.IUserUseCase
	publicRoute *gin.RouterGroup
}

func NewUserApi(publicRoute *gin.RouterGroup, usecase usecase.IUserUseCase) (*UserApi, error) {
	if publicRoute == nil || usecase == nil {
		return nil, errors.New("Empty Router or UseCase")
	}

	studentApi := UserApi{
		usecase:     usecase,
		publicRoute: publicRoute,
	}
	studentApi.InitRouter()
	return &studentApi, nil
}

func (api *UserApi) InitRouter() {
	studentRoute := api.publicRoute.Group("/user")
	studentRoute.GET("/info", api.getUserInfo)
}


// getUserInfo func for get info.
// @Description for get info.
// @Summary for get info
// @Tags User
// @Accept json
// @Produce json
// @Success 201 {object} httputil.ResponseMessage{data=model.User}
// @Failure 400 {object} httputil.ResponseMessage
// @Failure 404 {object} httputil.ResponseMessage
// @Security ApiKeyAuth
// @Router /user/info [get]
func (api *UserApi) getUserInfo(c *gin.Context) {
	res:= api.usecase.UserInfo()
	res.Send(c)
}
