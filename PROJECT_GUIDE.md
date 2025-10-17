# HTMX 管理后台学习项目

这是一个基于 HTMX、Bulma 和 Alpine.js 构建的现代化 SPA 风格管理后台系统，用于学习 HTMX 的最佳实践。

## 技术栈

- **后端框架**: Go + Freedom Framework
- **前端库**: 
  - HTMX v2.0.7 - 处理部分页面更新和 AJAX 请求
  - Bulma v1.0.4 - CSS 框架
  - Alpine.js v3.15.0 - 纯前端交互
- **架构模式**: SPA（单页应用）

## 项目结构

```
gohtml/
├── adapter/
│   └── controller/          # 控制器层
│       ├── default.go       # 默认控制器（根路由）
│       ├── dashboard.go     # 仪表盘控制器
│       ├── user.go          # 用户管理控制器
│       ├── product.go       # 商品管理控制器
│       ├── order.go         # 订单管理控制器
│       └── setting.go       # 系统设置控制器
├── domain/
│   └── vo/                  # 值对象（View Objects）
│       ├── common.go        # 通用结构体（分页等）
│       ├── user.go          # 用户相关结构体
│       ├── product.go       # 商品相关结构体
│       ├── order.go         # 订单相关结构体
│       └── dashboard.go     # 仪表盘结构体
├── web/
│   ├── static/
│   │   └── css/
│   │       └── admin.css    # 自定义样式
│   ├── tmplfuncs/
│   │   └── endpoint.go      # 模板辅助函数
│   └── views/
│       ├── layout.html      # 主布局（SPA Shell）
│       ├── components/      # 通用组件
│       │   ├── sidebar.html # 侧边栏导航
│       │   └── header.html  # 顶部栏
│       ├── dashboard/       # 仪表盘视图
│       ├── users/           # 用户管理视图
│       ├── products/        # 商品管理视图
│       ├── orders/          # 订单管理视图
│       └── settings/        # 系统设置视图
├── main.go                  # 应用入口
└── PROJECT_GUIDE.md         # 项目说明文档
```

## 功能模块

### 1. 仪表盘 (Dashboard)
- **路由**: `/dashboard`
- **功能**:
  - 展示统计卡片（用户数、商品数、订单数等）
  - 自动刷新数据（每 30 秒）
  - 最近订单、最近用户、热门商品
- **HTMX 特性**:
  - `hx-trigger="every 30s"` - 定时自动刷新
  - `hx-swap-oob` - Out of Band 更新页面标题

### 2. 用户管理 (Users)
- **路由**: `/users`
- **功能**:
  - 用户列表（分页、搜索、筛选）
  - 新增用户（模态框表单）
  - 编辑用户（模态框表单 + 单行更新）
  - 删除用户（带确认对话框）
- **HTMX 特性**:
  - `hx-get` - 加载表单和列表
  - `hx-post` - 创建用户
  - `hx-put` - 更新用户（单行更新）
  - `hx-delete` - 删除用户
  - `hx-trigger="keyup changed delay:500ms"` - 实时搜索
  - `hx-include` - 包含其他表单字段
  - `hx-indicator` - 加载指示器
  - `hx-target` + `hx-swap` - 精确控制更新位置
- **Alpine.js 应用**:
  - 模态框打开/关闭
  - 表单验证
  - 删除确认对话框

### 3. 商品管理 (Products)
- **路由**: `/products`
- **功能**:
  - 商品网格展示
  - 分页、搜索、分类筛选
  - 完整 CRUD 操作
  - 单卡更新
- **HTMX 特性**:
  - 实时搜索和筛选
  - 卡片式布局的局部更新
  - 表单验证和提交
- **Alpine.js 应用**:
  - 模态框管理
  - 表单验证
  - 确认对话框

### 4. 订单管理 (Orders)
- **路由**: `/orders`
- **功能**:
  - 订单列表（分页、搜索、状态筛选）
  - 订单详情（模态框显示）
  - 状态更新（下拉菜单）
  - 取消订单
- **HTMX 特性**:
  - `hx-vals` - 添加额外参数
  - 下拉菜单状态更新
  - 订单详情动态加载
- **Alpine.js 应用**:
  - 下拉菜单控制
  - 状态更新
  - 确认对话框

### 5. 系统设置 (Settings)
- **路由**: `/settings`
- **功能**:
  - 分标签页配置（基本信息、联系方式、区域设置）
  - 表单提交和保存
  - 成功/失败反馈
- **HTMX 特性**:
  - 表单提交
  - Toast 通知
- **Alpine.js 应用**:
  - Tab 切换
  - 表单验证
  - 提交状态管理

## HTMX 特性展示

本项目展示了以下 HTMX 特性的使用：

