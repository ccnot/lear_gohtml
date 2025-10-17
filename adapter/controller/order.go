// Package controller 订单管理控制器
package controller

import (
	"fmt"
	"gohtml/domain/vo"
	"gohtml/infra"
	"math"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/8treenet/freedom"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		// 绑定订单控制器到 /orders 路由
		initiator.BindController("/orders", &OrderController{})
	})
}

// OrderController 订单管理控制器
type OrderController struct {
	Worker  freedom.Worker
	Request *infra.Request
}

// mockOrders 模拟订单数据库
var mockOrders = []vo.Order{}

// init 初始化订单 mock 数据
func init() {
	customers := []string{"张三", "李四", "王五", "赵六", "孙七", "周八", "吴九", "郑十"}
	statuses := []string{"pending", "paid", "shipped", "completed", "cancelled"}
	payments := []string{"支付宝", "微信支付", "银行卡", "货到付款"}

	for i := 0; i < 30; i++ {
		// 生成订单项
		itemCount := rand.Intn(3) + 1
		items := make([]vo.OrderItem, itemCount)
		totalAmount := 0.0

		for j := 0; j < itemCount; j++ {
			price := float64(rand.Intn(500)+50) + 0.99
			quantity := rand.Intn(3) + 1
			subtotal := price * float64(quantity)
			totalAmount += subtotal

			items[j] = vo.OrderItem{
				ID:          int64(j + 1),
				ProductName: fmt.Sprintf("商品-%d", j+1),
				SKU:         fmt.Sprintf("SKU%05d", rand.Intn(1000)),
				Quantity:    quantity,
				Price:       price,
				Subtotal:    subtotal,
			}
		}

		order := vo.Order{
			ID:            int64(i + 1),
			OrderNo:       fmt.Sprintf("ORD%s%04d", time.Now().Format("20060102"), i+1),
			CustomerName:  customers[i%len(customers)],
			CustomerEmail: fmt.Sprintf("customer%d@example.com", i+1),
			TotalAmount:   totalAmount,
			Status:        statuses[i%len(statuses)],
			PaymentMethod: payments[i%len(payments)],
			Items:         items,
			CreatedAt:     time.Now().Add(-time.Duration(i) * 24 * time.Hour),
			UpdatedAt:     time.Now().Add(-time.Duration(i) * time.Hour),
		}
		mockOrders = append(mockOrders, order)
	}
}

// Get 获取订单列表
// GET /orders
func (c *OrderController) Get() freedom.Result {
	// 读取查询参数
	var params vo.SearchParams
	if err := c.Request.ReadQuery(&params, false); err != nil {
		params = vo.SearchParams{}
	}

	// 设置默认值
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	// 过滤和搜索
	filteredOrders := c.filterOrders(params)

	// 分页
	total := int64(len(filteredOrders))
	totalPages := int(math.Ceil(float64(total) / float64(params.PageSize)))
	start := (params.Page - 1) * params.PageSize
	end := start + params.PageSize
	if end > len(filteredOrders) {
		end = len(filteredOrders)
	}
	if start > len(filteredOrders) {
		start = len(filteredOrders)
	}

	pageOrders := filteredOrders[start:end]

	data := vo.OrderListData{
		Orders: pageOrders,
		PageInfo: vo.PageInfo{
			Page:       params.Page,
			PageSize:   params.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
		Query:  params.Keyword,
		Status: params.Status,
	}

	return &infra.ViewResponse{
		Name: "orders/list.html",
		Data: data,
	}
}

// GetBy 获取订单详情
// GET /orders/{id}
func (c *OrderController) GetBy(id int64) freedom.Result {
	order := c.findOrderByID(id)
	if order == nil {
		return &infra.JSONResponse{
			Code:  404,
			Error: fmt.Errorf("订单不存在"),
		}
	}

	data := vo.OrderDetailData{
		Order: *order,
	}

	return &infra.ViewResponse{
		Name: "orders/detail.html",
		Data: data,
	}
}

// PutStatusBy 更新订单状态
// PUT /orders/{id}/status
func (c *OrderController) PutStatusBy(id int64) freedom.Result {
	var statusData struct {
		Status string `form:"status" validate:"required"`
	}

	if err := c.Request.ReadForm(&statusData, true); err != nil {
		c.Worker.IrisContext().Header("X-Toast-Message", url.QueryEscape("状态更新失败: "+err.Error()))
		c.Worker.IrisContext().Header("X-Toast-Type", "error")
		return &infra.JSONResponse{Error: err}
	}

	// 查找并更新订单状态
	found := false
	for i, o := range mockOrders {
		if o.ID == id {
			mockOrders[i].Status = statusData.Status
			mockOrders[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		c.Worker.IrisContext().Header("X-Toast-Message", url.QueryEscape("订单不存在"))
		c.Worker.IrisContext().Header("X-Toast-Type", "error")
		return &infra.JSONResponse{
			Code:  404,
			Error: fmt.Errorf("订单不存在"),
		}
	}

	// 设置成功提示
	c.Worker.IrisContext().Header("X-Toast-Message", url.QueryEscape("订单状态更新成功"))
	c.Worker.IrisContext().Header("X-Toast-Type", "success")

	// 返回更新后的订单行
	order := c.findOrderByID(id)
	return &infra.ViewResponse{
		Name: "orders/row.html",
		Data: order,
	}
}

// DeleteBy 取消订单
// DELETE /orders/{id}
func (c *OrderController) DeleteBy(id int64) freedom.Result {
	// 查找并更新订单状态为已取消
	found := false
	for i, o := range mockOrders {
		if o.ID == id {
			mockOrders[i].Status = "cancelled"
			mockOrders[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		c.Worker.IrisContext().Header("X-Toast-Message", url.QueryEscape("订单不存在"))
		c.Worker.IrisContext().Header("X-Toast-Type", "error")
		return &infra.JSONResponse{
			Code:  404,
			Error: fmt.Errorf("订单不存在"),
		}
	}

	// 设置成功提示
	c.Worker.IrisContext().Header("X-Toast-Message", url.QueryEscape("订单已取消"))
	c.Worker.IrisContext().Header("X-Toast-Type", "success")

	// 返回更新后的订单行
	order := c.findOrderByID(id)
	return &infra.ViewResponse{
		Name: "orders/row.html",
		Data: order,
	}
}

// BeforeActivation 配置路由
func (c *OrderController) BeforeActivation(b freedom.BeforeActivation) {
	b.Handle("GET", "/{id:int64}", "GetBy")
	b.Handle("PUT", "/{id:int64}/status", "PutStatusBy")
	b.Handle("DELETE", "/{id:int64}", "DeleteBy")
}

// filterOrders 过滤订单
func (c *OrderController) filterOrders(params vo.SearchParams) []vo.Order {
	filtered := []vo.Order{}

	for _, order := range mockOrders {
		// 搜索过滤
		if params.Keyword != "" {
			keyword := strings.ToLower(params.Keyword)
			if !strings.Contains(strings.ToLower(order.OrderNo), keyword) &&
				!strings.Contains(strings.ToLower(order.CustomerName), keyword) &&
				!strings.Contains(strings.ToLower(order.CustomerEmail), keyword) {
				continue
			}
		}

		// 状态过滤
		if params.Status != "" && order.Status != params.Status {
			continue
		}

		filtered = append(filtered, order)
	}

	return filtered
}

// findOrderByID 根据 ID 查找订单
func (c *OrderController) findOrderByID(id int64) *vo.Order {
	for _, order := range mockOrders {
		if order.ID == id {
			return &order
		}
	}
	return nil
}
