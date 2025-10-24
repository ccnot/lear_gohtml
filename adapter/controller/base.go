// Package controller 基础控制器
package controller

import (
	"fmt"
	"godash/domain/vo"
	"godash/infra"
	"math"
	"net/url"
	"strings"

	"github.com/8treenet/freedom"
)

// BaseController 基础控制器，提供通用功能
type BaseController struct {
	Worker  freedom.Worker
	Request *infra.Request
}

// PaginationResult 分页结果
type PaginationResult struct {
	Page       int
	PageSize   int
	Total      int64
	TotalPages int
	Start      int
	End        int
}

// SearchHelper 搜索和分页助手
func (c *BaseController) SearchHelper(params vo.SearchParams) (vo.SearchParams, PaginationResult) {
	// 设置默认值
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	return params, PaginationResult{
		Page:     params.Page,
		PageSize: params.PageSize,
	}
}

// Paginate 分页处理
func (c *BaseController) Paginate(items []interface{}, pagination PaginationResult) ([]interface{}, PaginationResult) {
	total := int64(len(items))
	totalPages := int(math.Ceil(float64(total) / float64(pagination.PageSize)))
	start := (pagination.Page - 1) * pagination.PageSize
	end := start + pagination.PageSize

	if end > len(items) {
		end = len(items)
	}
	if start > len(items) {
		start = len(items)
	}

	pagination.Total = total
	pagination.TotalPages = totalPages
	pagination.Start = start
	pagination.End = end

	if start >= len(items) {
		return []interface{}{}, pagination
	}

	return items[start:end], pagination
}

// SetToastMessage 设置 Toast 消息
func (c *BaseController) SetToastMessage(message, toastType string) {
	c.Worker.IrisContext().Header("X-Toast-Message", url.QueryEscape(message))
	c.Worker.IrisContext().Header("X-Toast-Type", toastType)
}

// SetErrorToast 设置错误 Toast
func (c *BaseController) SetErrorToast(message string) {
	c.SetToastMessage(message, "error")
}

// SetSuccessToast 设置成功 Toast
func (c *BaseController) SetSuccessToast(message string) {
	c.SetToastMessage(message, "success")
}

// HandleNotFoundError 处理资源不存在错误
func (c *BaseController) HandleNotFoundError(resource string) freedom.Result {
	c.SetErrorToast(resource + "不存在")
	return &infra.JSONResponse{
		Code:  404,
		Error: fmt.Errorf(resource + "不存在"),
	}
}

// HandleValidationError 处理验证错误
func (c *BaseController) HandleValidationError(err error, viewName string, data interface{}) freedom.Result {
	c.SetErrorToast("表单验证失败: " + err.Error())
	return &infra.ViewResponse{
		Name: viewName,
		Data: data,
	}
}

// SPA Navigation
func (c *BaseController) NavigateTo(path string) {
	c.Worker.IrisContext().Header("HX-Location", `{"path":"`+path+`","target":"#main-container","swap":"innerHTML"}`)
}

// FilterHelper 通用过滤助手
func (c *BaseController) FilterHelper(keyword, searchField string) bool {
	if keyword == "" {
		return true
	}

	lowerKeyword := strings.ToLower(keyword)
	lowerField := strings.ToLower(searchField)

	return strings.Contains(lowerField, lowerKeyword)
}

// CreatePageInfo 创建页面信息
func (c *BaseController) CreatePageInfo(pagination PaginationResult) vo.PageInfo {
	return vo.PageInfo{
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		Total:      pagination.Total,
		TotalPages: pagination.TotalPages,
	}
}
