package api

import (
	"github.com/epulskyyy/majoo-test-2022/config"
	"github.com/epulskyyy/majoo-test-2022/delivery"
	"github.com/epulskyyy/majoo-test-2022/manager"
	"log"
)

type Server interface {
	Run()
}

type server struct {
	config  *config.Config
	infra   manager.Infra
	usecase manager.UseCaseManager
}

func NewApiServer() Server {
	appConfig := config.NewConfig()
	infra := manager.NewInfra(appConfig)
	client := manager.NewClientManager(infra)
	repo := manager.NewRepoManager(infra)
	usecase := manager.NewUseCaseManger(repo, client)
	return &server{
		config:  appConfig,
		infra:   infra,
		usecase: usecase,
	}
}

func (s *server) Run() {
	err := delivery.NewServer(s.config.RouterEngine, s.usecase)
	err = s.config.RouterEngine.Run(s.config.ApiBaseUrl)
	if err != nil {
		log.Fatal(err)
	}
}
