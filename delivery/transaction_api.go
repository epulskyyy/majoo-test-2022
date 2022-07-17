package delivery

import (
	"errors"
	"github.com/epulskyyy/majoo-test-2022/usecase"
	"github.com/gin-gonic/gin"
)

type TransactionApi struct {
	usecase     usecase.ITransactionUseCase
	publicRoute *gin.RouterGroup
}

func NewTransactionApi(publicRoute *gin.RouterGroup, usecase usecase.ITransactionUseCase) (*TransactionApi, error) {
	if publicRoute == nil || usecase == nil {
		return nil, errors.New("Empty Router or UseCase")
	}

	studentApi := TransactionApi{
		usecase:     usecase,
		publicRoute: publicRoute,
	}
	studentApi.InitRouter()
	return &studentApi, nil
}

func (api *TransactionApi) InitRouter() {
	studentRoute := api.publicRoute.Group("/transaction")
	studentRoute.GET("/:id", api.getTransactionById)
	studentRoute.GET("/", api.getTransactions)
	studentRoute.GET("/report/", api.getReportTrx)
	studentRoute.GET("/report/csv", api.createCSV)
}


// getTransactionById func for get transaction by id.
// @Description for get transaction by id.
// @Summary for get transaction by id
// @Tags Transaction
// @Accept json
// @Produce json
// @Success 201 {object} httputil.ResponseMessage{data=model.Transaction}
// @Failure 400 {object} httputil.ResponseMessage
// @Failure 404 {object} httputil.ResponseMessage
// @Security ApiKeyAuth
// @Router /api/transaction/:id [get]
func (api *TransactionApi) getTransactionById(c *gin.Context) {
	id := c.Param("id")
	res:= api.usecase.GetById(id)
	res.Send(c)
}

// getTransactions func for get transactions.
// @Description for get transactions.
// @Summary for get transactions
// @Tags Transaction
// @Accept json
// @Produce json
// @Param page query string false "page"
// @Param limit query string false "limit"
// @Success 201 {object} httputil.ResponseMessage{data=pagination.Pagination}
// @Failure 400 {object} httputil.ResponseMessage
// @Failure 404 {object} httputil.ResponseMessage
// @Security ApiKeyAuth
// @Router /transaction/:id [get]
func (api *TransactionApi) getTransactions(c *gin.Context) {
	page := c.Query("page")
	limit := c.Query("limit")
	res:= api.usecase.GetAllTrx(limit,page)
	res.Send(c)
}
// getReportTrxs func create report transaction.
// @Description create report transaction.
// @Summary create report transaction
// @Tags Transaction
// @Accept json
// @Produce json
// @Param limit query string false "limit"
// @Param page query string false "page"
// @Param year query string false "year"
// @Param month query string false "month"
// @Security ApiKeyAuth
// @Router /transaction/report/ [get]
func (api *TransactionApi) getReportTrx(c *gin.Context) {
	page := c.Query("page")
	limit := c.Query("limit")
	month := c.Query("month")
	year := c.Query("year")
	res:= api.usecase.GetReportTrx(limit,page,month,year)
	res.Send(c)
}

// createCSV func create report transaction.
// @Description create report transaction.
// @Summary create report transaction
// @Tags Transaction
// @Accept json
// @Produce json
// @Param year query string false "year"
// @Param month query string false "month"
// @Security ApiKeyAuth
// @Router /transaction/report/csv [get]
func (api *TransactionApi) createCSV(c *gin.Context) {
	month := c.Query("month")
	year := c.Query("year")
	api.usecase.CreateCSV(month,year,c)
}

