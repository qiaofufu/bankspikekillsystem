package captcha

import (
	"bank-product-spike-system/global"
	"context"
	"fmt"
	"time"
)

type RedisStore struct{}

func (receiver RedisStore) Set(id string, value string) error {
	fmt.Println(value)
	ctx := context.Background()
	status := global.REDIS.Set(ctx, id, value, time.Second*180)
	return status.Err()
}

func (receiver RedisStore) Get(id string, clear bool) string {
	ctx := context.Background()
	values := global.REDIS.Get(ctx, id)
	v := values.Val()
	if clear {
		global.REDIS.Del(ctx, id)
	}
	return v
}

func (receiver RedisStore) Verify(id, answer string, clear bool) bool {
	ans := receiver.Get(id, clear)
	if ans == answer {
		return true
	} else {
		return false
	}
}
