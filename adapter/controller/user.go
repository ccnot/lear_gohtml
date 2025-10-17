// Package controller 用户管理控制器
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
		// 绑定用户控制器到 /users 路由
		initiator.BindController("/users", &UserController{})
	})
}

// UserController 用户管理控制器
type UserController struct {
	Worker  freedom.Worker
	Request *infra.Request
}

// mockUsers 模拟用户数据库（使用 map 存储）
var mockUsers = make(map[int64]vo.User)
var userIDCounter int64 = 30

// init 初始化 mock 数据（30条）
func init() {
	names := []string{"张伟", "王芳", "李娜", "刘洋", "陈静", "杨军", "赵敏", "孙涛", "周杰", "吴彦祖",
		"郑爽", "黄晓明", "林志玲", "范冰冰", "李冰冰", "章子怡", "周润发", "刘德华", "张国荣", "梁朝伟",
		"赵本山", "小沈阳", "宋丹丹", "蔡明", "潘长江", "郭德纲", "于谦", "岳云鹏", "孙越", "张云雷"}
	roles := []string{"admin", "editor", "viewer"}
	statuses := []string{"active", "inactive"}

	for i := 0; i < 30; i++ {
		user := vo.User{
			ID:        int64(i + 1),
			Username:  fmt.Sprintf("user%d", i+1),
			Email:     fmt.Sprintf("user%d@example.com", i+1),
			RealName:  names[i],
			Phone:     fmt.Sprintf("138%08d", i+1),
			Role:      roles[i%len(roles)],
			Status:    statuses[i%len(statuses)],
			Avatar:    "/static/images/zxg.jpg",
			CreatedAt: time.Now().Add(-time.Duration(i) * 24 * time.Hour),
			UpdatedAt: time.Now().Add(-time.Duration(i) * time.Hour),
		}
		mockUsers[user.ID] = user
	}
}

// Get 获取用户列表
// GET /users
func (c *UserController) Get() freedom.Result {
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
	filteredUsers := c.filterUsers(params)

	// 分页
	total := int64(len(filteredUsers))
	totalPages := int(math.Ceil(float64(total) / float64(params.PageSize)))
	start := (params.Page - 1) * params.PageSize
	end := start + params.PageSize
	if end > len(filteredUsers) {
		end = len(filteredUsers)
	}
	if start > len(filteredUsers) {
		start = len(filteredUsers)
	}

	pageUsers := filteredUsers[start:end]

	data := vo.UserListData{
		Users: pageUsers,
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
		Name: "users/list.html",
		Data: data,
	}
}

// GetNew 显示新增用户页面
// GET /users/new
func (c *UserController) GetNew() freedom.Result {
	return &infra.ViewResponse{
		Name: "users/new.html",
		Data: nil,
	}
}

// GetBy 获取单个用户（用于编辑）
// GET /users/{id}
func (c *UserController) GetBy(id int64) freedom.Result {
	user := c.findUserByID(id)
	if user == nil {
		c.Worker.IrisContext().Header("X-Toast-Message", url.QueryEscape("用户不存在"))
		c.Worker.IrisContext().Header("X-Toast-Type", "error")
		return c.Get()
	}

	return &infra.ViewResponse{
		Name: "users/edit.html",
		Data: map[string]interface{}{
			"User": user,
		},
	}
}

