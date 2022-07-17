package main

import (
	"github.com/epulskyyy/majoo-test-2022/api"
)

// @title           Majoo Test API
// @version         1.0
// @description     This is a sample Go API Documentation.

// @contact.name   Epul
// @contact.email  saepulstr@gmail.com

// @host      localhost:8080
// @BasePath  /api

// @query.collection.format multi


// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @x-extension-openapi {"example": "value on a json format"}
func main() {
	api.NewApiServer().Run()
}
