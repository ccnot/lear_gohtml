# HTMX 管理后台学习项目

一个基于 HTMX v2.0.7、DaisyUI 5.3.7 和 Alpine.js v3.15.0 构建的现代化 SPA 风格管理后台系统。

## 项目简介

这是一个完整的管理后台示例项目，专为学习 HTMX 的最佳实践而设计。项目采用 SPA 架构，使用 HTMX 处理部分页面更新，DaisyUI 提供美观的 UI，Alpine.js 处理纯前端交互。

## 核心功能

✅ **仪表盘** - 统计数据展示、自动刷新、最近数据
✅ **用户管理** - 完整 CRUD、搜索筛选、分页
✅ **商品管理** - 卡片展示、实时搜索、分类筛选
✅ **订单管理** - 订单列表、详情查看、状态更新
✅ **系统设置** - Tab 切换、表单提交、配置管理

## 技术栈

- **后端**: Go + Freedom Framework v1.9.7
- **前端**: HTMX v2.0.7 + DaisyUI 5.3.7 + Tailwind CSS 4 + Alpine.js v3.15.0
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

## HTMX v2.0.7 特性展示

本项目全面展示了 HTMX v2.0.7 的核心特性：

- ✅ `hx-get/post/put/delete` - 基本 HTTP 请求
- ✅ `hx-target` - 精确控制更新目标
- ✅ `hx-swap` - 多种内容替换方式
- ✅ `hx-trigger` - 自定义触发事件（实时搜索、定时刷新）
- ✅ `hx-indicator` - 加载状态指示
- ✅ `hx-push-url` - 浏览器历史管理
- ✅ `hx-vals` - 额外请求参数
- ✅ `hx-swap-oob` - Out of Band 更新
- ✅ `hx-on:*` - 事件处理
- ✅ Toast 通知系统
- ✅ 自定义确认对话框
- ✅ View Transitions 支持

## 项目结构

```
godash/
├── adapter/controller/    # 控制器层（Freedom Framework）
├── domain/vo/            # 值对象（数据结构定义）
├── web/
│   ├── static/css/       # 精简的自定义样式
│   ├── static/js/        # 前端交互逻辑（Alpine.js + HTMX）
│   ├── tmplfuncs/        # 模板辅助函数
│   └── views/            # HTML 视图模板
│       ├── layout.html   # SPA 主布局（重构版）
│       ├── components/   # 通用组件（sidebar, header）
│       ├── dashboard/    # 仪表盘视图
│       ├── users/        # 用户管理视图
│       ├── products/     # 商品管理视图
│       ├── orders/       # 订单管理视图
│       └── settings/     # 系统设置视图
├── main.go              # 应用入口
├── config/              # 配置文件
├── go.mod               # Go 依赖管理
├── go.sum               # Go 依赖锁定
└── README.MD            # 项目说明
```

## 学习要点

### SPA 架构（重构版）
- 主布局只加载一次，内容区域动态更新
- 统一的 HTMX v2.0.7 请求处理
- 优化的 Alpine.js 状态管理
- 响应式设计支持

### 现代前端技术栈
- **DaisyUI 5.3.7**: 组件库，提供美观的 UI 组件
- **Tailwind CSS 4**: 原子化 CSS 框架
- **Alpine.js v3.15.0**: 轻量级前端状态管理
- **HTMX v2.0.7**: 无需 JavaScript 的 AJAX 交互

### 部分页面更新
- 服务器返回 HTML 片段，只更新需要的部分
- `hx-swap-oob` 支持，一次请求更新多个区域
- 智能的加载状态管理
- 优化的用户体验

### 重构后的改进
- **代码简化**: 移除冗余代码，提高可维护性
- **版本统一**: 统一使用最新稳定版本
- **性能优化**: CDN 资源加载，减少本地依赖
- **最佳实践**: 遵循 HTMX v2.0.7 和 DaisyUI 5.3.7 最佳实践

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

## 重构亮点

### 🎯 技术栈现代化
- 从 Bulma 迁移到 DaisyUI 5.3.7 + Tailwind CSS 4
- 统一使用 CDN 资源，减少项目体积
- 优化了 Tailwind 配置，支持自定义主题

### 📦 代码精简
- **CSS**: 从 414 行精简到 76 行，减少 82%
- **JavaScript**: 从 515 行优化到 324 行，减少 37%
- 移除冗余的本地依赖文件

### 🚀 性能优化
- 统一使用 CDN 加载，提高加载速度
- 优化了 Alpine.js 组件结构
- 改进了 HTMX 事件处理机制

### 🛠️ 开发体验
- 统一的版本管理
- 清晰的代码结构
- 完善的错误处理

## 学习资源

- [HTMX v2.0.7 官方文档](https://htmx.org/)
- [DaisyUI 5.3.7 文档](https://daisyui.com/)
- [Tailwind CSS 4 文档](https://tailwindcss.com/)
- [Alpine.js v3.15.0 文档](https://alpinejs.dev/)
- [Freedom 框架](https://github.com/8treenet/freedom)

## 特点

✨ **现代化设计** - 基于 DaisyUI 5.3.7 的美观 UI
🚀 **高性能** - HTMX v2.0.7 部分页面更新
📱 **响应式设计** - 桌面优先，移动端友好
🎯 **最佳实践** - 遵循官方推荐用法
📝 **重构优化** - 代码精简，易于维护
🔧 **易于扩展** - 模块化设计

## 许可证

本项目仅用于学习目的。

---

**重构完成！Happy Learning! 🎉**
