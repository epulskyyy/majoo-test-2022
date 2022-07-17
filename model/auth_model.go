package model

import (
	"encoding/json"
	"github.com/epulskyyy/majoo-test-2022/httputil"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"log"
	"net/http"
	"strings"
)

type Auth struct {
	Username string `json:"user_name"`
	Password string `json:"password"`
}
type AccessDetails struct {
	AccessUuid string   `json:"access_uuid"`
	Id         uint `json:"id"`
	Username   string   `json:"user_name"`
}
type AuthReq struct {
	Id       uint `json:"id"`
	Username string   `json:"user_name"`
}
type TokenDetail struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}
type AuthResp struct {
	AccessToken  string `json:"access_token" `
	RefreshToken string `json:"refresh_token"`
}
type AuthToken struct {
	Token string `json:"token"`
}

func (a *Auth) ValidateRequest(context *gin.Context) *httputil.ResponseMessage {
	log.Println("VALIDATE REQUEST")
	var res httputil.ResponseMessage
	errorList := make(map[string]string)
	if err := context.ShouldBind(&a); err != nil {
		authJson, _ := json.Marshal(&a)
		errorList["message"] = "the body of the request is not valid JSON; example = " + string(authJson)
		log.Println("the body of the request is not valid JSON", err.Error())
		res.Errors(http.StatusBadRequest, "000", errorList)
		return &res
	}
	err := validation.Errors{
		"password": validation.Validate(a.Password, validation.Required),
		"user_name": validation.Validate(a.Username, validation.Required),
	}.Filter()
	if err != nil {
		errArr := strings.Split(err.Error(), ";")
		for i := 0; i < len(errArr); i++ {
			errSp := strings.Split(errArr[i], ":")
			key := strings.TrimSpace(errSp[0])
			value := strings.TrimSpace(errSp[1])
			errorList[key] = value
		}
	}
	if len(errorList) > 0 {
		res.Errors(http.StatusBadRequest, "000", errorList)
		return &res
	}
	return nil
}
