// Package controller 用户管理控制器
package controller

import (
	"fmt"
	"gohtml/domain/vo"
	"gohtml/infra"
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
	BaseController
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

	// 使用基础控制器的搜索助手
	params, pagination := c.SearchHelper(params)
	filteredUsers := c.filterUsers(params)

	// 转换为 interface{} 进行分页
	users := make([]interface{}, len(filteredUsers))
	for i, user := range filteredUsers {
		users[i] = user
	}

	pagedUsers, pagination := c.Paginate(users, pagination)

	// 转换回用户类型
	result := make([]vo.User, len(pagedUsers))
	for i, user := range pagedUsers {
		result[i] = user.(vo.User)
	}

	data := vo.UserListData{
		Users:    result,
		PageInfo: c.CreatePageInfo(pagination),
		Query:    params.Keyword,
		Status:   params.Status,
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
		c.SetErrorToast("用户不存在")
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
		return c.HandleValidationError(err, "users/new.html", formData)
	}

	// 检查用户名是否已存在
	if c.isUsernameExists(formData.Username) {
		c.SetErrorToast("用户名已存在")
		return &infra.ViewResponse{
			Name: "users/new.html",
			Data: formData,
		}
	}

	// 创建新用户
	newID := c.generateUserID()
	newUser := c.createUser(formData, newID)
	mockUsers[newID] = newUser

	// 设置成功提示并导航
	c.NavigateTo("/users")
	c.SetSuccessToast("用户创建成功")

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
		user, _ := mockUsers[id]
		return c.HandleValidationError(err, "users/edit.html", map[string]interface{}{
			"User":     user,
			"FormData": formData,
		})
	}

	// 查找并更新用户
	user, exists := mockUsers[id]
	if !exists {
		c.Worker.IrisContext().Header("HX-Redirect", "/users")
		return c.HandleNotFoundError("用户")
	}

	// 更新用户信息
	c.updateUser(&user, formData)
	mockUsers[id] = user

	// 设置成功提示并导航
	c.NavigateTo("/users")
	c.SetSuccessToast("用户更新成功")

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
		c.Worker.IrisContext().StatusCode(404)
		return c.HandleNotFoundError("用户")
	}

	// 删除用户
	delete(mockUsers, id)

	// 设置成功提示
	c.SetSuccessToast("用户删除成功")
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
		// 使用基础控制器的过滤助手
		if params.Keyword != "" {
			if !c.FilterHelper(params.Keyword, user.Username) &&
				!c.FilterHelper(params.Keyword, user.Email) &&
				!c.FilterHelper(params.Keyword, user.RealName) {
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
	c.sortUsersByID(filtered)

	return filtered
}

// findUserByID 根据 ID 查找用户
func (c *UserController) findUserByID(id int64) *vo.User {
	if user, ok := mockUsers[id]; ok {
		return &user
	}
	return nil
}

// isUsernameExists 检查用户名是否已存在
func (c *UserController) isUsernameExists(username string) bool {
	for _, user := range mockUsers {
		if user.Username == username {
			return true
		}
	}
	return false
}

// generateUserID 生成新的用户ID
func (c *UserController) generateUserID() int64 {
	userIDCounter++
	return userIDCounter
}

// createUser 创建用户对象
func (c *UserController) createUser(formData vo.UserFormData, id int64) vo.User {
	return vo.User{
		ID:        id,
		Username:  formData.Username,
		Email:     formData.Email,
		RealName:  formData.RealName,
		Phone:     formData.Phone,
		Role:      formData.Role,
		Status:    formData.Status,
		Avatar:    fmt.Sprintf("https://i.pravatar.cc/150?img=%d", (id%70)+1),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// updateUser 更新用户信息
func (c *UserController) updateUser(user *vo.User, formData vo.UserFormData) {
	user.Email = formData.Email
	user.RealName = formData.RealName
	user.Phone = formData.Phone
	user.Role = formData.Role
	user.Status = formData.Status
	user.UpdatedAt = time.Now()
}

// sortUsersByID 按 ID 降序排序用户
func (c *UserController) sortUsersByID(users []vo.User) {
	for i := 0; i < len(users)-1; i++ {
		for j := i + 1; j < len(users); j++ {
			if users[i].ID < users[j].ID {
				users[i], users[j] = users[j], users[i]
			}
		}
	}
}
