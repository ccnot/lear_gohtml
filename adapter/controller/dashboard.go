// Package controller 仪表盘控制器
package controller

import (
	"godash/domain/vo"
	"godash/infra"
	"math/rand"
	"time"

	"github.com/8treenet/freedom"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		// 绑定仪表盘控制器到 /dashboard 路由
		initiator.BindController("/dashboard", &DashboardController{})
	})
}

// DashboardController 仪表盘控制器
type DashboardController struct {
	Worker  freedom.Worker
	Request *infra.Request
}

// Get 获取仪表盘数据
// GET /dashboard
func (c *DashboardController) Get() freedom.Result {
	// 生成 mock 统计数据
	stats := vo.DashboardStats{
		TotalUsers:    1234,
		TotalProducts: 567,
		TotalOrders:   8901,
		TotalRevenue:  234567.89,
		ActiveUsers:   892,
		PendingOrders: 23,
		LowStock:      8,
		TodayOrders:   45,
	}

	// 生成 mock 最近订单
	recentOrders := c.generateMockOrders(5)

	// 生成 mock 最近用户
	recentUsers := c.generateMockUsers(5)

	// 生成 mock 热门商品
	topProducts := c.generateMockProducts(4)

	data := vo.DashboardData{
		Stats:        stats,
		RecentOrders: recentOrders,
		RecentUsers:  recentUsers,
		TopProducts:  topProducts,
	}

	return &infra.ViewResponse{
		Name: "dashboard/index.html",
		Data: data,
	}
}

// GetStats 获取统计数据（用于定时刷新）
// GET /dashboard/stats
func (c *DashboardController) GetStats() freedom.Result {
	// 生成动态变化的统计数据，模拟实时更新
	rand.Seed(time.Now().UnixNano())

	stats := vo.DashboardStats{
		TotalUsers:    1234 + int64(rand.Intn(10)),
		TotalProducts: 567,
		TotalOrders:   8901 + int64(rand.Intn(5)),
		TotalRevenue:  234567.89 + float64(rand.Intn(1000)),
		ActiveUsers:   892 + int64(rand.Intn(20)),
		PendingOrders: 23 + int64(rand.Intn(3)),
		LowStock:      8,
		TodayOrders:   45 + int64(rand.Intn(5)),
	}

	return &infra.ViewResponse{
		Name: "dashboard/stats.html",
		Data: stats,
	}
}

// BeforeActivation 配置路由
func (c *DashboardController) BeforeActivation(b freedom.BeforeActivation) {
	b.Handle("GET", "/stats", "GetStats")
}

// generateMockOrders 生成 mock 订单数据
func (c *DashboardController) generateMockOrders(count int) []vo.Order {
	orders := make([]vo.Order, count)
	statuses := []string{"pending", "paid", "shipped", "completed"}
	customers := []string{"张三", "李四", "王五", "赵六", "孙七"}

	for i := 0; i < count; i++ {
		orders[i] = vo.Order{
			ID:            int64(i + 1),
			OrderNo:       generateOrderNo(),
			CustomerName:  customers[rand.Intn(len(customers))],
			CustomerEmail: "customer" + string(rune(i+1)) + "@example.com",
			TotalAmount:   float64(rand.Intn(5000)+100) + 0.99,
			Status:        statuses[rand.Intn(len(statuses))],
			PaymentMethod: "支付宝",
			CreatedAt:     time.Now().Add(-time.Duration(i) * time.Hour),
			UpdatedAt:     time.Now().Add(-time.Duration(i) * time.Minute * 30),
		}
	}

	return orders
}

// generateMockUsers 生成 mock 用户数据
func (c *DashboardController) generateMockUsers(count int) []vo.User {
	users := make([]vo.User, count)
	names := []string{"张伟", "王芳", "李娜", "刘洋", "陈静"}
	roles := []string{"admin", "editor", "viewer"}

	for i := 0; i < count; i++ {
		users[i] = vo.User{
			ID:        int64(i + 1),
			Username:  "user" + string(rune(i+1)),
			Email:     "user" + string(rune(i+1)) + "@example.com",
			RealName:  names[i%len(names)],
			Phone:     "138****" + string(rune(1000+i)),
			Role:      roles[rand.Intn(len(roles))],
			Status:    "active",
			CreatedAt: time.Now().Add(-time.Duration(i) * 24 * time.Hour),
			UpdatedAt: time.Now(),
		}
	}

	return users
}

// generateMockProducts 生成 mock 商品数据
func (c *DashboardController) generateMockProducts(count int) []vo.Product {
	products := make([]vo.Product, count)
	names := []string{"无线蓝牙耳机", "智能手环", "机械键盘", "高清摄像头"}
	categories := []string{"电子产品", "数码配件", "办公用品", "智能设备"}

	for i := 0; i < count; i++ {
		products[i] = vo.Product{
			ID:          int64(i + 1),
			Name:        names[i%len(names)],
			SKU:         "SKU" + string(rune(10000+i)),
			Category:    categories[i%len(categories)],
			Price:       float64(rand.Intn(1000)+50) + 0.99,
			Stock:       rand.Intn(500) + 10,
			Status:      "active",
			Image:       "https://via.placeholder.com/300x200?text=" + names[i%len(names)],
			Description: "这是一个优质的商品",
			CreatedAt:   time.Now().Add(-time.Duration(i) * 7 * 24 * time.Hour),
			UpdatedAt:   time.Now(),
		}
	}

	return products
}

// generateOrderNo 生成订单编号
func generateOrderNo() string {
	return "ORD" + time.Now().Format("20060102150405") + string(rune(rand.Intn(9000)+1000))
}
