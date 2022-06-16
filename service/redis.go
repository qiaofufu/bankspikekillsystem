package service

import (
	"bank-product-spike-system/global"
	"bank-product-spike-system/models"

	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type RedisServer struct{}

const (
	Stock = "ProductStock="
)

func (r RedisServer) setActivity(activity models.SpikeActivity) (err error) {
	datas, err1 := json.Marshal(activity)
	if err1 != nil {
		err = errors.New("设置缓存失败，内部错误")
		return
	}
	key := fmt.Sprintf("spike_activity_id=%d", activity.ID)
	ctx := context.Background()
	_, err = global.REDIS.Set(ctx, key, datas, time.Minute*10).Result()
	return
}

func (r RedisServer) GetActivityBaseInfo(activityID uint) (activity *models.SpikeActivity, err error) {
	key := fmt.Sprintf("spike_activity_id=%d", activityID)
	ctx := context.Background()
	val, err1 := global.REDIS.Get(ctx, key).Result()
	switch {
	case err1 == redis.Nil:
		fmt.Println("活动未加载")
		activity, err = SpikeActivityService{}.GetActivityBaseInfoByID(activityID)
		if err != nil {
			return
		}
		err = r.setActivity(*activity)
		if err != nil {
			return
		}

	case err1 != nil:
		fmt.Println("内部错误")
		err = errors.New("内部错误")
		return
	default:
		json.Unmarshal([]byte(val), &activity)
	}
	return
}

func (r RedisServer) setSpikeResult(res models.SpikeResult) (err error) {
	datas, err := json.Marshal(res)
	if err != nil {
		err = errors.New("设置缓存失败，内部错误")
		return
	}
	key := fmt.Sprintf("spike_result{activity_id=%d,user_id=%d}", res.ActivityID, res.UserID)
	ctx := context.Background()
	_, err = global.REDIS.Set(ctx, key, datas, time.Minute*10).Result()
	return
}

//func (r RedisServer) GetSpikeResult(activityID uint, userID uint) (res *models.SpikeResult, err error) {
//	key := fmt.Sprintf("spike_result{activity_id=%d,user_id=%d}", activityID, userID)
//	ctx := context.Background()
//	val, err1 := global.REDIS.Get(ctx, key).Result()
//	switch {
//	case err1 == redis.Nil:
//		fmt.Println("结果不存在")
//		result, err1 := SpikeActivityService{}.GetSpikeResult(activityID, userID)
//		if err1 != nil {
//			err = err1
//			return
//		}
//		err = r.setSpikeResult(result)
//		if err != nil {
//			return
//		}
//		res = &result
//	case err1 != nil:
//		fmt.Println("内部错误")
//		err = errors.New("内部错误")
//		return
//	default:
//		json.Unmarshal([]byte(val), &res)
//	}
//	return
//}

// SetProductStock 设置产品库存
func (r RedisServer) SetProductStock(productID uint, number int64, expiration time.Duration) (err error) {
	var product models.Product
	err = global.DB.Where("id = ?", productID).First(&product).Error

	ctx := context.Background()
	key := Stock + fmt.Sprintf("%d", productID)
	if number != 0 {
		_, err = global.REDIS.Set(ctx, key, number, expiration).Result()
	} else {
		_, err = global.REDIS.Set(ctx, key, product.ProductNumber-product.SoldNumber, expiration).Result()
	}
	if err != nil {
		errors.New("设置缓存信息失败")
	}

	return
}

// GetProductStock 获取产品库存
func (r RedisServer) GetProductStock(productID uint) (stock int64, err error) {
	ctx := context.Background()
	key := Stock + fmt.Sprintf("%d", productID)
	val, err1 := global.REDIS.Get(ctx, key).Result()
	switch {
	case err1 == redis.Nil:
		err = errors.New(fmt.Sprintf("库存信息不存在"))
		return
	case err1 != nil:
		err = errors.New("内部错误，err:" + err1.Error())
		return
	}
	stock, err = strconv.ParseInt(val, 10, 64)
	return
}

// DeleteProductStock 删除产品库存
func (r RedisServer) DeleteProductStock(productID uint) (ordStock int64, err error) {
	ctx := context.Background()
	key := Stock + fmt.Sprintf("%d", productID)

	cnt, existErr := global.REDIS.Exists(ctx, key).Result()
	if existErr != nil {
		err = errors.New(fmt.Sprintf("检查存在错误"))
		return
	}

	if cnt <= 0 {
		err = errors.New("库存不存在")
		return
	}

	val, delErr := global.REDIS.GetDel(ctx, key).Result()
	if delErr != nil {
		err = errors.New(fmt.Sprintf("删除产品库存缓存失败"))
		return
	}
	ordStock, err = strconv.ParseInt(val, 10, 64)
	return
}

// GetRemainTime 获取缓存剩余时间
func (r RedisServer) GetRemainTime(key string) (expiration time.Duration, err error) {
	ctx := context.Background()
	expiration, ttlErr := global.REDIS.TTL(ctx, key).Result()
	if ttlErr != nil {
		err = errors.New(fmt.Sprintf("获取失败"))
	}
	return
}

// DecreaseStock 减少库存
func (r RedisServer) DecreaseStock(productID uint) (err error) {
	key := Stock + fmt.Sprintf("%d", productID)
	// 获取缓存
	stock, getErr := r.GetProductStock(productID)
	if getErr != nil {
		err = errors.New(fmt.Sprintf("库存服务未加载"))
		return
	}
	if stock <= 0 {
		err = errors.New(fmt.Sprintf("库存不足，秒杀失败"))
		return
	}
	// 获取旧缓存的剩余时间
	expiration, err := r.GetRemainTime(key)
	// 删除旧的缓存
	ordStock, err := r.DeleteProductStock(productID)
	if err != nil {
		return
	}
	// 设置新库存
	err = r.SetProductStock(productID, ordStock-1, expiration)
	return
}
