package delivery

import (
	"errors"
	"github.com/epulskyyy/majoo-test-2022/httputil"
	"github.com/epulskyyy/majoo-test-2022/model"
	"github.com/epulskyyy/majoo-test-2022/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthApi struct {
	usecase     usecase.IAuthUseCase
	publicRoute *gin.RouterGroup
}

func NewAuthApi(publicRoute *gin.RouterGroup, usecase usecase.IAuthUseCase) (*AuthApi, error) {
	if publicRoute == nil || usecase == nil {
		return nil, errors.New("Empty Router or UseCase")
	}
	studentApi := AuthApi{
		usecase:     usecase,
		publicRoute: publicRoute,
	}
	studentApi.InitRouter()
	return &studentApi, nil
}

func (api *AuthApi) InitRouter() {
	studentRoute := api.publicRoute.Group("/")
	studentRoute.POST("/login", api.login)
	studentRoute.POST("/logout", api.logout)
	studentRoute.POST("/refresh", api.refresh)
}

// Login func Login.
// @Description Login.
// @Summary Login
// @Tags Auth
// @Accept json
// @Produce json
// @Param account_req_attrs body model.Auth true "Auth attributes"
// @Success 201 {object} httputil.ResponseMessage{data=model.AuthResp}
// @Failure 400 {object} httputil.ResponseMessage
// @Failure 404 {object} httputil.ResponseMessage
// @Router /login [post]
func (api *AuthApi) login(context *gin.Context) {
	var auth model.Auth
	res := auth.ValidateRequest(context)
	if res != nil {
		res.Send(context)
		return
	}
	res = api.usecase.Login(auth)
	res.Send(context)
}


// Logout func Logout.
// @Description Logout.
// @Summary Logout
// @Tags Auth
// @Accept json
// @Produce json
// @Success 201 {object} httputil.ResponseMessage
// @Failure 400 {object} httputil.ResponseMessage
// @Failure 404 {object} httputil.ResponseMessage
// @Security ApiKeyAuth
// @Router /logout [post]
func (api *AuthApi) logout(context *gin.Context) {
	res:= api.usecase.Logout(context.Request)
	res.Send(context)
}

// Refresh func Refresh token.
// @Description Refresh token.
// @Summary Refresh
// @Tags Auth
// @Accept json
// @Produce json
// @Param account_req_attrs body model.AuthToken true "AuthToken attributes"
// @Success 201 {object} httputil.ResponseMessage{data=model.AuthResp}
// @Failure 400 {object} httputil.ResponseMessage
// @Failure 404 {object} httputil.ResponseMessage
// @Router /refresh [post]
func (api *AuthApi) refresh(context *gin.Context) {
	var authToken model.AuthToken
	var res *httputil.ResponseMessage
	errorList 	:= make(map[string]string)
	err := context.ShouldBind(&authToken)
	if err != nil {
		errorList["message"] = err.Error()
		res.Errors(http.StatusBadRequest,  "000", errorList)
	}
	res = api.usecase.Refresh(authToken.Token)
	res.Send(context)
}
