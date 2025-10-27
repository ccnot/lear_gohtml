# AI 编程参考规则

## 项目概述

这是一个基于 HTMX v2.0.8、DaisyUI 5.3.7 和 Tailwind CSS 4 构建的 SPA（单页应用）管理后台系统。后端使用 Go + Freedom Framework v1.9.7。

## 核心技术栈

- **前端**: HTMX v2.0.8 + DaisyUI 5.3.7 + Tailwind CSS 4 + Alpine.js v3.15.0
- **后端**: Go + Freedom Framework v1.9.7
- **架构**: SPA（单页应用）
- **设计**: 现代、简洁、美观（桌面端优先）

## 核心原则

1. **HTMX 优先**: 所有数据交互和页面导航使用 HTMX 实现
2. **DaisyUI 为主**: UI 组件优先使用 DaisyUI，不足时使用 Tailwind CSS 补充
3. **Alpine.js 最小化**: 仅用于必要的 UI 状态管理，避免复杂逻辑
4. **SPA 架构**: 保持单页应用架构，避免传统页面跳转

## 1. HTMX 使用规则

### 1.1 核心原则

- **HTMX 优先**: 所有数据交互和页面导航必须使用 HTMX 实现
- **SPA 架构**: 主布局只加载一次，内容区域动态更新
- **目标容器**: 使用 `hx-target="main"` 更新主内容区域
- **历史管理**: 使用 `hx-push-url="true"` 管理浏览器历史记录
- **页面标题**: 使用 `hx-swap-oob="true"` 更新页面标题

### 1.2 请求模式

- **GET 请求**: 用于页面导航和数据获取
- **POST/PUT/DELETE**: 用于表单提交和数据操作
- **部分更新**: 使用 `hx-select` 选择特定内容片段
- **加载指示**: 使用 `hx-indicator` 显示加载状态

### 1.3 事件处理

- **实时搜索**: `hx-trigger="keyup changed delay:500ms"`
- **定时刷新**: `hx-trigger="every 30s"`
- **表单提交**: `hx-trigger="submit"`
- **确认对话框**: `hx-confirm="确定要删除吗？"`

### 1.4 内容交换策略

- **innerHTML**: 默认选项，替换目标元素内容
- **outerHTML**: 替换整个目标元素
- **beforebegin/afterbegin**: 在目标元素前后插入
- **beforeend/afterend**: 在目标元素内部前后插入

## 2. DaisyUI 使用规则

### 2.1 组件优先级

1. **优先使用 DaisyUI 组件**: 所有 UI 元素首先考虑使用 DaisyUI 提供的组件
2. **Tailwind CSS 补充**: 只有在 DaisyUI 组件不够用时，才使用 Tailwind CSS 自定义
3. **避免自定义 CSS**: 尽量不编写自定义 CSS，使用 DaisyUI 和 Tailwind CSS 类

### 2.2 参考指南

- **完整组件文档**: 参考 `doc/llms.txt` 获取所有 DaisyUI 5 组件的详细用法和示例
- **何时参考**: 需要使用特定组件时，先查看 llms.txt 中的组件说明和示例代码
- **组件选择**: 根据需求选择合适的组件，优先使用 DaisyUI 提供的解决方案

### 2.3 颜色系统

- **语义化颜色**: 使用 `primary`, `secondary`, `accent`, `success`, `warning`, `error`
- **中性色**: 使用 `neutral`, `base-100`, `base-200`, `base-300`
- **内容色**: 使用 `*-content` 变量确保对比度

### 2.4 响应式设计

- **断点前缀**: `sm:`, `md:`, `lg:`, `xl:`
- **桌面优先**: 默认样式针对桌面端，移动端使用断点前缀覆盖

## 3. Alpine.js 使用限制

### 3.1 使用原则

- **最小化使用**: 只有在 DaisyUI 和 HTMX 无法满足需求时才使用 Alpine.js
- **状态管理**: 仅用于简单的 UI 状态管理（如侧边栏开关、主题切换）
- **避免复杂逻辑**: 不使用 Alpine.js 处理复杂业务逻辑

### 3.2 允许的使用场景

- **侧边栏状态管理**: 控制侧边栏的开关状态
- **主题切换**: 控制明暗主题切换
- **简单的显示/隐藏**: 纯前端UI元素的显示隐藏

### 3.3 禁止的使用场景

