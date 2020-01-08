package redis

import (
	"fmt"
	"sync"

	"demo/pkg/config"

	"github.com/go-redis/redis"
)

var once sync.Once
var client *redis.Client

func Init() {
	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     config.GetString("redis.client.addr"),
			Password: config.GetString("redis.client.password"), // no password set
			DB:       config.GetInt("redis.client.db"),          // use default DB
			PoolSize: config.GetInt("redis.client.pool"),
		})

		pong, err := client.Ping().Result()
		if err == nil {
			fmt.Printf("\033[1;30;42m[info]\033[0m redis [riskcenter] connect success %s\n", pong)
		} else {
			panic(fmt.Sprintf("\033[1;30;41m[error]\033[0m redis connect error %s\n", err.Error()))
		}
	})
}

// RiskCenterClient 获取风控中心 Redis 连接
func GetClient() *redis.Client {
	return client
}
