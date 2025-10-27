package tmplfuncs

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/8treenet/iris/v12/view"
)

// Register 注册模板辅助函数
func Register(engine *view.HTMLEngine) {
	// 字符串函数
	engine.AddFunc("toUpper", strings.ToUpper)
	engine.AddFunc("toLower", strings.ToLower)
	engine.AddFunc("substr", substr)

	// 数学函数
	engine.AddFunc("add", add)
	engine.AddFunc("sub", sub)
	engine.AddFunc("mul", mul)
	engine.AddFunc("div", div)

	// 迭代函数
	engine.AddFunc("iterate", iterate)
	engine.AddFunc("pageRange", pageRange)

	// 字典函数
	engine.AddFunc("dict", dict)

	// JSON 序列化
	engine.AddFunc("toJSON", toJSON)

	// 日期时间格式化
	engine.AddFunc("formatTime", formatTime)
	engine.AddFunc("formatDate", formatDate)
	engine.AddFunc("formatDateTime", formatDateTime)
	engine.AddFunc("formatDateTimeFull", formatDateTimeFull)
}

// substr 截取字符串
func substr(s string, start, length int) string {
	runes := []rune(s)
	if start >= len(runes) {
		return ""
	}
	end := start + length
	if end > len(runes) {
		end = len(runes)
	}
	return string(runes[start:end])
}

// add 加法
func add(a, b int) int {
	return a + b
}

// sub 减法
func sub(a, b int) int {
	return a - b
}

// mul 乘法
func mul(a, b int) int {
	return a * b
}

// div 除法
func div(a, b int) int {
	if b == 0 {
		return 0
	}
	return a / b
}

// iterate 生成从 1 到 n 的序列（用于分页）
func iterate(count int) []int {
	result := make([]int, count)
	for i := 0; i < count; i++ {
		result[i] = i + 1
	}
	return result
}

// pageRange 生成分页显示的页码范围，只包含当前页前后2页和首尾页
func pageRange(currentPage, totalPages int) []int {
	var result []int

	// 总是包含第1页
	if currentPage > 1 {
		result = append(result, 1)
	}

	// 计算显示范围：当前页前2页到后2页
	start := currentPage - 2
	if start < 2 {
		start = 2
	}

	end := currentPage + 2
	if end > totalPages-1 {
		end = totalPages - 1
	}

	// 添加中间页码
	for i := start; i <= end; i++ {
		if i >= 2 && i <= totalPages-1 {
			result = append(result, i)
		}
	}

	// 总是包含最后一页（如果与第1页不同）
	if totalPages > 1 && currentPage != totalPages {
		result = append(result, totalPages)
	}
	return result
}

// dict 创建字典映射，用于在模板中传递多个参数
func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("invalid dict call，需要偶数个参数")
	}
	d := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict keys must be strings")
		}
		d[key] = values[i+1]
	}
	return d, nil
}

// toJSON 将对象序列化为 JSON 字符串，用于在模板中传递数据到 JavaScript
func toJSON(v interface{}) template.JS {
	b, err := json.Marshal(v)
	if err != nil {
		return template.JS("{}")
	}
	return template.JS(b)
}

// formatTime 格式化时间为指定格式
func formatTime(t time.Time, layout string) string {
	return t.Format(layout)
}

// formatDate 格式化日期（年-月-日）
func formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// formatDateTime 格式化日期时间（年-月-日 时:分）
func formatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}

// formatDateTimeFull 格式化完整日期时间（年-月-日 时:分:秒）
func formatDateTimeFull(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
