// Package controller 订单管理控制器
package controller

import (
	"fmt"
	"godash/domain/vo"
	"godash/infra"
	"math/rand"
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
	BaseController
}

// mockOrders 模拟订单数据库
var mockOrders = []vo.Order{}

// init 初始化订单 mock 数据
func init() {
	customers := []string{"张三", "李四", "王五", "赵六", "孙七", "周八", "吴九", "郑十"}
	statuses := []string{"pending", "paid", "shipped", "completed", "cancelled"}
	payments := []string{"支付宝", "微信支付", "银行卡", "货到付款"}

	for i := 0; i < 300; i++ {
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

	// 使用基础控制器的搜索助手
	params, pagination := c.SearchHelper(params)
	filteredOrders := c.filterOrders(params)

	// 转换为 interface{} 进行分页
	orders := make([]interface{}, len(filteredOrders))
	for i, order := range filteredOrders {
		orders[i] = order
	}

	pagedOrders, pagination := c.Paginate(orders, pagination)

	// 转换回订单类型
	result := make([]vo.Order, len(pagedOrders))
	for i, order := range pagedOrders {
		result[i] = order.(vo.Order)
	}

	data := vo.OrderListData{
		Orders:   result,
		PageInfo: c.CreatePageInfo(pagination),
		Query:    params.Keyword,
		Status:   params.Status,
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

	// 检查是否为模态框请求（通过检查请求头或查询参数）
	hxTarget := c.Worker.IrisContext().GetHeader("HX-Target")
	isModal := hxTarget == "order-modal-content" || hxTarget == "#order-modal-content" ||
		c.Worker.IrisContext().URLParamExists("modal")

	data := vo.OrderDetailData{
		Order:   *order,
		IsModal: isModal,
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
		Return string `form:"return"`
	}

	if err := c.Request.ReadForm(&statusData, true); err != nil {
		c.SetErrorToast("状态更新失败: " + err.Error())
		return &infra.JSONResponse{Error: err}
	}

	// 查找并更新订单状态
	if !c.updateOrderStatus(id, statusData.Status) {
		return c.HandleNotFoundError("订单")
	}

	c.SetSuccessToast("订单状态更新成功")

	// 根据 return 参数决定返回订单行还是订单详情
	order := c.findOrderByID(id)
	if statusData.Return == "detail" {
		// 检查是否为模态框请求
		hxTarget := c.Worker.IrisContext().GetHeader("HX-Target")
		isModal := hxTarget == "order-modal-content" || hxTarget == "#order-modal-content"

		data := vo.OrderDetailData{
			Order:   *order,
			IsModal: isModal,
		}
		return &infra.ViewResponse{
			Name: "orders/detail.html",
			Data: data,
		}
	}

	// 默认返回订单行（用于列表页面）
	return &infra.ViewResponse{
		Name: "orders/row.html",
		Data: order,
	}
}

// DeleteBy 取消订单
// DELETE /orders/{id}
func (c *OrderController) DeleteBy(id int64) freedom.Result {
	if !c.updateOrderStatus(id, "cancelled") {
		return c.HandleNotFoundError("订单")
	}

	c.SetSuccessToast("订单已取消")

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
		// 使用基础控制器的过滤助手
		if params.Keyword != "" {
			if !c.FilterHelper(params.Keyword, order.OrderNo) &&
				!c.FilterHelper(params.Keyword, order.CustomerName) &&
				!c.FilterHelper(params.Keyword, order.CustomerEmail) {
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

// updateOrderStatus 更新订单状态
func (c *OrderController) updateOrderStatus(id int64, status string) bool {
	for i, order := range mockOrders {
		if order.ID == id {
			mockOrders[i].Status = status
			mockOrders[i].UpdatedAt = time.Now()
			return true
		}
	}
	return false
}
