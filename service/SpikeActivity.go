package service

import (
	"bank-product-spike-system/global"
	"bank-product-spike-system/models"
	"bank-product-spike-system/utils/SendMQ"

	"bank-product-spike-system/utils"
	"bank-product-spike-system/utils/Random"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"time"
)

type SpikeActivityService struct{}

// getCacheKey 获取活动缓存Key
func (receiver SpikeActivityService) getCacheKey(activityID uint) (key string) {
	return fmt.Sprintf("spike_activiyt_cache=%d", activityID)
}

// getSpikeResultCacheKey 获取秒杀结果缓存
func (receiver SpikeActivityService) getSpikeResultCacheKey(activityID uint, productID uint, userID uint) string {
	return fmt.Sprintf("spike_result_cache=[activity=%d;product=%d;user=%d]", activityID, productID, userID)
}

// ReleaseSpikeActivity 发布秒杀活动
func (receiver SpikeActivityService) ReleaseSpikeActivity(spikeActivity *models.SpikeActivity) (err error) {
	result := global.DB.Create(&spikeActivity)
	if result.Error != nil {
		err = errors.New(fmt.Sprintf("发布秒杀活动失败，%v", result.Error))
		return
	}
	return
}

// UpdateSpikeActivity 更新秒杀活动
func (receiver SpikeActivityService) UpdateSpikeActivity(spike *models.SpikeActivity) (err error) {
	result := global.DB.Updates(&spike)
	if result.Error != nil {
		err = errors.New(fmt.Sprintf("更新秒杀活动失败，%v", result.Error))
		return
	}
	return
}

// DeleteSpikeActivity 删除秒杀活动
func (receiver SpikeActivityService) DeleteSpikeActivity(activity *models.SpikeActivity) (err error) {
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		//err1 := tx.Where("spike_activity_id = ?", activity.ID).Delete(&models.ActivityProduct{}).Error
		//if err1 != nil {
		//	return err1
		//}

		err1 := tx.Delete(activity).Error
		if err1 != nil {
			return err1
		}
		return nil
	})
	return
}

// GetAllActivityBaseInfo 获取所有活动基本信息
func (receiver SpikeActivityService) GetAllActivityBaseInfo(startTime string, endingTime string, activityStatus int, auditStatus int, page int, pageSize int) (list []models.SpikeActivity, cnt int64, err error) {
	result := global.DB.Model(&models.SpikeActivity{})
	if startTime != "" {
		sTime, err1 := utils.GetTime(startTime)
		if err1 != nil {
			err = errors.New("开始时间格式错误")
			return
		}
		result = result.Where("start_time >= ?", sTime)
		if result.Error != nil {
			err = errors.New("内部错误1")
			return
		}
	}
	if endingTime != "" {
		eTime, err1 := utils.GetTime(endingTime)
		if err1 != nil {
			err = errors.New("开始时间格式错误")
			return
		}
		result = result.Where("start_time <= ?", eTime)
		if result.Error != nil {
			err = errors.New("内部错误1")
			return
		}
	}

	if activityStatus != 0 {
		result = result.Where("activity_status = ?", activityStatus)
		if result.Error != nil {
			err = errors.New("内部错误2")
			return
		}
	}

	if auditStatus != 0 {
		result = result.Where("audit_status = ?", auditStatus)
		if result.Error != nil {
			err = errors.New("内部错误3")
			return
		}
	}

	result.Count(&cnt)
	err = result.Preload("Admin").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return
}

// GetRunningActivity 获取正在进行中的活动
func (receiver SpikeActivityService) GetRunningActivity() (list []models.SpikeActivity, err error) {
	nowTime := time.Now()
	err = global.DB.Where("start_time <= ? and ending_time >= ? and audit_status = ?", nowTime, nowTime, models.AuditPassed).Preload("Product").Find(&list).Error
	return
}

// Attention 关注活动
func (receiver SpikeActivityService) Attention(dto models.ActivityAttention) (err error) {
	if receiver.IsAttention(dto.ActivityID, dto.UserID) {
		err = errors.New("以成功关注，请勿重复关注")
		return
	}
	err1 := global.DB.Create(&dto).Error
	if err1 != nil {
		err = errors.New("关注失败1")
		return
	}
	err1 = global.DB.Exec("update spike_activities set number_of_favorites = number_of_favorites + 1 where id = ?", dto.ActivityID).Error
	if err1 != nil {
		err = errors.New("关注失败2")
		return
	}
	return
}

