package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/epulskyyy/majoo-test-2022/httputil"
	"github.com/epulskyyy/majoo-test-2022/model"
	"github.com/epulskyyy/majoo-test-2022/repository"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type IAuthUseCase interface {
	Login(auth model.Auth) *httputil.ResponseMessage
	ExtractTokenMetadata(r *http.Request) (*model.AccessDetails, error)
	FetchAuth(authD *model.AccessDetails) error
	Refresh(token string) *httputil.ResponseMessage
	Logout(c *http.Request) *httputil.ResponseMessage
}

type AuthUseCase struct {
	user        *model.User
	client      *redis.Client
	res         httputil.ResponseMessage
	repo        repository.IUserRepository
	tokenDetail model.TokenDetail
}

func (a AuthUseCase) Logout(r *http.Request) *httputil.ResponseMessage {
	errorList := make(map[string]string)
	metadata, err := a.ExtractTokenMetadata(r)
	if err != nil {
		errorList["message"] = err.Error()
		log.Println("[auth_usecase] error Extract Token ", err.Error())
		a.res.Errors(http.StatusBadRequest, "000", errorList)
		return &a.res
	}
	err = a.deleteTokens(metadata)
	if err != nil {
		errorList["message"] = err.Error()
		log.Println("[auth_usecase] error delete Token ", err.Error())
		a.res.Errors(http.StatusBadRequest, "000", errorList)
		return &a.res
	}
	a.res.Success(http.StatusOK, "000", "", "success logout")
	return &a.res
}

func (a AuthUseCase) Login(auth model.Auth) *httputil.ResponseMessage {
	errorList := make(map[string]string)
	var err error
	a.user, err = a.repo.GetOneByUserName(auth.Username)
	if err != nil {
		errorList["user_name"] = "Username not registered"
		log.Println("[auth_usecase] error Username not registered")
		a.res.Errors(http.StatusNotFound, "000", errorList)
		return &a.res
	}
	if errCompare := bcrypt.CompareHashAndPassword([]byte(a.user.Password), []byte(auth.Password)); errCompare != nil {
		// If the two passwords don't match, return a 401 status
		errorList["password"] = "Password wrong"
		log.Println("[auth_usecase] error Password wrong")
		a.res.Errors(http.StatusUnauthorized, "000", errorList)
		return &a.res
	}
	errCreateToken := a.createToken()
	if errCreateToken != nil {
		errorList["message"] = errCreateToken.Error()
		log.Println("[auth_usecase] error CreateToken ", errCreateToken.Error())
		a.res.Errors(http.StatusUnprocessableEntity, "000", errorList)
		return &a.res
	}
	errCreateAuth := a.createAuth()
	if errCreateAuth != nil {
		errorList["message"] = errCreateAuth.Error()
		log.Println("[auth_usecase] error CreateAuth ", errCreateAuth.Error())
		a.res.Errors(http.StatusUnprocessableEntity, "000", errorList)
		return &a.res
	}
	result := &model.AuthResp{
		AccessToken:  a.tokenDetail.AccessToken,
		RefreshToken: a.tokenDetail.RefreshToken,
	}
	a.res.Success(http.StatusOK, "000", "", result)
	return &a.res
}

func (a AuthUseCase) ExtractTokenMetadata(r *http.Request) (*model.AccessDetails, error) {
	token, err := a.verifyToken(r)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		id,_:= claims["id"].(float64)
		return &model.AccessDetails{
			AccessUuid: accessUuid,
			Id: uint(id)      ,
			Username:   claims["user_name"].(string),
		}, nil
	}
	return nil, err
}

func (a AuthUseCase) TokenValid(r *http.Request) error {
	panic("implement me")
}

func (a AuthUseCase) FetchAuth(authD *model.AccessDetails) error {
	var err error
	a.user, err = a.repo.GetOneById(strconv.Itoa(int(authD.Id)))
	if err != nil {
		return errors.New("User not found")
	}
	userInfo, err := json.Marshal(a.user)
	if err != nil {
		log.Println("ERROR Marshal ", err)
		return err
	}

	errSetUser := a.client.Set("user_info",string(userInfo),0)
	if errSetUser != nil {
		return errSetUser.Err()
	}
	errSetUserId := a.client.Set("user_id",a.user.Id,0)
	if errSetUserId != nil {
		return errSetUserId.Err()
	}
	if authD.Username != a.user.UserName {
		return errors.New("unauthorized")
	}
	return nil
}

