package httputil

import (
	"github.com/gin-gonic/gin"
)

type ResponseMessage struct {
	Status      string      `json:"status"`
	Code        string      `json:"code"`
	HttpCode    int      	`json:"http_code,omitempty"`
	Description string      `json:"message,omitempty"`
	Data        interface{} `json:"data,omitempty"`
	Error 		interface{} `json:"error,omitempty"`
}

func (r *ResponseMessage) sendData(httpCode int, code, desk string, data interface{})  {
	r.Status = "Success"
	r.HttpCode = httpCode
	r.Code = code
	r.Description = desk
	r.Data = data
}

func (r *ResponseMessage) sendError(httpCode int, code string, err interface{})  {
	r.Status = "Failure"
	r.HttpCode = httpCode
	r.Code = code
	r.Data = nil
	r.Error = err
}

func (r *ResponseMessage) Errors(httpCode int, code string, error interface{}) {
	r.sendError(httpCode, code, error)
}
func (r *ResponseMessage) Success(httpCode int, code, desk string, data interface{}) {
	r.sendData(httpCode, code, desk, data)
}
func (r ResponseMessage) Send(context *gin.Context)  {
	context.JSON(r.HttpCode, r)
	return
}