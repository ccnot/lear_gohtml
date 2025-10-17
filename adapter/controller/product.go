// Package controller 商品管理控制器
package controller

import (
	"fmt"
	"gohtml/domain/vo"
	"gohtml/infra"
	"math"
	"net/url"
	"strings"
	"time"

	"github.com/8treenet/freedom"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		// 绑定商品控制器到 /products 路由
		initiator.BindController("/products", &ProductController{})
	})
}

// ProductController 商品管理控制器
type ProductController struct {
	Worker  freedom.Worker
	Request *infra.Request
}

// mockProducts 模拟商品数据库（使用 map 存储）
var mockProducts = make(map[int64]vo.Product)
var productIDCounter int64 = 30

// init 初始化商品 mock 数据（30条）
func init() {
	names := []string{
		"无线蓝牙耳机", "智能手环", "机械键盘", "高清摄像头", "笔记本电脑",
		"显示器", "鼠标垫", "USB充电器", "移动硬盘", "路由器",
		"智能音箱", "平板电脑", "游戏手柄", "麦克风", "电脑椅",
		"台灯", "手机支架", "数据线", "蓝牙音箱", "投影仪",
		"扫描仪", "打印机", "绘图板", "读卡器", "散热器",
		"电源适配器", "网线", "HDMI线", "耳机架", "桌面支架",
	}
	categories := []string{"电子产品", "数码配件", "办公用品", "智能设备", "电脑配件"}
	statuses := []string{"active", "inactive", "out_of_stock"}

	for i := 0; i < 30; i++ {
		product := vo.Product{
			ID:          int64(i + 1),
			Name:        names[i],
			SKU:         fmt.Sprintf("SKU%05d", i+1),
			Category:    categories[i%len(categories)],
			Price:       float64((i+1)*50) + 99.99,
			Stock:       (i+1)*10 - (i % 3 * 5),
			Status:      statuses[i%len(statuses)],
			Image:       fmt.Sprintf("https://via.placeholder.com/300x200?text=%s", names[i]),
			Description: fmt.Sprintf("这是一款优质的%s，性能卓越，品质保证。", names[i]),
			CreatedAt:   time.Now().Add(-time.Duration(i) * 24 * time.Hour),
			UpdatedAt:   time.Now().Add(-time.Duration(i) * time.Hour),
		}
		mockProducts[product.ID] = product
	}
}

// Get 获取商品列表
// GET /products
func (c *ProductController) Get() freedom.Result {
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
		params.PageSize = 12
	}

	// 过滤和搜索
	filteredProducts := c.filterProducts(params)

	// 分页
	total := int64(len(filteredProducts))
	totalPages := int(math.Ceil(float64(total) / float64(params.PageSize)))
	start := (params.Page - 1) * params.PageSize
	end := start + params.PageSize
	if end > len(filteredProducts) {
		end = len(filteredProducts)
	}
	if start > len(filteredProducts) {
		start = len(filteredProducts)
	}

	pageProducts := filteredProducts[start:end]

	data := vo.ProductListData{
		Products: pageProducts,
		PageInfo: vo.PageInfo{
			Page:       params.Page,
			PageSize:   params.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
		Query:    params.Keyword,
		Category: params.Status, // 这里复用 Status 字段作为分类筛选
	}

	return &infra.ViewResponse{
		Name: "products/list.html",
		Data: data,
	}
}

// GetNew 显示新增商品页面
// GET /products/new
func (c *ProductController) GetNew() freedom.Result {
	return &infra.ViewResponse{
		Name: "products/new.html",
		Data: nil,
	}
}

// GetBy 获取单个商品（用于编辑）
// GET /products/{id}
func (c *ProductController) GetBy(id int64) freedom.Result {
	product := c.findProductByID(id)
	if product == nil {
		c.Worker.IrisContext().Header("X-Toast-Message", url.QueryEscape("商品不存在"))
		c.Worker.IrisContext().Header("X-Toast-Type", "error")
		return c.Get()
	}

	return &infra.ViewResponse{
		Name: "products/edit.html",
		Data: map[string]interface{}{
			"Product": product,
		},
	}
}

