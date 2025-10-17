package vo

// DashboardStats 仪表盘统计数据
type DashboardStats struct {
	TotalUsers    int64   `json:"total_users"`    // 总用户数
	TotalProducts int64   `json:"total_products"` // 总商品数
	TotalOrders   int64   `json:"total_orders"`   // 总订单数
	TotalRevenue  float64 `json:"total_revenue"`  // 总收入
	ActiveUsers   int64   `json:"active_users"`   // 活跃用户
	PendingOrders int64   `json:"pending_orders"` // 待处理订单
	LowStock      int64   `json:"low_stock"`      // 低库存商品
	TodayOrders   int64   `json:"today_orders"`   // 今日订单
}

// DashboardData 仪表盘数据
type DashboardData struct {
	Stats        DashboardStats `json:"stats"`
	RecentOrders []Order        `json:"recent_orders"` // 最近订单
	RecentUsers  []User         `json:"recent_users"`  // 最近用户
	TopProducts  []Product      `json:"top_products"`  // 热门商品
}

// SettingsData 系统设置数据（扩展版本）
type SettingsData struct {
	// 基本信息
	SiteName        string `json:"site_name" form:"site_name" validate:"required"`
	SiteDescription string `json:"site_description" form:"site_description"`
	SiteLogo        string `json:"site_logo" form:"site_logo"`
	SiteURL         string `json:"site_url" form:"site_url"`

	// 联系方式
	ContactEmail   string `json:"contact_email" form:"contact_email" validate:"required,email"`
	ContactPhone   string `json:"contact_phone" form:"contact_phone"`
	ContactAddress string `json:"contact_address" form:"contact_address"`

	// 区域设置
	Currency   string `json:"currency" form:"currency" validate:"required"`
	Timezone   string `json:"timezone" form:"timezone" validate:"required"`
	Language   string `json:"language" form:"language" validate:"required"`
	DateFormat string `json:"date_format" form:"date_format"`

	// 主题设置
	ThemeColor   string `json:"theme_color" form:"theme_color"`
	SidebarColor string `json:"sidebar_color" form:"sidebar_color"`

	// 功能开关
	EnableRegistration  bool `json:"enable_registration" form:"enable_registration"`
	EnableComments      bool `json:"enable_comments" form:"enable_comments"`
	EnableNotifications bool `json:"enable_notifications" form:"enable_notifications"`
	MaintenanceMode     bool `json:"maintenance_mode" form:"maintenance_mode"`

	// 其他
	ItemsPerPage   int `json:"items_per_page" form:"items_per_page"`
	SessionTimeout int `json:"session_timeout" form:"session_timeout"`
}

// SettingsStats 设置统计信息
type SettingsStats struct {
	LastSaveTime  string `json:"last_save_time"`
	TotalSaves    int    `json:"total_saves"`
	ConfigVersion string `json:"config_version"`
	SystemUptime  string `json:"system_uptime"`
}