// UnAttention 取消关注活动
func (receiver SpikeActivityService) UnAttention(activityID uint, userID uint) (err error) {
	result := global.DB.Where("activity_id = ? and user_id = ?", activityID, userID).Delete(&models.ActivityAttention{})
	if result.Error != nil {
		err = errors.New(fmt.Sprintf("取消关注活动失败， err = %v", result.Error))
	}
	err1 := global.DB.Exec("update spike_activities set number_of_favorites = number_of_favorites - 1 where id = ?", activityID).Error
	if err1 != nil {
		err = errors.New("关注失败2")
		return
	}
	return
}

// IsAttention 是否关注
func (receiver SpikeActivityService) IsAttention(activityID uint, userID uint) bool {
	var cnt int64
	global.DB.Debug().Model(models.ActivityAttention{}).Where("activity_id = ? and user_id = ?", activityID, userID).Count(&cnt)
	if cnt >= 1 {
		return true
	} else {
		return false
	}
}

// GetAttentionListByUserID 获取关注列表通过UserID
func (receiver SpikeActivityService) GetAttentionListByUserID(userID uint) (list []models.ActivityAttention, err error) {
	err = global.DB.Where("user_id = ?", userID).Find(&list).Error
	return
}

// GetUpcomingActivity 获取即将开放的活动
func (receiver SpikeActivityService) GetUpcomingActivity() (list []models.SpikeActivity, err error) {
	nowTime := time.Now()
	err = global.DB.Where("start_time > ? and Date(start_time) = ? and audit_status = ?", nowTime, nowTime.Format("2006-01-02"), models.AuditPassed).Preload("Product").Find(&list).Error
	return
}

// GetActivityBaseInfoByID 获取活动基本信息
func (receiver SpikeActivityService) GetActivityBaseInfoByID(activityID uint) (activity *models.SpikeActivity, err error) {
	err = global.DB.Where("id = ?", activityID).Preload("Product").First(&activity).Error
	return
}

// GetActivityInfoByID 获取活动全部信息，通过活动ID
func (receiver SpikeActivityService) GetActivityInfoByID(activityID uint) (activity *models.SpikeActivity, err error) {
	err = global.DB.Where("id = ?", activityID).Preload("Product.ProductTags").Preload(clause.Associations).First(&activity).Error
	return
}

// AddActivityGuidanceInformation 添加活动引导信息
func (receiver SpikeActivityService) AddActivityGuidanceInformation(ctx *gin.Context, formName string) (err error) {
	activityID, err1 := strconv.Atoi(ctx.PostForm("activity_id"))
	if err1 != nil {
		err = errors.New("活动id指定不明确")
		return
	}

	var activity models.SpikeActivity
	err = global.DB.Where("id = ?", activityID).First(&activity).Error
	if err != nil {
		return
	}

	if activity.ActivityStatus < models.ActivityResourceConfiguration {
		err = errors.New("活动状态错误，非法访问")
		return
	}

	result, err1 := utils.UploadMultipartFile(ctx, formName)
	if err1 != nil {
		err = err1
		return
	}
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		for _, v := range result {
			err = tx.Create(&models.ActivityGuidanceInformation{
				ActivityID: uint(activityID),
				Url: models.Url{
					Host:         v.Host,
					RelativePATH: v.RelativePath,
				},
			}).Error
			if err != nil {
				return err
			}
		}
		activity.ActivityStatus = models.ActivityWaitAudit
		activity.AuditStatus = models.WaitAudit
		return tx.Save(&activity).Error
	})
	return
}