1. **基本 HTTP 请求**: `hx-get`, `hx-post`, `hx-put`, `hx-delete`
2. **目标控制**: `hx-target` - 指定更新的元素
3. **内容替换**: `hx-swap` - innerHTML, outerHTML, beforeend 等
4. **触发器**: `hx-trigger` - 自定义触发事件（keyup, change, every 30s 等）
5. **加载指示**: `hx-indicator` - 显示加载状态
6. **URL 管理**: `hx-push-url` - 更新浏览器历史
7. **内容选择**: `hx-select` - 从响应中选择特定内容
8. **额外参数**: `hx-vals` - 添加请求参数
9. **包含字段**: `hx-include` - 包含其他输入字段
10. **Out of Band**: `hx-swap-oob` - 同时更新多个位置

## Alpine.js 应用场景

本项目展示了以下 Alpine.js 的使用：

1. **模态框管理** - 打开/关闭模态框
2. **下拉菜单** - 用户菜单、通知、筛选菜单
3. **表单验证** - 客户端实时验证
4. **Tab 切换** - 设置页面的标签页
5. **确认对话框** - 删除操作确认
6. **Toast 通知** - 全局通知管理
7. **状态管理** - 组件内部状态

## 启动项目

1. **安装依赖**:
```bash
go mod download
```

2. **运行项目**:
```bash
go run main.go
```

3. **访问应用**:
打开浏览器访问 `http://localhost:8000`

## 学习要点

### 1. SPA 架构
- 主布局只加载一次（`layout.html`）
- 内容区域通过 HTMX 动态更新
- 使用 `hx-push-url` 管理浏览器历史
- 侧边栏使用 `hx-get` + `hx-target="#main-container"` 实现导航

### 2. 部分页面更新
- HTMX 只请求和更新需要的部分
- 服务器返回 HTML 片段，不是完整页面
- 减少数据传输，提高性能

### 3. 渐进增强
- 表单可以正常提交（如果没有 HTMX）
- HTMX 增强用户体验，但不是必需的
- 优雅降级

### 4. 状态管理
- 服务器端保存状态（mock 数据在内存中）
- 客户端无状态（Alpine.js 只管理 UI 状态）
- 简化前端复杂度

### 5. 模块化设计
- 每个功能模块独立的控制器
- 视图文件分离（list、form、row/card）
- 可复用的组件（sidebar、header）

## 代码示例

### 实时搜索
```html
<input class="input" 
       type="search" 
       name="search"
       hx-get="/users"
       hx-trigger="keyup changed delay:500ms, search"
       hx-target="#user-table-container"
       hx-swap="outerHTML"
       hx-include="[name='status']">
```

### 单行更新
```html
<form hx-put="/users/{{.ID}}"
      hx-target="#user-row-{{.ID}}"
      hx-swap="outerHTML">
  <!-- 表单字段 -->
</form>
```

### 定时刷新
```html
<div hx-get="/dashboard/stats" 
     hx-trigger="every 30s" 
     hx-swap="innerHTML">
  <!-- 统计数据 -->
</div>
```

### Out of Band 更新
```html
<!-- 在响应中更新页面标题 -->
<div id="page-title" hx-swap-oob="true">用户管理</div>
```

### Alpine.js 模态框
```html
<div class="modal" x-ref="userModal">
  <button @click="$refs.userModal.classList.add('is-active')">
    打开模态框
  </button>
</div>
```

## 最佳实践

1. **语义化 HTML** - 使用正确的 HTTP 方法（GET, POST, PUT, DELETE）
2. **RESTful 路由** - 清晰的 URL 结构
3. **渐进增强** - 在没有 JavaScript 的情况下也能工作
4. **无障碍访问** - 使用适当的 ARIA 属性
5. **响应式设计** - Bulma 的栅格系统
6. **错误处理** - 友好的错误提示
7. **加载状态** - 使用 `hx-indicator` 显示加载状态
8. **代码注释** - 详细的中文注释

## 注意事项

1. **Mock 数据** - 所有数据都是内存中的 mock 数据，重启后会丢失
2. **无认证** - 这是一个学习项目，没有实现认证和授权
3. **无数据库** - 没有真实的数据持久化
4. **桌面优先** - 未针对移动端优化（企业应用）
5. **无深色主题** - 只有亮色主题

## 扩展建议

如果你想进一步学习，可以尝试：

1. 添加真实的数据库（MySQL、PostgreSQL）
2. 实现用户认证和权限管理
3. 添加文件上传功能
4. 实现 WebSocket 实时通知
5. 添加数据导出功能（Excel、CSV）
6. 实现更复杂的筛选和排序
7. 添加数据可视化图表
8. 实现批量操作功能
9. 添加操作日志
10. 实现多语言支持

## 相关资源

- [HTMX 官方文档](https://htmx.org/)
- [Bulma 官方文档](https://bulma.io/)
- [Alpine.js 官方文档](https://alpinejs.dev/)
- [Freedom 框架文档](https://github.com/8treenet/freedom)

## 许可证

本项目仅用于学习目的。