- **数据获取**: 使用 HTMX 而非 Alpine.js
- **表单处理**: 使用 HTMX 而非 Alpine.js
- **复杂状态管理**: 使用后端状态而非前端状态
- **路由管理**: 使用 HTMX 的 `hx-push-url` 而非 Alpine.js

## 4. 代码组织规则

### 4.1 文件结构

```
web/
├── views/
│   ├── layout.html          # 主布局
│   ├── components/         # 通用组件
│   │   ├── header.html
│   │   ├── sidebar.html
│   │   └── pagination.html
│   ├── dashboard/          # 仪表盘页面
│   ├── users/             # 用户管理页面
│   ├── products/          # 商品管理页面
│   ├── orders/            # 订单管理页面
│   └── settings/          # 系统设置页面
├── static/
│   ├── css/admin.css      # 最小化自定义样式
│   └── js/app.js          # 最小化自定义脚本
└── tmplfuncs/            # 模板辅助函数
```

### 4.2 模板复用

- **组件化**: 将可复用的部分提取为独立组件
- **模板继承**: 使用 `{{template}}` 复用模板片段
- **参数传递**: 使用 `dict` 传递复杂参数

```html
<!-- 分页组件使用示例 -->
{{$ctx := dict "BaseURL" "/users" "PageInfo" .PageInfo "TargetContainer" "user-table-container" "ExtraParams" (dict "keyword" .Query "status" .Status)}}
{{template "components/pagination.html" $ctx}}
```

### 4.3 样式组织

- **最小化自定义**: 只在必要时添加自定义样式

### 4.4 Freedom 框架模板使用规则

- **路径直接引入**: 模板文件无需手动 {{define}} 即可通过路径直接引入（如 {{template "components/sidebar.html"}}）
- **选择性内容使用**: 若需从一个模板文件中选择性使用部分内容（如仅引入 10 个 div 中的一个），应将各部分分别用 {{define "name"}}...{{end}} 命名定义，再通过模板名按需调用（如 {{template "name" .}}）

## 5. 最佳实践

### 5.1 性能优化
- **按需加载**: 使用 HTMX 按需加载页面内容

### 5.2 用户体验

- **加载状态**: 使用 `hx-indicator` 显示加载状态
- **错误处理**: 统一的错误提示机制
- **确认对话框**: 危险操作使用确认对话框
- **Toast 通知**: 操作反馈使用 Toast 通知

### 5.3 可维护性

- **命名规范**: 使用语义化的 ID 和类名
- **代码注释**: 为复杂逻辑添加注释
- **组件化**: 将可复用部分提取为组件
- **一致性**: 保持代码风格一致

## 6. 开发指导

### 6.1 组件选择流程

1. **数据交互**: 优先使用 HTMX 实现所有数据请求和页面更新
2. **UI 组件**: 优先使用 DaisyUI 组件，不足时使用 Tailwind CSS 补充
3. **状态管理**: 仅在必要时使用 Alpine.js 处理简单 UI 状态

### 6.2 文档参考策略

- **DaisyUI 组件**: 参考 `doc/llms.txt` 获取完整组件列表和用法，建议先 grep，文档较大。
- **HTMX 功能**: 参考 HTMX 官方文档了解高级特性
- **项目结构**: 参考现有代码文件保持一致性

## 7. 注意事项

1. **避免过度使用 Alpine.js**: 优先使用 HTMX 和 DaisyUI
2. **保持 SPA 架构**: 不要使用传统的页面跳转
3. **统一错误处理**: 使用 Toast 通知统一处理错误
4. **响应式设计**: 确保移动端体验
5. **性能考虑**: 避免不必要的请求和 DOM 操作
6. **可访问性**: 添加适当的 ARIA 标签和语义化 HTML

## 8. 参考资源

- [HTMX v2.0.8 文档](https://htmx.org/)
- [DaisyUI 5.3.7 文档](https://daisyui.com/)
- [Tailwind CSS 4 文档](https://tailwindcss.com/)
- [Alpine.js v3.15.0 文档](https://alpinejs.dev/)
- [Freedom Framework 文档](https://github.com/8treenet/freedom)

---

**重要提醒**: 本规则文档旨在为 AI 编程提供参考，确保代码风格一致性和最佳实践。在实际开发中，请根据具体需求灵活应用这些规则。

**文档目的**: 本文档不是项目总结，而是为 AI 编程提供的规则指南，注重实用性和可操作性，确保生成的代码符合项目架构和最佳实践。