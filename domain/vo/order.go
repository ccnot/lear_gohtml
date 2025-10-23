package vo

import "time"

// Order 订单信息
type Order struct {
	ID            int64       `json:"id"`
	OrderNo       string      `json:"order_no"`        // 订单编号
	CustomerName  string      `json:"customer_name"`   // 客户名称
	CustomerEmail string      `json:"customer_email"`  // 客户邮箱
	TotalAmount   float64     `json:"total_amount"`    // 总金额
	Status        string      `json:"status"`          // pending, paid, shipped, completed, cancelled
	PaymentMethod string      `json:"payment_method"`  // 支付方式
	Items         []OrderItem `json:"items,omitempty"` // 订单项
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

// OrderItem 订单项
type OrderItem struct {
	ID          int64   `json:"id"`
	ProductName string  `json:"product_name"`
	SKU         string  `json:"sku"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Subtotal    float64 `json:"subtotal"`
}

// OrderListData 订单列表数据
type OrderListData struct {
	Orders   []Order  `json:"orders"`
	PageInfo PageInfo `json:"page_info"`
	Query    string   `json:"query"`  // 当前搜索关键词
	Status   string   `json:"status"` // 当前状态筛选
}

// OrderDetailData 订单详情数据
type OrderDetailData struct {
	Order    Order `json:"order"`
	IsModal  bool  `json:"is_modal"`
}
