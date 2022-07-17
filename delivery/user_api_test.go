package delivery

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/epulskyyy/majoo-test-2022/model"
	"github.com/epulskyyy/majoo-test-2022/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var dummyUsers = []model.User{
	{
		ID:        big.NewInt(1),
		Name:      "Epul",
		UserName:  "Epul1",
		CreatedAt: time.Time{},
		CreatedBy: big.NewInt(1),
		UpdatedAt: time.Time{},
		UpdatedBy: big.NewInt(1),
	},
	{
		ID:        big.NewInt(1),
		Name:      "Ani",
		UserName:  "Epul1",
		CreatedAt: time.Time{},
		CreatedBy: big.NewInt(1),
		UpdatedAt: time.Time{},
		UpdatedBy: big.NewInt(1),
	},
}

type MockResponse struct {
	Message model.User
}
type MockErrorResponse struct {
	Message string
}

type userUseCaseMock struct {
	mock.Mock
}


/*
Untuk HTTP client test
http://hassansin.github.io/Unit-Testing-http-client-in-Go
*/

func (s *userUseCaseMock) FindAllUser() ([]model.User, error) {
	panic("implement me")
}

func (s *userUseCaseMock) NewRegistration(student model.User) (*model.User, error) {
	args := s.Called(student)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), nil
}

func (s *userUseCaseMock) FindUserInfoById(idCard string) (*model.User, error) {
	args := s.Called(idCard)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), nil
}

type UserApiTestSuite struct {
	suite.Suite
	useCaseTest     usecase.IUserUseCase
	routerTest      *gin.Engine
	routerGroupTest *gin.RouterGroup
}

func (suite *UserApiTestSuite) SetupTest() {
	suite.useCaseTest = new(userUseCaseMock)
	suite.routerTest = gin.Default()
	suite.routerGroupTest = suite.routerTest.Group("/api")
}
func (suite *UserApiTestSuite) Test_CreateNewUserAPI_Success() {
	student, err := NewUserApi(suite.routerGroupTest, suite.useCaseTest)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), student)

}
func (suite *UserApiTestSuite) TestUserApi_CreateUser_Success() {
	dummyUser := dummyUsers[1]
	suite.useCaseTest.(*userUseCaseMock).On("NewRegistration", dummyUser).Return(&dummyUser, nil)
	studentApi, _ := NewUserApi(suite.routerGroupTest, suite.useCaseTest)
	handler := studentApi.createUser
	suite.routerTest.POST("", handler)

	rr := httptest.NewRecorder()
	reqBody, _ := json.Marshal(dummyUser)
	request, _ := http.NewRequest(http.MethodPost, "/api/student", bytes.NewBuffer(reqBody))
	request.Header.Set("Content-Type", "application/json")

	suite.routerTest.ServeHTTP(rr, request)
	assert.Equal(suite.T(), rr.Code, 200)

	//expectedRespBody, _ := json.Marshal(gin.H{
	//	"message": dummyUser,
	//})
	//assert.Equal(suite.T(), expectedRespBody, rr.Body.Bytes())
	a := rr.Body.String()
	actualUser := new(MockResponse)
	json.Unmarshal([]byte(a), actualUser)
	assert.Equal(suite.T(), dummyUser.Name, actualUser.Message.Name)
}
func (suite *UserApiTestSuite) TestUserApi_CreateUser_FailedBinding() {
	suite.useCaseTest.(*userUseCaseMock).On("NewRegistration", nil).Return(nil, errors.New("failed"))
	studentApi, _ := NewUserApi(suite.routerGroupTest, suite.useCaseTest)
	handler := studentApi.createUser
	suite.routerTest.POST("", handler)

	rr := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/api/student", nil)
	request.Header.Set("Content-Type", "application/json")

	suite.routerTest.ServeHTTP(rr, request)
	assert.Equal(suite.T(), rr.Code, 400)
}

func (suite *UserApiTestSuite) TestUserApi_CreateUser_FailedUseCase() {
	dummyUser := dummyUsers[1]
	suite.useCaseTest.(*userUseCaseMock).On("NewRegistration", dummyUser).Return(nil, errors.New("failed"))
	studentApi, _ := NewUserApi(suite.routerGroupTest, suite.useCaseTest)
	handler := studentApi.createUser
	suite.routerTest.POST("", handler)

	rr := httptest.NewRecorder()
	reqBody, _ := json.Marshal(dummyUser)
	request, _ := http.NewRequest(http.MethodPost, "/api/student", bytes.NewBuffer(reqBody))
	request.Header.Set("Content-Type", "application/json")

	suite.routerTest.ServeHTTP(rr, request)
	assert.Equal(suite.T(), rr.Code, 500)
	a := rr.Body.String()
	actualError := new(MockErrorResponse)
	json.Unmarshal([]byte(a), actualError)
	assert.Equal(suite.T(), "failed", actualError.Message)
}

func (suite *UserApiTestSuite) TestUserApi_GetById_Success() {
	dummyUser := dummyUsers[0]
	suite.useCaseTest.(*userUseCaseMock).On("FindUserInfoById", "2").Return(&dummyUser, nil)
	studentApi, _ := NewUserApi(suite.routerGroupTest, suite.useCaseTest)
	handler := studentApi.getUserById
	suite.routerTest.GET("/:idcard", handler)

	rr := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/api/student/2", nil)
	suite.routerTest.ServeHTTP(rr, request)
	assert.Equal(suite.T(), rr.Code, 200)

	a := rr.Body.String()
	actualUser := new(MockResponse)
	json.Unmarshal([]byte(a), actualUser)
	assert.Equal(suite.T(), dummyUser.Name, actualUser.Message.Name)
}
func (suite *UserApiTestSuite) TestUserApi_GetById_Failed() {
	suite.useCaseTest.(*userUseCaseMock).On("FindUserInfoById", "1").Return(nil, errors.New("failed"))
	studentApi, _ := NewUserApi(suite.routerGroupTest, suite.useCaseTest)
	handler := studentApi.getUserById
	suite.routerTest.GET("/:idcard", handler)

	rr := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/api/student/1", nil)
	request.Header.Set("Content-Type", "application/json")

	suite.routerTest.ServeHTTP(rr, request)
	assert.Equal(suite.T(), rr.Code, 400)
}
func TestUserApiTestSuite(t *testing.T) {
	suite.Run(t, new(UserApiTestSuite))
}
