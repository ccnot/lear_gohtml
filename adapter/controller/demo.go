// Package controller 特性演示控制器
package controller

import (
	"fmt"
	"gohtml/infra"
	"math/rand"
	"net/url"
	"time"

	"github.com/8treenet/freedom"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		// 绑定演示控制器到 /demo 路由
		initiator.BindController("/demo", &DemoController{})
	})
}

// DemoController 特性演示控制器
type DemoController struct {
	Worker  freedom.Worker
	Request *infra.Request
}

// mockDemoData 模拟演示数据
var mockDemoData = map[string]interface{}{
	"counter":       0,
	"theme":         "light",
	"users":         []string{"Alice", "Bob", "Charlie", "David", "Eve"},
	"lastSave":      time.Now().Format("2006-01-02 15:04:05"),
	"autoSaveCount": 0,
}

// BeforeActivation 配置路由
func (c *DemoController) BeforeActivation(b freedom.BeforeActivation) {
	b.Handle("GET", "/", "Get")
	b.Handle("POST", "/autosave", "PostAutosave")
	b.Handle("GET", "/stats", "GetStats")
	b.Handle("POST", "/theme", "PostTheme")
	b.Handle("GET", "/users", "GetUsers")
	b.Handle("POST", "/validate", "PostValidate")
	b.Handle("GET", "/poll", "GetPoll")
	b.Handle("POST", "/slow", "PostSlow")
}

// Get 获取演示页面
// GET /demo
func (c *DemoController) Get() freedom.Result {
	return &infra.ViewResponse{
		Name: "demo/features.html",
		Data: mockDemoData,
	}
}

// PostAutosave 自动保存演示
// POST /demo/autosave
func (c *DemoController) PostAutosave() freedom.Result {
	var data struct {
		Field string `form:"field"`
		Value string `form:"value"`
	}

	if err := c.Request.ReadForm(&data, true); err != nil {
		c.Worker.IrisContext().Header("X-Toast-Message", url.QueryEscape("保存失败"))
		c.Worker.IrisContext().Header("X-Toast-Type", "error")
		c.Worker.IrisContext().WriteString("保存失败")
		return nil
	}

	// 模拟保存
	mockDemoData[data.Field] = data.Value
	mockDemoData["lastSave"] = time.Now().Format("2006-01-02 15:04:05")
	count := mockDemoData["autoSaveCount"].(int)
	mockDemoData["autoSaveCount"] = count + 1

	// 返回保存状态 HTML
	c.Worker.IrisContext().ContentType("text/html")
	html := fmt.Sprintf(`
		<div class="notification is-success is-light" style="padding: 0.5rem 1rem; margin: 0;">
			<span class="icon"><i class="fas fa-check"></i></span>
			<span>已自动保存于 %s</span>
		</div>
	`, mockDemoData["lastSave"])
	c.Worker.IrisContext().WriteString(html)
	return nil
}

// GetStats 获取统计信息（用于轮询演示）
// GET /demo/stats
func (c *DemoController) GetStats() freedom.Result {
	stats := map[string]interface{}{
		"time":          time.Now().Format("15:04:05"),
		"random":        rand.Intn(100),
		"autoSaveCount": mockDemoData["autoSaveCount"],
		"lastSave":      mockDemoData["lastSave"],
	}

	c.Worker.IrisContext().ContentType("text/html")
	html := fmt.Sprintf(`
		<div style="display: grid; grid-template-columns: repeat(2, 1fr); gap: 1rem;">
			<div class="box has-text-centered">
				<p class="heading">当前时间</p>
				<p class="title is-5">%s</p>
			</div>
			<div class="box has-text-centered">
				<p class="heading">随机数</p>
				<p class="title is-5">%d</p>
			</div>
			<div class="box has-text-centered">
				<p class="heading">自动保存次数</p>
				<p class="title is-5">%d</p>
			</div>
			<div class="box has-text-centered">
				<p class="heading">最后保存</p>
				<p class="title is-6" style="font-size: 0.9rem;">%s</p>
			</div>
		</div>
	`, stats["time"], stats["random"], stats["autoSaveCount"], stats["lastSave"])
	c.Worker.IrisContext().WriteString(html)
	return nil
}

