package service

import (
	"bank-product-spike-system/global"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type CacheService struct{}

// Exist 检查缓存是否存在
func (receiver CacheService) Exist(key string) (isExist bool, err error) {
	ctx := context.Background()

	cnt, existErr := global.REDIS.Exists(ctx, key).Result()
	if existErr != nil {
		err = errors.New(fmt.Sprintf("缓存内部错误"))
		return false, err
	}

	if cnt <= 0 {
		return false, nil
	}

	return true, nil
}

// Get 获取缓存
func (receiver CacheService) Get(key string) (data string, err error) {
	ctx := context.Background()

	data, getErr := global.REDIS.Get(ctx, key).Result()

	switch {
	case getErr == redis.Nil:
		err = errors.New(fmt.Sprintf("缓存信息不存在"))
		return
	case getErr != nil:
		err = errors.New(fmt.Sprintf("缓存内部错误"))
		return
	}
	return
}

// Set 设置缓存 expiration:0 持久保存
func (receiver CacheService) Set(key string, data interface{}, expiration time.Duration) (err error) {
	serializedData, err := json.Marshal(data)
	if err != nil {
		err = errors.New("设置缓存失败，数据解析错误")
		return
	}
	ctx := context.Background()
	_, err = global.REDIS.Set(ctx, key, serializedData, expiration).Result()
	fmt.Println(key)
	fmt.Println(expiration)
	return
}

// Delete 删除缓存
func (receiver CacheService) Delete(key string) (err error) {
	ctx := context.Background()
	_, err = global.REDIS.Del(ctx, key).Result()
	return
}

// Decrease 递减缓存数据（要求value为整型）
func (receiver CacheService) Decrease(key string, number int) (err error) {

	ctx := context.Background()
	// Get Expiration
	expiration, err1 := global.REDIS.TTL(ctx, key).Result()
	if err1 != nil {
		err = err1
		return
	}
	// Get old
	data, err1 := global.REDIS.Get(ctx, key).Result()
	// Delete Cache
	_, err1 = global.REDIS.Del(ctx, key).Result()
	if err1 != nil {
		err = err1
		return
	}
	// Set New Cache
	oldData, err1 := strconv.Atoi(data)
	if err1 != nil {
		err = err1
		return
	}
	err = receiver.Set(key, oldData-number, expiration)

	return
}
