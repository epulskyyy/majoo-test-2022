package delivery

import (
	"errors"
	"github.com/epulskyyy/majoo-test-2022/usecase"
	"github.com/gin-gonic/gin"
)

type MerchantApi struct {
	usecase     usecase.IMerchantUseCase
	publicRoute *gin.RouterGroup
}

func NewMerchantApi(publicRoute *gin.RouterGroup, usecase usecase.IMerchantUseCase) (*MerchantApi, error) {
	if publicRoute == nil || usecase == nil {
		return nil, errors.New("Empty Router or UseCase")
	}

	studentApi := MerchantApi{
		usecase:     usecase,
		publicRoute: publicRoute,
	}
	studentApi.InitRouter()
	return &studentApi, nil
}

func (api *MerchantApi) InitRouter() {
	studentRoute := api.publicRoute.Group("/merchant")
	studentRoute.GET("/:id", api.getMerchantById)
}


// getMerchantById func for get merchant by id.
// @Description for get merchant by id.
// @Summary for get merchant by id
// @Tags Merchant
// @Accept json
// @Produce json
// @Success 201 {object} httputil.ResponseMessage{data=model.Merchant}
// @Failure 400 {object} httputil.ResponseMessage
// @Failure 404 {object} httputil.ResponseMessage
// @Security ApiKeyAuth
// @Router /merchant/:id [get]
func (api *MerchantApi) getMerchantById(c *gin.Context) {
	id := c.Param("id")
	res:= api.usecase.GetById(id)
	res.Send(c)
}
