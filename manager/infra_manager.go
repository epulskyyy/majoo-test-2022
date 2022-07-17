package manager

import (
	"fmt"
	"github.com/epulskyyy/majoo-test-2022/config"
	"github.com/epulskyyy/majoo-test-2022/model"
	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

type Infra interface {
	GormDB() *gorm.DB
	RedisClient() *redis.Client
}

type infra struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewInfra(config *config.Config) Infra {
	fmt.Println(config.DataSourceName)
	db, redisClient, err := initDbResource(config.DataSourceName)
	if err != nil {
		log.Panicln(err)
	}
	return &infra{
		db:          db,
		redisClient: redisClient,
	}
}

func (i *infra) GormDB() *gorm.DB {
	return i.db
}
func (i *infra) RedisClient() *redis.Client {
	return i.redisClient
}

func initDbResource(dataSourceName string) (*gorm.DB, *redis.Client, error) {
	db, errOpen := gorm.Open(mysql.Open(dataSourceName))
	if errOpen != nil {
		return nil, nil, errOpen
	}

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Merchant{})
	db.AutoMigrate(&model.Outlet{})
	db.AutoMigrate(&model.Transaction{})

	redisClient := redis.NewClient(&redis.Options{
		Addr:         os.Getenv("REDIS_URL"),
		MinIdleConns: 20,
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Panic(err)
		fmt.Print(err.Error())
	}

	return db, redisClient, nil
}
