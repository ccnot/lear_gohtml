package tmplfuncs

import (
	"strings"

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
