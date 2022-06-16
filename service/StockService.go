package service

import (
	"fmt"
	"strconv"
	"time"
)

type StockService struct{}

const (
	Running = 1
	Ending  = -1
)

// GetCacheKey 获取缓存KEY
func (receiver StockService) GetCacheKey(productID uint) (key string) {
	return fmt.Sprintf("StockProductID=%d", productID)
}

func (receiver StockService) GetStatusCacheKey(productID uint) (key string) {
	return fmt.Sprintf("ProductStatusID=%d", productID)
}

// SetProductStockCache 设置产品缓存
func (receiver StockService) SetProductStockCache(productID uint, number int64, expiration time.Duration) (err error) {
	key := receiver.GetCacheKey(productID)
	err = CacheService{}.Set(key, number, expiration)
	return
}

// SetProductStatusCache 设置产品状态缓存
func (receiver StockService) SetProductStatusCache(productID uint, status int, expiration time.Duration) (err error) {
	key := receiver.GetStatusCacheKey(productID)
	err = CacheService{}.Delete(key)
	err = CacheService{}.Set(key, status, expiration)
	return
}

// GetProductStatus 获取产品状态
func (receiver StockService) GetProductStatus(productID uint) (status int, err error) {
	key := receiver.GetStatusCacheKey(productID)
	data, getErr := CacheService{}.Get(key)
	if getErr != nil {
		err = getErr
		return
	}
	status, err = strconv.Atoi(data)
	return
}

// DecreaseStock 减少库存
func (receiver StockService) DecreaseStock(productID uint, number int) (err error) {
	key := receiver.GetCacheKey(productID)
	err = CacheService{}.Decrease(key, number)
	return
}

// GetProductStock 获取产品库存
func (receiver StockService) GetProductStock(productID uint) (stock int64, err error) {
	key := receiver.GetCacheKey(productID)
	data, getErr := CacheService{}.Get(key)
	if getErr != nil {
		err = getErr
		return
	}
	stock, err = strconv.ParseInt(data, 10, 64)
	return
}