// PostTheme 主题切换
// POST /demo/theme
func (c *DemoController) PostTheme() freedom.Result {
	var data struct {
		Theme string `form:"theme"`
	}

	c.Request.ReadForm(&data, true)
	mockDemoData["theme"] = data.Theme

	c.Worker.IrisContext().Header("X-Toast-Message", url.QueryEscape("主题已切换"))
	c.Worker.IrisContext().Header("X-Toast-Type", "success")

	c.Worker.IrisContext().ContentType("text/html")
	c.Worker.IrisContext().WriteString(fmt.Sprintf("当前主题: %s", data.Theme))
	return nil
}

// GetUsers 获取用户列表（用于 hx-select 演示）
// GET /demo/users
func (c *DemoController) GetUsers() freedom.Result {
	users := mockDemoData["users"].([]string)

	c.Worker.IrisContext().ContentType("text/html")
	html := `<div id="user-list">`
	for _, user := range users {
		html += fmt.Sprintf(`
			<div class="box" style="padding: 1rem; margin-bottom: 0.5rem;">
				<span class="icon-text">
					<span class="icon has-text-info">
						<i class="fas fa-user"></i>
					</span>
					<span>%s</span>
				</span>
			</div>
		`, user)
	}
	html += `</div>`

	c.Worker.IrisContext().WriteString(html)
	return nil
}

// PostValidate 表单验证演示
// POST /demo/validate
func (c *DemoController) PostValidate() freedom.Result {
	var data struct {
		Email    string `form:"email"`
		Password string `form:"password"`
	}

	c.Request.ReadForm(&data, true)

	// 模拟验证
	if data.Email == "" || data.Password == "" {
		c.Worker.IrisContext().Header("X-Toast-Message", url.QueryEscape("请填写所有字段"))
		c.Worker.IrisContext().Header("X-Toast-Type", "error")
		c.Worker.IrisContext().ContentType("text/html")
		c.Worker.IrisContext().WriteString(`
			<div class="notification is-danger">
				<strong>验证失败：</strong>请填写所有必填字段
			</div>
		`)
		return nil
	}

	c.Worker.IrisContext().Header("X-Toast-Message", url.QueryEscape("验证通过！"))
	c.Worker.IrisContext().Header("X-Toast-Type", "success")
	c.Worker.IrisContext().ContentType("text/html")
	c.Worker.IrisContext().WriteString(`
		<div class="notification is-success">
			<strong>验证成功！</strong>表单数据已通过所有验证
		</div>
	`)
	return nil
}

// GetPoll 轮询数据
// GET /demo/poll
func (c *DemoController) GetPoll() freedom.Result {
	// 模拟实时数据
	timestamp := time.Now().Format("15:04:05")
	value := rand.Intn(100)

	c.Worker.IrisContext().ContentType("text/html")
	html := fmt.Sprintf(`
		<div class="notification is-info is-light">
			<strong>%s</strong> - 实时数值: <strong>%d</strong>
		</div>
	`, timestamp, value)
	c.Worker.IrisContext().WriteString(html)
	return nil
}

// PostSlow 慢速响应（用于测试 loading 状态）
// POST /demo/slow
func (c *DemoController) PostSlow() freedom.Result {
	// 模拟慢速响应
	time.Sleep(2 * time.Second)

	c.Worker.IrisContext().Header("X-Toast-Message", url.QueryEscape("慢速请求完成"))
	c.Worker.IrisContext().Header("X-Toast-Type", "success")
	c.Worker.IrisContext().ContentType("text/html")
	c.Worker.IrisContext().WriteString(`
		<div class="notification is-success">
			<strong>完成！</strong>慢速请求已处理完毕（2秒延迟）
		</div>
	`)
	return nil
}
