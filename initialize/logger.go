package initialize

import "bank-product-spike-system/global"

func InitLogger() {
	global.Logger = global.NewLog()
}
