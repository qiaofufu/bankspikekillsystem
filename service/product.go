package service

import (
	"bank-product-spike-system/global"
	"bank-product-spike-system/models"

	"errors"
	"fmt"
	"gorm.io/gorm"
)

type ProductService struct{}

// ReleaseProduct 生成产品
func (receiver ProductService) ReleaseProduct(product *models.Product) (err error) {
	if err = global.DB.Where("product_name = ?", product.ProductName).First(&product).Error; err != gorm.ErrRecordNotFound {
		err = errors.New("产品名已存在，请重新输入产品名")
		return
	}
	err = nil
	result := global.DB.Create(product)
	if result.Error != nil {
		err = errors.New(fmt.Sprintf("生成产品失败，%v\n", result.Error))
		return err
	}
	return
}

// UpdateProduct 更新产品
func (receiver ProductService) UpdateProduct(product *models.Product) (err error) {
	global.DB.Transaction(func(tx *gorm.DB) error {
		err = tx.Where("product_id = ?", product.ID).Delete(&models.ProductTags{}).Error
		if err != nil {
			return err
		}
		err = tx.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&product).Error
		if err != nil {
			err = errors.New(fmt.Sprintf("更新产品失败，%v\n", tx.Error))
			return err
		}
		return nil
	})
	return
}

// DeleteProduct 删除产品
func (receiver ProductService) DeleteProduct(product *models.Product) (err error) {
	result := global.DB.Delete(&product)
	if result.Error != nil {
		err = errors.New(fmt.Sprintf("删除产品失败，%v\n", result.Error))
		return
	}
	return
}

// GetAllProduct  获取所有产品
func (receiver ProductService) GetAllProduct(productName string, interestRateType int, status int, page int, pageSize int) (products []models.Product, cnt int64, err error) {

	result := global.DB.Model(&models.Product{})

	if productName != "" {
		result = result.Where("product_name = ?", productName)
		if result.Error != nil {
			err = errors.New("内部错误")
			return
		}
	}

	if interestRateType != 0 {
		result = result.Where("interest_rate_type = ?", interestRateType)
		if result.Error != nil {
			err = errors.New("内部错误")
			return
		}
	}

	if status != 0 {
		result = result.Where("audit_status = ?", status)
		if result.Error != nil {
			err = errors.New("内部错误")
			return
		}
	}

	result.Count(&cnt)
	result = result.Preload("ProductTags").Offset((page - 1) * pageSize).Limit(pageSize).Find(&products)
	if result.Error != nil {
		err = errors.New(fmt.Sprintf("获取产品失败，%v", result.Error))
		return
	}
	return
}

// GetVisibleProductListOrderByWeight 获取可见产品按权重排序
func (receiver ProductService) GetVisibleProductListOrderByWeight(page int, pageSize int) (list []models.Product, err error) {
	result := global.DB.Model(&models.Product{}).Order("weights desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list)
	if result.Error != nil {
		err = errors.New(fmt.Sprintf("获取产品失败"))
	}
	return
}

// GetVisibleProduct 根据ID获取可见产品
func (receiver ProductService) GetVisibleProduct(id uint) (product models.Product, err error) {
	result := global.DB.Where("id = ? && is_visible = ? && audit_status", id, true, models.AuditPassed).Preload("ProductTags").First(&product)
	if result.Error != nil {
		err = errors.New(fmt.Sprintf("删除产品失败，%v\n", result.Error))
		return
	}
	return
}

// GetProductByID 根据ID获取产品
func (receiver ProductService) GetProductByID(id uint) (product models.Product, err error) {
	result := global.DB.Where("id = ?", id).Preload("ProductTags").First(&product)
	if result.Error != nil {
		err = errors.New(fmt.Sprintf("获取产品失败，%v\n", result.Error))
		return
	}
	return
}

// GetsProductsByIDArr 根据产品ID数组获取产品
func (receiver ProductService) GetsProductsByIDArr(id []uint) (products []models.Product, err error) {
	for _, v := range id {
		elem, err1 := receiver.GetProductByID(v)
		if err1 != nil {
			err = err1
			return
		}
		products = append(products, elem)
	}
	return
}

// GetTagsList 获取产品标签列表
func (receiver ProductService) GetTagsList(page int, pageSize int) (tags []models.Tag, cnt int64, err error) {
	global.DB.Model(&models.Tag{}).Count(&cnt)
	result := global.DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&tags)
	if result.Error != nil {
		err = errors.New(fmt.Sprintf("获取产品标签列表失败，%v", result.Error))
		return
	}
	return
}
