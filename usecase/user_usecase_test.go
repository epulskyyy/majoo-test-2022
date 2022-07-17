package usecase

import (
	"github.com/epulskyyy/majoo-test-2022/model"
	"github.com/epulskyyy/majoo-test-2022/repository"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

var dummyUsers = []model.User{
}

type repoMock struct {
	mock.Mock
}

func (r *repoMock) GetOneByUserName(username string) (*model.User, error) {
	args := r.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), nil
}

func (r *repoMock) GetAll() ([]model.User, error) {
	args := r.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.User), nil
}

func (r *repoMock) GetOneById(idCard string) (*model.User, error) {
	args := r.Called(idCard)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), nil
}

func (r *repoMock) CreateOne(student model.User) (*model.User, error) {
	args := r.Called(student)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), nil
}

type UserUseCaseTestSuite struct {
	suite.Suite
	repoTest repository.IUserRepository
}

//Dipanggil setiap kali menjalankan test
func (suite *UserUseCaseTestSuite) SetupTest() {
	suite.repoTest = new(repoMock)
}

//Semua method yang diawali kata Test, dianggap sebagai unit test di test suite
func TestUserUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserUseCaseTestSuite))
}