func (a AuthUseCase) Refresh(refreshToken string) *httputil.ResponseMessage {
	errorList := make(map[string]string)
	//verify the token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println(fmt.Errorf("unexpected signing method: %v", token.Header["alg"]))
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	//if there is an error, the token must
	//have expired
	if err != nil {
		errorList["message"] = err.Error()
		log.Println("[auth_usecase] token expired")
		a.res.Errors(http.StatusUnauthorized, "000", errorList)
		return &a.res
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		errorList["message"] = "token not valid"
		log.Println("[auth_usecase] error token not valid")
		a.res.Errors(http.StatusUnauthorized, "000", errorList)
		return &a.res
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, okRef := claims["refresh_uuid"].(string) //convert the interface to string
		if !okRef {
			errorList["message"] = "unauthorized"
			log.Println("[auth_usecase] error unauthorized")
			a.res.Errors(http.StatusUnauthorized, "000", errorList)
			return &a.res
		}
		//Delete the previous Refresh Token
		deleted, errDel := a.client.Del(refreshUuid).Result()
		if errDel != nil || deleted == 0 { //if any goes wrong
			errorList["message"] = "unauthorized"
			log.Println("[auth_usecase] error Delete Auth")
			a.res.Errors(http.StatusUnauthorized, "000", errorList)
			return &a.res
		}
		id,_:= claims["id"].(float64)
		a.user.Id = uint(id)
		a.user.UserName = claims["user_name"].(string)

		//Create new pairs of refresh and access tokens
		errCreateToken := a.createToken()
		if errCreateToken != nil {
			errorList["message"] = errCreateToken.Error()
			log.Println("[auth_usecase] error CreateToken ", errCreateToken.Error())
			a.res.Errors(http.StatusUnprocessableEntity, "000", errorList)
			return &a.res
		}
		//save the tokens metadata to redis
		errCreateAuth := a.createAuth()
		if errCreateAuth != nil {
			errorList["message"] = errCreateAuth.Error()
			log.Println("[auth_usecase] error CreateAuth ", errCreateAuth.Error())
			a.res.Errors(http.StatusUnprocessableEntity, "000", errorList)
			return &a.res
		}
		result := &model.AuthResp{
			AccessToken:  a.tokenDetail.AccessToken,
			RefreshToken: a.tokenDetail.RefreshToken,
		}
		a.res.Success(http.StatusOK, "000", "", result)
		return &a.res
	} else {
		errorList["message"] = "unauthorized"
		log.Println("[auth_usecase] unauthorized")
		a.res.Errors(http.StatusUnprocessableEntity, "000", errorList)
		return &a.res
	}
}

func NewAuthUseCase(userRepository repository.IUserRepository, redisClient *redis.Client) IAuthUseCase {
	return &AuthUseCase{repo: userRepository, client: redisClient}
}

func (a *AuthUseCase) createToken() error {
	sessionExpired := os.Getenv("AUTH_EXPIRED")
	if sessionExpired == "" {
		sessionExpired = "60"
	}

	sessionExpiredInt, _ := strconv.Atoi(sessionExpired)

	a.tokenDetail.AtExpires = time.Now().Add(time.Minute * time.Duration(sessionExpiredInt)).Unix()
	a.tokenDetail.AccessUuid = uuid.New().String()

	a.tokenDetail.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	a.tokenDetail.RefreshUuid = a.tokenDetail.AccessUuid + "++" + strconv.Itoa(int(a.user.Id))

	var err error

	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = a.tokenDetail.AccessUuid
	atClaims["id"] = a.user.Id
	atClaims["user_name"] = a.user.UserName
	atClaims["exp"] = a.tokenDetail.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	a.tokenDetail.AccessToken, err = at.SignedString([]byte(os.Getenv("AUTH_ACCESS_SECRET")))
	if err != nil {
		return err
	}

	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = a.tokenDetail.RefreshUuid
	rtClaims["id"] = a.user.Id
	rtClaims["user_name"] = a.user.UserName
	rtClaims["exp"] = a.tokenDetail.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)

	a.tokenDetail.RefreshToken, err = rt.SignedString([]byte(os.Getenv("AUTH_REFRESH_SECRET")))
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthUseCase) createAuth() error {
	//converting Unix to UTC(to Time object)
	at := time.Unix(a.tokenDetail.AtExpires, 0)
	rt := time.Unix(a.tokenDetail.RtExpires, 0)

	now := time.Now()

	authJson, err := json.Marshal(&model.AuthReq{Id: a.user.Id, Username: a.user.UserName})
	if err != nil {
		log.Println("ERROR Marshal ", err)
		return err
	}
	errAccess := a.client.Set(a.tokenDetail.AccessUuid, string(authJson), at.Sub(now)).Err()
	if errAccess != nil {
		log.Println("ERROR set Access UUID ", err)
		return errAccess
	}
	errRefresh := a.client.Set(a.tokenDetail.RefreshUuid, string(authJson), rt.Sub(now)).Err()
	if errRefresh != nil {
		log.Println("ERROR set Refresh UUID ", err)
		return errRefresh
	}
	return nil
}

func (a AuthUseCase) verifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := a.extractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("AUTH_ACCESS_SECRET")), nil
	})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return token, nil
}

func (a AuthUseCase) extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func (a AuthUseCase) deleteTokens(details *model.AccessDetails) error {
	refreshUuid := fmt.Sprintf("%s++%s", details.AccessUuid, details.Id)
	log.Println(details.AccessUuid, refreshUuid)
	//delete access token
	deletedAt, err := a.client.Del(details.AccessUuid).Result()
	if err != nil {
		return err
	}
	//delete refresh token
	deletedRt, err := a.client.Del(refreshUuid).Result()
	if err != nil {
		return err
	}
	//When the record is deleted, the return value is 1
	if deletedAt != 1 || deletedRt != 1 {
		return errors.New("something went wrong")
	}
	return nil
}
