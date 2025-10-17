package vo

import "time"

// Product 商品信息
type Product struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	SKU         string    `json:"sku"`         // 商品编码
	Category    string    `json:"category"`    // 分类
	Price       float64   `json:"price"`       // 价格
	Stock       int       `json:"stock"`       // 库存
	Status      string    `json:"status"`      // active, inactive, out_of_stock
	Image       string    `json:"image"`       // 商品图片URL
	Description string    `json:"description"` // 描述
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProductListData 商品列表数据
type ProductListData struct {
	Products []Product `json:"products"`
	PageInfo PageInfo  `json:"page_info"`
	Query    string    `json:"query"`    // 当前搜索关键词
	Category string    `json:"category"` // 当前分类筛选（保留，已经有了）
}

// ProductFormData 商品表单数据
type ProductFormData struct {
	ID          int64   `json:"id" form:"id"`
	Name        string  `json:"name" form:"name" validate:"required"`
	SKU         string  `json:"sku" form:"sku" validate:"required"`
	Category    string  `json:"category" form:"category" validate:"required"`
	Price       float64 `json:"price" form:"price" validate:"required,gt=0"`
	Stock       int     `json:"stock" form:"stock" validate:"required,gte=0"`
	Status      string  `json:"status" form:"status" validate:"required"`
	Description string  `json:"description" form:"description"`
}
