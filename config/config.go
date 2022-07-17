package config

import (
	"fmt"
	docs "github.com/epulskyyy/majoo-test-2022/docs" // you need to update github.com/rizalgowandy/go-swag-sample with your own project path
	"github.com/epulskyyy/majoo-test-2022/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
)

type Config struct {
	RouterEngine   *gin.Engine
	DataSourceName string
	ApiBaseUrl     string
}

func NewConfig() *Config {

	docs.SwaggerInfo.Title = "Majoo Test API"
	docs.SwaggerInfo.Version = "1.0"

	config := new(Config)
	apiHost := os.Getenv("API_HOST")
	apiPort := os.Getenv("API_PORT")

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)
	fmt.Println(dsn)
	config.DataSourceName = dsn

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	config.RouterEngine = r

	config.ApiBaseUrl = fmt.Sprintf("%s:%s", apiHost, apiPort)
	return config
}
