package utils

import (
	"bank-product-spike-system/global"
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

var credential *common.Credential
var cpf *profile.ClientProfile

func InitCredential() {
	credential = common.NewCredential(viper.GetString("sms.secretID"), viper.GetString("sms.secretKey"))
	cpf = profile.NewClientProfile()

}

// SendSMS
// phone -- eg. +8618547304726
// code --  eg. 165923
func SendSMS(phone string, code string) error {
	client, _ := sms.NewClient(credential, "ap-guangzhou", cpf)
	request := sms.NewSendSmsRequest()
	request.SmsSdkAppId = common.StringPtr(viper.GetString("sms.appID"))
	request.SignName = common.StringPtr(viper.GetString("sms.signature"))
	request.TemplateParamSet = common.StringPtrs([]string{code, "5"})
	request.TemplateId = common.StringPtr(viper.GetString("sms.templateID"))
	request.PhoneNumberSet = common.StringPtrs([]string{phone})
	response, err := client.SendSms(request)
	if err != nil {
		return err
	}
	b, _ := json.Marshal(response.Response)
	// 打印返回的json字符串
	fmt.Printf("%s", b)
	return nil
}

// VerifySMS
// 验证SMS
func VerifySMS(phone string, code string) bool {
	ctx := context.Background()
	result := global.REDIS.Get(ctx, "phone:"+phone)

	if result.Val() == code {
		global.REDIS.Del(ctx, "phone_"+phone)
		return true
	}
	return false
}
