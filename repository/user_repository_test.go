package repository

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/epulskyyy/majoo-test-2022/model"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"testing"
)

var dummyUsers = []model.User{
}

type UserRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	mock sqlmock.Sqlmock
}

type Student struct {
	ID   string
	Name string
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	var (
		db  *sql.DB
		err error
	)
	db, suite.mock, err = sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	if db == nil {
		log.Fatalf("mock db is null")
	}
	if suite.mock == nil {
		log.Fatalf("sqlmock is null")
	}
	suite.db, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		panic(err) // Error here
	}
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
