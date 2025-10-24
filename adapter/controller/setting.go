// Package controller 系统设置控制器
package controller

import (
	"godash/domain/vo"
	"godash/infra"
	"net/url"

	"github.com/8treenet/freedom"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		// 绑定设置控制器到 /settings 路由
		initiator.BindController("/settings", &SettingController{})
	})
}

// SettingController 系统设置控制器
type SettingController struct {
	Worker  freedom.Worker
	Request *infra.Request
}

// mockSettings 模拟系统设置数据（实际项目中应该从数据库读取）
var mockSettings = vo.SettingsData{
	SiteName:        "HTMX 管理后台",
	SiteDescription: "基于 HTMX、Bulma 和 Alpine.js 的现代化管理后台系统",
	ContactEmail:    "admin@example.com",
	ContactPhone:    "400-123-4567",
	Currency:        "CNY",
	Timezone:        "Asia/Shanghai",
	Language:        "zh-CN",
}

// Get 获取系统设置
// GET /settings
func (c *SettingController) Get() freedom.Result {
	return &infra.ViewResponse{
		Name: "settings/form.html",
		Data: mockSettings,
	}
}

// Post 保存系统设置
// POST /settings
func (c *SettingController) Post() freedom.Result {
	var formData vo.SettingsData
	if err := c.Request.ReadForm(&formData, true); err != nil {
		c.Worker.IrisContext().Header("X-Toast-Message", url.QueryEscape("表单验证失败: "+err.Error()))
		c.Worker.IrisContext().Header("X-Toast-Type", "error")
		return c.Get()
	}

	// 更新设置
	mockSettings = formData

	// 设置成功提示
	c.Worker.IrisContext().Header("X-Toast-Message", url.QueryEscape("设置保存成功"))
	c.Worker.IrisContext().Header("X-Toast-Type", "success")

	// 返回更新后的表单
	return &infra.ViewResponse{
		Name: "settings/form.html",
		Data: mockSettings,
	}
}
