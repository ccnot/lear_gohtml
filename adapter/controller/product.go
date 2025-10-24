// Package controller 商品管理控制器
package controller

import (
	"fmt"
	"godash/domain/vo"
	"godash/infra"
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
	BaseController
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

	// 使用基础控制器的搜索助手，设置商品默认页面大小
	params.PageSize = 12 // 商品默认页面大小
	params, pagination := c.SearchHelper(params)
	filteredProducts := c.filterProducts(params)

	// 转换为 interface{} 进行分页
	products := make([]interface{}, len(filteredProducts))
	for i, product := range filteredProducts {
		products[i] = product
	}

	pagedProducts, pagination := c.Paginate(products, pagination)

	// 转换回商品类型
	result := make([]vo.Product, len(pagedProducts))
	for i, product := range pagedProducts {
		result[i] = product.(vo.Product)
	}

	data := vo.ProductListData{
		Products: result,
		PageInfo: c.CreatePageInfo(pagination),
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
		c.SetErrorToast("商品不存在")
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
		return c.HandleValidationError(err, "products/new.html", map[string]interface{}{
			"FormData": formData,
		})
	}

	// 检查 SKU 是否已存在
	if c.isSKUExists(formData.SKU) {
		c.SetErrorToast("SKU 已存在，请使用其他 SKU")
		return &infra.ViewResponse{
			Name: "products/new.html",
			Data: map[string]interface{}{
				"FormData": formData,
				"Error":    "SKU 已存在，请使用其他 SKU",
			},
		}
	}

	// 创建新商品
	newID := c.generateProductID()
	newProduct := c.createProduct(formData, newID)
	mockProducts[newID] = newProduct

	// 设置成功提示并导航
	c.NavigateTo("/products")
	c.SetSuccessToast("商品创建成功")

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
		product, _ := mockProducts[id]
		return c.HandleValidationError(err, "products/edit.html", map[string]interface{}{
			"Product":  product,
			"FormData": formData,
		})
	}

	// 查找并更新商品
	product, exists := mockProducts[id]
	if !exists {
		c.Worker.IrisContext().Header("HX-Redirect", "/products")
		return c.HandleNotFoundError("商品")
	}

	// 更新商品信息
	c.updateProduct(&product, formData)
	mockProducts[id] = product

	// 设置成功提示并导航
	c.NavigateTo("/products")
	c.SetSuccessToast("商品更新成功")

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
		c.Worker.IrisContext().StatusCode(404)
		return c.HandleNotFoundError("商品")
	}

	// 删除商品
	delete(mockProducts, id)

	// 设置成功提示
	c.SetSuccessToast("商品删除成功")
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
		// 使用基础控制器的过滤助手
		if params.Keyword != "" {
			if !c.FilterHelper(params.Keyword, product.Name) &&
				!c.FilterHelper(params.Keyword, product.SKU) &&
				!c.FilterHelper(params.Keyword, product.Description) {
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
	c.sortProductsByID(filtered)

	return filtered
}

// findProductByID 根据 ID 查找商品
func (c *ProductController) findProductByID(id int64) *vo.Product {
	if product, ok := mockProducts[id]; ok {
		return &product
	}
	return nil
}

// isSKUExists 检查 SKU 是否已存在
func (c *ProductController) isSKUExists(sku string) bool {
	for _, product := range mockProducts {
		if product.SKU == sku {
			return true
		}
	}
	return false
}

// generateProductID 生成新的商品ID
func (c *ProductController) generateProductID() int64 {
	productIDCounter++
	return productIDCounter
}

// createProduct 创建商品对象
func (c *ProductController) createProduct(formData vo.ProductFormData, id int64) vo.Product {
	return vo.Product{
		ID:          id,
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
}

// updateProduct 更新商品信息
func (c *ProductController) updateProduct(product *vo.Product, formData vo.ProductFormData) {
	product.Name = formData.Name
	product.Category = formData.Category
	product.Price = formData.Price
	product.Stock = formData.Stock
	product.Status = formData.Status
	product.Description = formData.Description
	product.UpdatedAt = time.Now()
}

// sortProductsByID 按 ID 降序排序商品
func (c *ProductController) sortProductsByID(products []vo.Product) {
	for i := 0; i < len(products)-1; i++ {
		for j := i + 1; j < len(products); j++ {
			if products[i].ID < products[j].ID {
				products[i], products[j] = products[j], products[i]
			}
		}
	}
}
