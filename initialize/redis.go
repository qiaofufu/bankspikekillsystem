package initialize

import (
	"bank-product-spike-system/global"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

func InitRedis() {

	client := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
	global.REDIS = client
}
