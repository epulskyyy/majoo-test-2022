package delivery

import (
	"github.com/epulskyyy/majoo-test-2022/manager"
	"github.com/epulskyyy/majoo-test-2022/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type useCaseManagerMock struct {
	mock.Mock
}

func (uc *useCaseManagerMock) AuthUseCase() usecase.IAuthUseCase {
	panic("implement me")
}

func (uc *useCaseManagerMock) UserUseCase() usecase.IUserUseCase {
	args := uc.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(usecase.IUserUseCase)
}

type DeliveryInitTestSuite struct {
	suite.Suite
	routerTest         *gin.Engine
	useCaseManagerTest manager.UseCaseManager
}

func (suite *DeliveryInitTestSuite) SetupTest() {
	suite.routerTest = gin.Default()
	suite.useCaseManagerTest = new(useCaseManagerMock)
}
func (suite *DeliveryInitTestSuite) TestDeliveryInit_CreateServer_Success() {
	suite.useCaseManagerTest.(*useCaseManagerMock).On("UserUseCase").Return(usecase.NewUserUseCase(nil))
	err := NewServer(suite.routerTest, suite.useCaseManagerTest)
	assert.Nil(suite.T(), err)
}
func (suite *DeliveryInitTestSuite) TestDeliveryInit_CreateServer_Failed() {
	suite.useCaseManagerTest.(*useCaseManagerMock).On("UserUseCase").Return(nil)
	err := NewServer(suite.routerTest, suite.useCaseManagerTest)
	assert.NotNil(suite.T(), err)
}

func TestDeliveryInitTestSuite(t *testing.T) {
	suite.Run(t, new(DeliveryInitTestSuite))
}
