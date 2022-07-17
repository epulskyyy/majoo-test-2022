package manager

import (
	"github.com/go-redis/redis"
)

type ClientManager interface {
	Redis() *redis.Client
}

type clientManager struct {
	infra Infra
}

func (c clientManager) Redis() *redis.Client {
	return c.infra.RedisClient()
}

func NewClientManager(infra Infra) ClientManager {
	return &clientManager{infra}
}
