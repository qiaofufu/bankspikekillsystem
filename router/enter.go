package router

import "bank-product-spike-system/router/system"

type RouterGroup struct {
	SystemRouter system.RouterGroup
}

var RouterGroupAPP = new(RouterGroup)
