package cache

import (
	"context"
	"github/SmoothieNoIce/dcard-backend-intern-assignment-2024/config"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var Rdb0 *redis.Client

type redisCache struct {
	client redis.Client
}

func SetUpDefaultDB() {
	Rdb0 = redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.Redis.Host,
		Password: config.AppConfig.Redis.Password,
		DB:       0,
	})
}
