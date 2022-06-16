package response

import "bank-product-spike-system/models"

type ProductList struct {
	Total int64
	List  []models.Product
}

type TagsList struct {
	Total int64
	List  []models.Tag
}

type ProductBaseInfo struct {
	ProductID           uint    `json:"product_id"`            // 产品ID
	ProductName         string  `json:"product_name"`          // 产品名称
	ProductInterestRate float64 `json:"product_interest_rate"` // 产品利率
}

type GetFeaturedProduct struct {
	List []ProductBaseInfo `json:"list"`
}