// Post 创建新用户
// POST /users
func (c *UserController) Post() freedom.Result {
	var formData vo.UserFormData
	if err := c.Request.ReadForm(&formData, true); err != nil {
		// 设置 Toast 消息
		c.Worker.IrisContext().Header("X-Toast-Message", url.QueryEscape("表单验证失败: "+err.Error()))
		c.Worker.IrisContext().Header("X-Toast-Type", "error")
		// 返回带数据的表单，保留用户输入
		return &infra.ViewResponse{
			Name: "users/new.html",
			Data: map[string]interface{}{
				"FormData": formData,
				"Error":    "表单验证失败: " + err.Error(),
			},
		}
	}

	// 检查用户名是否已存在
	for _, u := range mockUsers {
		if u.Username == formData.Username {
			// 设置 Toast 消息
			c.Worker.IrisContext().Header("X-Toast-Message", url.QueryEscape("用户名已存在"))
			c.Worker.IrisContext().Header("X-Toast-Type", "error")
			// 返回带数据的表单，保留用户输入
			return &infra.ViewResponse{
				Name: "users/new.html",
				Data: map[string]interface{}{
					"FormData": formData,
					"Error":    "用户名已存在，请使用其他用户名",
				},
			}
		}
	}

	// 生成新 ID
	userIDCounter++
	newID := userIDCounter

	// 创建新用户
	newUser := vo.User{
		ID:        newID,
		Username:  formData.Username,
		Email:     formData.Email,
		RealName:  formData.RealName,
		Phone:     formData.Phone,
		Role:      formData.Role,
		Status:    formData.Status,
		Avatar:    fmt.Sprintf("https://i.pravatar.cc/150?img=%d", (newID%70)+1),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockUsers[newID] = newUser

	// 设置成功提示并使用 HX-Location 进行 SPA 导航
	// 使用 JSON 格式指定目标容器和其他选项
	c.Worker.IrisContext().Header("HX-Location", `{"path":"/users","target":"#main-container","swap":"innerHTML"}`)
	c.Worker.IrisContext().Header("X-Toast-Message", url.QueryEscape("用户创建成功"))
	c.Worker.IrisContext().Header("X-Toast-Type", "success")

	return &infra.JSONResponse{
		Object: map[string]interface{}{
			"success": true,
			"id":      newID,
		},
	}
}

// PutBy 更新用户
// PUT /users/{id}
func (c *UserController) PutBy(id int64) freedom.Result {
	var formData vo.UserFormData
	if err := c.Request.ReadForm(&formData, true); err != nil {
		// 设置 Toast 消息
		c.Worker.IrisContext().Header("X-Toast-Message", url.QueryEscape("表单验证失败: "+err.Error()))
		c.Worker.IrisContext().Header("X-Toast-Type", "error")
		// 返回带数据的编辑表单，保留用户输入
		user, _ := mockUsers[id]
		return &infra.ViewResponse{
			Name: "users/edit.html",
			Data: map[string]interface{}{
				"User":     user,
				"FormData": formData,
				"Error":    "表单验证失败: " + err.Error(),
			},
		}
	}

	// 查找并更新用户
	user, exists := mockUsers[id]
	if !exists {
		c.Worker.IrisContext().Header("X-Toast-Message", url.QueryEscape("用户不存在"))
		c.Worker.IrisContext().Header("X-Toast-Type", "error")
		c.Worker.IrisContext().Header("HX-Redirect", "/users")
		return &infra.JSONResponse{
			Code:  404,
			Error: fmt.Errorf("用户不存在"),
		}
	}

	// 更新用户信息
	user.Email = formData.Email
	user.RealName = formData.RealName
	user.Phone = formData.Phone
	user.Role = formData.Role
	user.Status = formData.Status
	user.UpdatedAt = time.Now()
	mockUsers[id] = user

	// 设置成功提示并使用 HX-Location 进行 SPA 导航
	// 使用 JSON 格式指定目标容器和其他选项
	c.Worker.IrisContext().Header("HX-Location", `{"path":"/users","target":"#main-container","swap":"innerHTML"}`)
	c.Worker.IrisContext().Header("X-Toast-Message", url.QueryEscape("用户更新成功"))
	c.Worker.IrisContext().Header("X-Toast-Type", "success")

	return &infra.JSONResponse{
		Object: map[string]interface{}{
			"success": true,
			"id":      id,
		},
	}
}

// DeleteBy 删除用户
// DELETE /users/{id}
func (c *UserController) DeleteBy(id int64) freedom.Result {
	// 检查用户是否存在
	if _, exists := mockUsers[id]; !exists {
		c.Worker.IrisContext().Header("X-Toast-Message", url.QueryEscape("用户不存在"))
		c.Worker.IrisContext().Header("X-Toast-Type", "error")
		c.Worker.IrisContext().StatusCode(404)
		return &infra.JSONResponse{
			Code:  404,
			Error: fmt.Errorf("用户不存在"),
		}
	}

	// 删除用户
	delete(mockUsers, id)

	// 设置成功提示和状态码
	c.Worker.IrisContext().Header("X-Toast-Message", url.QueryEscape("用户删除成功"))
	c.Worker.IrisContext().Header("X-Toast-Type", "success")
	c.Worker.IrisContext().StatusCode(200)

	// 返回空响应，让 HTMX 用空内容替换目标元素（实现删除行的效果）
	c.Worker.IrisContext().ContentType("text/html")
	c.Worker.IrisContext().WriteString("")

	return nil
}

// GetRoles 角色管理页面
// GET /users/roles
func (c *UserController) GetRoles() freedom.Result {
	return &infra.ViewResponse{
		Name: "users/roles.html",
		Data: nil,
	}
}

// GetPermissions 权限管理页面
// GET /users/permissions
func (c *UserController) GetPermissions() freedom.Result {
	return &infra.ViewResponse{
		Name: "users/permissions.html",
		Data: nil,
	}
}

// BeforeActivation 配置路由
func (c *UserController) BeforeActivation(b freedom.BeforeActivation) {
	b.Handle("GET", "/new", "GetNew")
	b.Handle("GET", "/roles", "GetRoles")
	b.Handle("GET", "/permissions", "GetPermissions")
	b.Handle("GET", "/{id:int64}", "GetBy")
	b.Handle("PUT", "/{id:int64}", "PutBy")
	b.Handle("DELETE", "/{id:int64}", "DeleteBy")
}

// filterUsers 过滤用户
func (c *UserController) filterUsers(params vo.SearchParams) []vo.User {
	filtered := []vo.User{}

	for _, user := range mockUsers {
		// 搜索过滤
		if params.Keyword != "" {
			keyword := strings.ToLower(params.Keyword)
			if !strings.Contains(strings.ToLower(user.Username), keyword) &&
				!strings.Contains(strings.ToLower(user.Email), keyword) &&
				!strings.Contains(strings.ToLower(user.RealName), keyword) {
				continue
			}
		}

		// 状态过滤
		if params.Status != "" && user.Status != params.Status {
			continue
		}

		filtered = append(filtered, user)
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

// findUserByID 根据 ID 查找用户
func (c *UserController) findUserByID(id int64) *vo.User {
	if user, ok := mockUsers[id]; ok {
		return &user
	}
	return nil
}
