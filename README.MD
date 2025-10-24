# HTMX 管理后台学习项目

一个基于 HTMX v2.0.7、Bulma v1.0.4 和 Alpine.js v3.15.0 构建的现代化 SPA 风格管理后台系统。

## 项目简介

这是一个完整的管理后台示例项目，专为学习 HTMX 的最佳实践而设计。项目采用 SPA 架构，使用 HTMX 处理部分页面更新，Bulma 提供美观的 UI，Alpine.js 处理纯前端交互。

## 核心功能

✅ **仪表盘** - 统计数据展示、自动刷新、最近数据  
✅ **用户管理** - 完整 CRUD、搜索筛选、分页  
✅ **商品管理** - 卡片展示、实时搜索、分类筛选  
✅ **订单管理** - 订单列表、详情查看、状态更新  
✅ **系统设置** - Tab 切换、表单提交、配置管理  

## 技术栈

- **后端**: Go + Freedom Framework
- **前端**: HTMX v2.0.7 + Bulma v1.0.4 + Alpine.js v3.15.0
- **架构**: SPA（单页应用）
- **设计**: 现代、简洁、美观（桌面端优先）

## 快速开始

### 1. 安装依赖
```bash
go mod download
```

### 2. 运行项目
```bash
go run main.go
```

### 3. 访问应用
打开浏览器访问: http://localhost:8000

## HTMX 特性展示

本项目全面展示了 HTMX 的核心特性：

- ✅ `hx-get/post/put/delete` - 基本 HTTP 请求
- ✅ `hx-target` - 精确控制更新目标
- ✅ `hx-swap` - 多种内容替换方式
- ✅ `hx-trigger` - 自定义触发事件（实时搜索、定时刷新）
- ✅ `hx-indicator` - 加载状态指示
- ✅ `hx-push-url` - 浏览器历史管理
- ✅ `hx-vals` - 额外请求参数
- ✅ `hx-include` - 包含其他表单字段
- ✅ `hx-swap-oob` - Out of Band 更新
- ✅ Toast 通知系统

## 项目结构

```
godash/
├── adapter/controller/    # 控制器层（dashboard, user, product, order, setting）
├── domain/vo/            # 值对象（数据结构定义）
├── web/
│   ├── static/css/       # 自定义样式
│   ├── tmplfuncs/        # 模板辅助函数
│   └── views/            # HTML 视图模板
│       ├── layout.html   # SPA 主布局
│       ├── components/   # 通用组件（sidebar, header）
│       ├── dashboard/    # 仪表盘视图
│       ├── users/        # 用户管理视图
│       ├── products/     # 商品管理视图
│       ├── orders/       # 订单管理视图
│       └── settings/     # 系统设置视图
├── main.go              # 应用入口
├── README.MD            # 项目说明
└── PROJECT_GUIDE.md     # 详细学习指南
```

## 学习要点

### SPA 架构
- 主布局只加载一次
- 内容区域动态更新
- 浏览器历史管理
- 无页面刷新的流畅体验

### 部分页面更新
- 服务器返回 HTML 片段
- 只更新需要的部分
- 减少数据传输
- 提高响应速度

### 最佳实践
- 语义化 HTML
- RESTful 路由设计
- 渐进增强
- 详细的中文注释
- 模块化代码结构

## 代码示例

### 实时搜索
```html
<input type="search" 
       hx-get="/users"
       hx-trigger="keyup changed delay:500ms"
       hx-target="#user-table"
       hx-swap="outerHTML">
```

### 定时刷新
```html
<div hx-get="/dashboard/stats" 
     hx-trigger="every 30s">
  <!-- 统计数据 -->
</div>
```

### 单行更新
```html
<form hx-put="/users/123"
      hx-target="#user-row-123"
      hx-swap="outerHTML">
  <!-- 表单 -->
</form>
```

## 注意事项

⚠️ **这是一个学习项目**
- 所有数据都是内存中的 mock 数据
- 没有实现真实的数据库持久化
- 没有用户认证和授权
- 重启应用后数据会丢失

## 详细文档

查看 [PROJECT_GUIDE.md](./PROJECT_GUIDE.md) 获取更详细的说明，包括：
- 完整的功能说明
- HTMX 特性详解
- Alpine.js 应用场景
- 代码示例和最佳实践
- 扩展建议

## 学习资源

- [HTMX 官方文档](https://htmx.org/)
- [Bulma CSS 框架](https://bulma.io/)
- [Alpine.js 文档](https://alpinejs.dev/)
- [Freedom 框架](https://github.com/8treenet/freedom)

## 特点

✨ **现代化设计** - 简洁美观的 UI  
🚀 **高性能** - 部分页面更新  
📱 **桌面优先** - 企业级应用  
🎯 **最佳实践** - HTMX 正确用法  
📝 **详细注释** - 易于理解和学习  
🔧 **易于扩展** - 模块化设计  

## 许可证

本项目仅用于学习目的。

---

**Happy Learning! 🎉**