// Post 创建新商品
// POST /products
func (c *ProductController) Post() freedom.Result {
	var formData vo.ProductFormData
	if err := c.Request.ReadForm(&formData, true); err != nil {
		// 设置 Toast 消息
		c.Worker.IrisContext().Header("X-Toast-Message", url.QueryEscape("表单验证失败: "+err.Error()))
		c.Worker.IrisContext().Header("X-Toast-Type", "error")
		// 返回带数据的表单，保留用户输入
		return &infra.ViewResponse{
			Name: "products/new.html",
			Data: map[string]interface{}{
				"FormData": formData,
				"Error":    "表单验证失败: " + err.Error(),
			},
		}
	}

	// 检查 SKU 是否已存在
	for _, p := range mockProducts {
		if p.SKU == formData.SKU {
			// 设置 Toast 消息
			c.Worker.IrisContext().Header("X-Toast-Message", url.QueryEscape("SKU 已存在"))
			c.Worker.IrisContext().Header("X-Toast-Type", "error")
			// 返回带数据的表单，保留用户输入
			return &infra.ViewResponse{
				Name: "products/new.html",
				Data: map[string]interface{}{
					"FormData": formData,
					"Error":    "SKU 已存在，请使用其他 SKU",
				},
			}
		}
	}

	// 生成新 ID
	productIDCounter++
	newID := productIDCounter

	// 创建新商品
	newProduct := vo.Product{
		ID:          newID,
		Name:        formData.Name,
		SKU:         formData.SKU,
		Category:    formData.Category,
		Price:       formData.Price,
		Stock:       formData.Stock,
		Status:      formData.Status,
		Image:       fmt.Sprintf("https://via.placeholder.com/300x200?text=%s", formData.Name),
		Description: formData.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockProducts[newID] = newProduct

	// 设置成功提示并使用 HX-Location 进行 SPA 导航
	// 使用 JSON 格式指定目标容器和其他选项
	c.Worker.IrisContext().Header("HX-Location", `{"path":"/products","target":"#main-container","swap":"innerHTML"}`)
	c.Worker.IrisContext().Header("X-Toast-Message", url.QueryEscape("商品创建成功"))
	c.Worker.IrisContext().Header("X-Toast-Type", "success")

	return &infra.JSONResponse{
		Object: map[string]interface{}{
			"success": true,
			"id":      newID,
		},
	}
}

// PutBy 更新商品
// PUT /products/{id}
func (c *ProductController) PutBy(id int64) freedom.Result {
	var formData vo.ProductFormData
	if err := c.Request.ReadForm(&formData, true); err != nil {
		// 设置 Toast 消息
		c.Worker.IrisContext().Header("X-Toast-Message", url.QueryEscape("表单验证失败: "+err.Error()))
		c.Worker.IrisContext().Header("X-Toast-Type", "error")
		// 返回带数据的编辑表单，保留用户输入
		product, _ := mockProducts[id]
		return &infra.ViewResponse{
			Name: "products/edit.html",
			Data: map[string]interface{}{
				"Product":  product,
				"FormData": formData,
				"Error":    "表单验证失败: " + err.Error(),
			},
		}
	}

	// 查找并更新商品
	product, exists := mockProducts[id]
	if !exists {
		c.Worker.IrisContext().Header("X-Toast-Message", url.QueryEscape("商品不存在"))
		c.Worker.IrisContext().Header("X-Toast-Type", "error")
		c.Worker.IrisContext().Header("HX-Redirect", "/products")
		return &infra.JSONResponse{
			Code:  404,
			Error: fmt.Errorf("商品不存在"),
		}
	}

	// 更新商品信息
	product.Name = formData.Name
	product.Category = formData.Category
	product.Price = formData.Price
	product.Stock = formData.Stock
	product.Status = formData.Status
	product.Description = formData.Description
	product.UpdatedAt = time.Now()
	mockProducts[id] = product

	// 设置成功提示并使用 HX-Location 进行 SPA 导航
	// 使用 JSON 格式指定目标容器和其他选项
	c.Worker.IrisContext().Header("HX-Location", `{"path":"/products","target":"#main-container","swap":"innerHTML"}`)
	c.Worker.IrisContext().Header("X-Toast-Message", url.QueryEscape("商品更新成功"))
	c.Worker.IrisContext().Header("X-Toast-Type", "success")

	return &infra.JSONResponse{
		Object: map[string]interface{}{
			"success": true,
			"id":      id,
		},
	}
}

// DeleteBy 删除商品
// DELETE /products/{id}
func (c *ProductController) DeleteBy(id int64) freedom.Result {
	// 检查商品是否存在
	if _, exists := mockProducts[id]; !exists {
		c.Worker.IrisContext().Header("X-Toast-Message", url.QueryEscape("商品不存在"))
		c.Worker.IrisContext().Header("X-Toast-Type", "error")
		c.Worker.IrisContext().StatusCode(404)
		return &infra.JSONResponse{
			Code:  404,
			Error: fmt.Errorf("商品不存在"),
		}
	}

	// 删除商品
	delete(mockProducts, id)

	// 设置成功提示和状态码
	c.Worker.IrisContext().Header("X-Toast-Message", url.QueryEscape("商品删除成功"))
	c.Worker.IrisContext().Header("X-Toast-Type", "success")
	c.Worker.IrisContext().StatusCode(200)

	// 返回空响应，让 HTMX 用空内容替换目标元素（实现删除卡片的效果）
	c.Worker.IrisContext().ContentType("text/html")
	c.Worker.IrisContext().WriteString("")

	return nil
}

// BeforeActivation 配置路由
func (c *ProductController) BeforeActivation(b freedom.BeforeActivation) {
	b.Handle("GET", "/new", "GetNew")
	b.Handle("GET", "/{id:int64}", "GetBy")
	b.Handle("PUT", "/{id:int64}", "PutBy")
	b.Handle("DELETE", "/{id:int64}", "DeleteBy")
}

// filterProducts 过滤商品
func (c *ProductController) filterProducts(params vo.SearchParams) []vo.Product {
	filtered := []vo.Product{}

	for _, product := range mockProducts {
		// 搜索过滤
		if params.Keyword != "" {
			keyword := strings.ToLower(params.Keyword)
			if !strings.Contains(strings.ToLower(product.Name), keyword) &&
				!strings.Contains(strings.ToLower(product.SKU), keyword) &&
				!strings.Contains(strings.ToLower(product.Description), keyword) {
				continue
			}
		}

		// 分类过滤（复用 Status 字段）
		if params.Status != "" && product.Category != params.Status {
			continue
		}

		filtered = append(filtered, product)
	}

	// 按 ID 降序排序（最新的在前面）
	for i := 0; i < len(filtered)-1; i++ {
		for j := i + 1; j < len(filtered); j++ {
			if filtered[i].ID < filtered[j].ID {
				filtered[i], filtered[j] = filtered[j], filtered[i]
			}
		}
	}

	return filtered
}

// findProductByID 根据 ID 查找商品
func (c *ProductController) findProductByID(id int64) *vo.Product {
	if product, ok := mockProducts[id]; ok {
		return &product
	}
	return nil
}
