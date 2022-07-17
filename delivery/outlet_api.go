package delivery

import (
	"errors"
	"github.com/epulskyyy/majoo-test-2022/usecase"
	"github.com/gin-gonic/gin"
)

type OutletApi struct {
	usecase     usecase.IOutletUseCase
	publicRoute *gin.RouterGroup
}

func NewOutletApi(publicRoute *gin.RouterGroup, usecase usecase.IOutletUseCase) (*OutletApi, error) {
	if publicRoute == nil || usecase == nil {
		return nil, errors.New("Empty Router or UseCase")
	}

	studentApi := OutletApi{
		usecase:     usecase,
		publicRoute: publicRoute,
	}
	studentApi.InitRouter()
	return &studentApi, nil
}

func (api *OutletApi) InitRouter() {
	studentRoute := api.publicRoute.Group("/outlet")
	studentRoute.GET("/:id", api.getOutletById)
}


// getOutletById func for get outlet by id.
// @Description for get outlet by id.
// @Summary for get outlet by id
// @Tags Outlet
// @Accept json
// @Produce json
// @Success 201 {object} httputil.ResponseMessage{data=model.Outlet}
// @Failure 400 {object} httputil.ResponseMessage
// @Failure 404 {object} httputil.ResponseMessage
// @Security ApiKeyAuth
// @Router /outlet/:id [get]
func (api *OutletApi) getOutletById(c *gin.Context) {
	id := c.Param("id")
	res:= api.usecase.GetById(id)
	res.Send(c)
}
