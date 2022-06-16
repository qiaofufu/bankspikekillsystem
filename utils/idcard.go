package utils

import (
	idvalidator "github.com/guanguans/id-validator"
	"time"
)

type IdInfo struct {
	AddressCode   int
	Abandoned     int
	Address       string
	AddressTree   []string
	Birthday      time.Time
	Constellation string
	ChineseZodiac string
	Sex           int
	Length        int
	CheckBit      string
}

// GetInfoByIDCard
// 通过身份证获取信息
func GetInfoByIDCard(IDCard string) (IdInfo, error) {
	info, err := idvalidator.GetInfo(IDCard, true)
	return IdInfo{
		AddressCode:   info.AddressCode,
		Abandoned:     info.Abandoned,
		Address:       info.Address,
		AddressTree:   info.AddressTree,
		Birthday:      info.Birthday,
		Constellation: info.Constellation,
		ChineseZodiac: info.ChineseZodiac,
		Sex:           info.Sex,
		Length:        info.Length,
		CheckBit:      info.CheckBit,
	}, err
}

func VerifyIDCard(IDCard string) bool {
	return idvalidator.IsValid(IDCard, true)
}