// UpdateActivityGuidanceInformation 更新活动引导信息
func (receiver SpikeActivityService) UpdateActivityGuidanceInformation(ctx *gin.Context, formName string) (err error) {
	activityID, err1 := strconv.Atoi(ctx.PostForm("activity_id"))
	if err1 != nil {
		err = errors.New("活动id指定不明确")
		return
	}

	var activity models.SpikeActivity
	err = global.DB.Where("id = ?", activityID).First(&activity).Error
	switch {
	case err == gorm.ErrRecordNotFound:
		err = errors.New("活动不存在")
		return
	case err != nil:
		err = errors.New("内部错误")
		return
	}

	if activity.ActivityStatus < models.ActivityResourceConfiguration {
		err = errors.New("活动状态错误，非法访问")
		return
	}

	result, err1 := utils.UploadMultipartFile(ctx, formName)
	if err1 != nil {
		err = err1
		return
	}
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		err = tx.Where("activity_id = ?", activityID).Delete(&models.ActivityGuidanceInformation{}).Error
		if err != nil {
			return err
		}
		for _, v := range result {
			err = tx.Create(&models.ActivityGuidanceInformation{
				ActivityID: uint(activityID),
				Url: models.Url{
					Host:         v.Host,
					RelativePATH: v.RelativePath,
				},
			}).Error
			if err != nil {
				return err
			}
		}
		activity.AuditStatus = models.WaitAudit
		activity.ActivityStatus = models.ActivityWaitAudit
		return tx.Save(activity).Error
	})
	return
}

// GetActivityGuidanceInformation 获取活动引导信息
func (receiver SpikeActivityService) GetActivityGuidanceInformation(activityID uint) (info []models.ActivityGuidanceInformation, err error) {
	result := global.DB.Debug().Where("activity_id = ?", activityID).Find(&info)
	if result.Error != nil {
		err = errors.New(fmt.Sprintf("获取活动引导信息失败， err:%v", result.Error))
	}
	return
}

// GetActivityUserRecord 获取活动用户记录
func (receiver SpikeActivityService) GetActivityUserRecord(activityID uint) (record []models.ActivityUser, err error) {
	result := global.DB.Where("spike_activity_id = ?", activityID).Find(&record)
	if result.Error != nil {
		err = errors.New(fmt.Sprintf("获取记录失败， err:%v", result.Error))
	}
	return
}

// CheckEligibility 检查参与资格
func (receiver SpikeActivityService) CheckEligibility(activityID uint, productID uint, userID uint, verifyCode string, mode int) (err error) {
	if mode == 1 {
		// 检查验证码是否正确
		if !receiver.VerifySpikeCode(activityID, userID, verifyCode) {
			err = errors.New("秒杀路径错误")
			return
		}
	}
	// 获取秒杀缓存
	cache, err1 := receiver.GetSpikeCache(activityID)
	if err1 != nil {
		err = errors.New(fmt.Sprintf("活动缓存获取失败， err[%v]", err1))
		return
	}
	// 检查是否在开放期内
	nowTime := time.Now()
	if !(cache.StartTime.Unix() <= nowTime.Unix() && cache.EndingTime.Unix() >= nowTime.Unix()) {
		err = errors.New("不在活动开放期内")
		return
	}
	// 检查库存
	if cache.ProductNumber < 1 {
		err = errors.New("库存不足，秒杀参与失败")
		return
	}
	return
}

// ParticipateSpike 参与秒杀
func (receiver SpikeActivityService) ParticipateSpike(activityID uint, productID uint, userID uint) (err error) {
	// 检查秒杀记录缓存是否存在
	key := receiver.getSpikeResultCacheKey(activityID, productID, userID)
	isExist, existErr := CacheService{}.Exist(key)
	if existErr != nil {
		err = existErr
		return
	}
	if isExist == true { // 秒杀记录缓存存在, 直接返回
		return
	} else {
		// 秒杀记录缓存不存在， 进行秒杀请求处理
		global.DB.Raw("update spike_activities set number_of_participants = number_of_participants + 1 where id = ?", activityID)
		// 获取秒杀缓存
		cache, err1 := receiver.GetSpikeCache(activityID)
		if err1 != nil {
			err = errors.New(fmt.Sprintf("活动缓存获取失败， err[%v]", err1))
			return
		}
		// 检查库存
		if cache.ProductNumber <= 0 {
			err = errors.New("库存不足，秒杀参与失败")
			return
		}
		// 设置秒杀记录缓存
		result := models.SpikeResult{}
		result.ActivityID = activityID
		result.UserID = userID
		result.Status = models.SpikeWait
		setErr := CacheService{}.Set(key, result, cache.EndingTime.Sub(time.Now()))
		if setErr != nil {
			err = setErr
			return
		}
		// 加入消息队列
		order := models.Order{
			ActivityID:    activityID,
			ProductID:     productID,
			UserID:        userID,
			OrderStatus:   models.OrderWait,
			PaymentStatus: models.PaymentWaiting,
		}
		data, err1 := json.Marshal(order)
		if err1 != nil {
			err = errors.New("订单序列化失败")
			return
		}
		SendMQ.MQ.SendMessage(SendMQ.OrderQueue, data)
	}
	return
}

