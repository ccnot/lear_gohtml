package vo

import "time"

// User 用户信息
type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	RealName  string    `json:"real_name"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role"`   // admin, editor, viewer
	Status    string    `json:"status"` // active, inactive, banned
	Avatar    string    `json:"avatar"` // 头像URL
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserListData 用户列表数据
type UserListData struct {
	Users    []User   `json:"users"`
	PageInfo PageInfo `json:"page_info"`
	Query    string   `json:"query"`  // 当前搜索关键词
	Status   string   `json:"status"` // 当前状态筛选
}

// UserFormData 用户表单数据（用于新增/编辑）
type UserFormData struct {
	ID       int64  `json:"id" form:"id"`
	Username string `json:"username" form:"username" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	RealName string `json:"real_name" form:"real_name" validate:"required"`
	Phone    string `json:"phone" form:"phone"`
	Role     string `json:"role" form:"role" validate:"required"`
	Status   string `json:"status" form:"status" validate:"required"`
	Password string `json:"password" form:"password"` // 新增时必填，编辑时可选
}
