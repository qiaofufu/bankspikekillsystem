package service

import (
	"bank-product-spike-system/global"
	"bank-product-spike-system/models"

	"bank-product-spike-system/utils"
	"errors"
	"fmt"
	"time"
)

type UserService struct{}

// Register 用户注册
func (receiver UserService) Register(phone string, password string, realName string, IDCard string) error {
	if !utils.VerifyIDCard(IDCard) {
		return errors.New("身份证不合法")
	}
	info, err := utils.GetInfoByIDCard(IDCard)
	if err != nil {
		return errors.New("获取身份证信息错误")
	}
	user := models.User{
		PhoneNumber: phone,
		Password:    utils.MD5(password),
		RealInformation: models.RealInformation{
			RealName: realName,
			IDCard:   IDCard,
		},
		Address: models.Address{
			Province: info.AddressTree[0],
			City:     info.AddressTree[1],
			District: info.AddressTree[2],
		},
		Age:    time.Now().Year() - info.Birthday.Year(),
		Gender: info.Sex + 1,
	}
	result := global.DB.Create(&user)
	if result.Error != nil {
		global.Logger.ErrorPrintf("注册失败!\t错误信息:\t%v\n" + result.Error.Error())
		return errors.New("注册失败！")
	}
	return nil
}

// LoginByPassword 密码登陆
func (receiver UserService) LoginByPassword(user *models.User) (*models.User, error) {
	user.Password = utils.MD5(user.Password)
	if err := global.DB.Where("phone_number = ? and password = ?", user.PhoneNumber, user.Password).First(&user).Error; err != nil {
		fmt.Println("err")
		return nil, errors.New("账号或密码错误请重新输入")
	}
	global.DB.Save(user)
	return user, nil
}

// LoginByPhone 手机登陆
func (receiver UserService) LoginByPhone(user *models.User) (*models.User, error) {
	if err := global.DB.Where("phone_number = ?", user.PhoneNumber).First(&user).Error; err != nil {
		return nil, errors.New("没有此账号信息，请注册后登录！")
	}
	global.DB.Save(user)
	return user, nil
}

// GetUserByToken 获取用户信息通过Token
func (receiver UserService) GetUserByToken(tokenStr string) (user models.User, err error) {
	token, err1 := utils.GetTokenPayload(tokenStr)
	if err1 != nil {
		err = err1
		return
	}

	uuid, ok := token.Get("uuid")
	if !ok {
		err = errors.New("token信息不完整，请重新登录")
		return
	}

	result := global.DB.Where("uuid = ?", uuid).First(&user)
	if result.Error != nil {
		err = errors.New(fmt.Sprintf("获取用户信息失败，请重新登录，%v", result.Error))
		return
	}
	return
}

// GetUserList 获取用户列表
func (receiver UserService) GetUserList(page int, pageSize int, age int, gender int, status int) (list []models.User, total int64, err error) {
	tx := global.DB.Model(&models.User{})
	if status != 0 {
		tx = tx.Where("status = ?", status)
	}
	if age != 0 {
		tx = tx.Where("age = ?", age)
	}
	if gender != 0 {
		tx = tx.Where("gender = ?", gender)
	}
	err = tx.Count(&total).Error
	if err != nil {
		err = errors.New(fmt.Sprintf("获取总数失败"))
		return
	}
	err = tx.Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return
}

// GetUserIDByUUID 获取用户信息通过UUID
func (receiver UserService) GetUserIDByUUID(uuid string) (id uint, err error) {
	u := models.User{}
	result := global.DB.Where("uuid = ?", uuid).First(&u)
	if result.Error != nil {
		err = errors.New("获取用户信息错误")
	}
	id = u.ID
	return
}

// UpdateUserInfo 更新用户信息
func (receiver UserService) UpdateUserInfo(user models.User) (err error) {
	result := global.DB.Updates(&user)
	if result.Error != nil {
		err = errors.New(fmt.Sprintf("更新失败，%v", result.Error))
	}
	return
}

// ChangePassword 更改密码服务
func (receiver UserService) ChangePassword(user *models.User, newPassword string) (*models.User, error) {
	password := utils.MD5(user.Password)
	result := global.DB.Where("uuid = ? and password = ?", user.UUID, password).First(&user)
	if result.Error != nil {
		return nil, errors.New("账号或密码错误")
	}
	user.Password = utils.MD5(newPassword)
	result = global.DB.Save(user)
	if result.Error != nil {
		return nil, errors.New("内部错误更新密码失败")
	}
	return user, nil
}

func (receiver UserService) GetUserInfo(userID uint) (user models.User, err error) {
	err1 := global.DB.Where("id = ?", userID).First(&user).Error
	if err1 != nil {
		err = errors.New("获取用户信息失败")
		return
	}
	return
}