// Cache 缓存活动信息
func (receiver SpikeActivityService) Cache(activityID uint, time time.Duration) (err error) {
	key := receiver.getCacheKey(activityID)
	isExist, err1 := CacheService{}.Exist(key)
	if err1 != nil {
		err = errors.New("内部错误1")
		return
	}
	if !isExist {
		var activity models.SpikeActivity
		global.DB.Where("id = ?", activityID).Preload("Product").First(&activity)

		activityBaseInfoCache := models.ActivityBaseInfoCache{
			ActivityName: activity.ActivityName,
			ActivityID:   activity.ID,
			StartTime:    activity.StartTime,
			EndingTime:   activity.EndingTime,
			ProductTotal: activity.Product.ProductNumber,
			ProductID:    activity.ProductID,
		}
		err1 = StockService{}.SetProductStockCache(activity.ProductID, activity.Product.ProductNumber-activity.Product.SoldNumber, time)
		if err1 != nil {
			err = errors.New("内部错误2 " + err1.Error())
			return
		}
		err1 = CacheService{}.Set(key, activityBaseInfoCache, time)
		if err1 != nil {
			err = errors.New("内部错误3 " + err1.Error())
			return
		}
	}
	return
}

// GetSpikeCache 获取秒杀缓存
func (receiver SpikeActivityService) GetSpikeCache(activityID uint) (cache models.SpikeActivityCache, err error) {
	key := receiver.getCacheKey(activityID)
	isExist, err1 := CacheService{}.Exist(key)
	if err1 != nil {
		err = err1
		return
	}
	if isExist {

		var actInfo models.ActivityBaseInfoCache
		activityInfo, err1 := CacheService{}.Get(key)
		if err1 != nil {
			err = errors.New("获取缓存信息失败" + err.Error())
			return
		}
		err1 = json.Unmarshal([]byte(activityInfo), &actInfo)
		if err1 != nil {
			err = errors.New("序列化失败")
			return
		}
		stock, err1 := StockService{}.GetProductStock(actInfo.ProductID)
		if err1 != nil {
			err = errors.New("获取缓存信息失败" + err.Error())
			return
		}
		cache = models.SpikeActivityCache{
			ActivityBaseInfoCache: actInfo,
			ProductNumberCache:    models.ProductNumberCache{ProductNumber: stock},
		}
	} else {
		err = errors.New("缓存不存在")
		return
	}
	return
}

// GetSpikeVerifyCode 获取秒杀验证码
func (receiver SpikeActivityService) GetSpikeVerifyCode(activityID uint, userID uint) (verifyCode string) {
	verifyCode = string(Random.Random{}.RandomString(60))
	key := fmt.Sprintf("activityID=%duserID=%d", activityID, userID)
	err := CacheService{}.Set(key, utils.MD5("suse"+verifyCode+"spike"), time.Minute*10)
	fmt.Println(key)
	if err != nil {
		fmt.Println(err.Error())
	}
	return
}

// VerifySpikeCode 验证秒杀验证码
func (receiver SpikeActivityService) VerifySpikeCode(activityID uint, userID uint, verifyCode string) bool {
	key := fmt.Sprintf("activityID=%duserID=%d", activityID, userID)
	data, err := CacheService{}.Get(key)
	if err != nil {
		return false
	}
	if data != "\""+verifyCode+"\"" {

		return false
	}
	err = CacheService{}.Delete(key)
	if err != nil {
		return false
	}
	return true
}

// QuerySpikeResult 查询秒杀结果
func (receiver SpikeActivityService) QuerySpikeResult(activityID uint, product uint, userID uint) (result models.SpikeResult, err error) {
	key := receiver.getSpikeResultCacheKey(activityID, product, userID)
	data, getErr := CacheService{}.Get(key)
	if getErr != nil {
		err = getErr
		return
	}
	parseErr := json.Unmarshal([]byte(data), &result)
	if parseErr != nil {
		err = errors.New(fmt.Sprintf("序列化错误"))
		return
	}
	return
}
