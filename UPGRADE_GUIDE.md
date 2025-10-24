# Tailwind CSS 4 + daisyUI 5 升级指南

## 概述

本项目已成功从 Tailwind CSS 3.3.5 + daisyUI 4 升级到 Tailwind CSS 4 + daisyUI 5。

## 主要变更

### 1. CSS 框架升级
- **Tailwind CSS**: 3.3.5 → 4.x
- **daisyUI**: 4.x → 5.x

### 2. 新组件和特性

#### Theme Controller（主题控制器）
- 使用 daisyUI 5 的新 **Theme Controller** 组件
- 基于 `<label class="swap">` 和 `<input class="theme-controller">` 实现
- 自动保存主题到 localStorage
- 位置：`web/views/components/header.html`

```html
<label class="swap swap-rotate btn btn-ghost btn-circle">
    <input type="checkbox" class="theme-controller" value="dark" />
    <!-- 太阳/月亮图标 -->
</label>
```

#### 原生 Dialog 模态框
- 使用 HTML5 原生 `<dialog>` 元素
- 不再需要 `<input type="checkbox">` 的方式
- 使用 `dialog.showModal()` 和 `dialog.close()` 方法
- 位置：
  - `web/views/components/confirm-dialog.html`
  - `web/views/layout.html`（订单详情模态框）

```html
<dialog id="my-modal" class="modal">
    <div class="modal-box">
        <!-- 内容 -->
    </div>
    <form method="dialog" class="modal-backdrop">
        <button>close</button>
    </form>
</dialog>
```

#### 改进的 Dropdown
- 使用 daisyUI 5 优化的 dropdown 组件
- 更好的 z-index 管理
- 位置：
  - `web/views/components/header.html`（通知和用户菜单）
  - `web/views/orders/row.html`（订单状态更新）

#### 统一的分页组件
- 使用 daisyUI 5 的 `join` 组件
- 支持 HTMX 动态加载
- 位置：`web/views/components/pagination.html`

```html
<div class="join">
    <button class="join-item btn btn-sm">«</button>
    <button class="join-item btn btn-sm btn-active">1</button>
    <button class="join-item btn btn-sm">»</button>
</div>
```

#### Stats 统计卡片
- 完全使用 daisyUI 5 的 stats 组件
- 移除了自定义 CSS 样式
- 支持悬停效果
- 位置：`web/views/dashboard/stats.html`

```html
<div class="stats shadow-lg border border-base-300">
    <div class="stat bg-primary text-primary-content">
        <div class="stat-figure"><!-- 图标 --></div>
        <div class="stat-title">标题</div>
        <div class="stat-value">数值</div>
        <div class="stat-desc">描述</div>
    </div>
</div>
```

#### Collapse 折叠菜单
- 侧边栏使用 HTML5 `<details>` 元素
- 原生支持展开/折叠动画
- 位置：`web/views/components/sidebar.html`

```html
<details open>
    <summary>用户管理</summary>
    <ul>
        <li><a>用户列表</a></li>
        <li><a>角色管理</a></li>
    </ul>
</details>
```

### 3. 样式优化

#### admin.css 简化
- 移除了所有 `@apply` 指令（Tailwind CSS 4 不再推荐）
- 使用标准 CSS
- 保留必要的自定义样式和动画
- 位置：`web/static/css/admin.css`

#### SVG 图标替代 Font Awesome
- 主要组件使用 Heroicons（SVG）替代部分 Font Awesome 图标
- 更好的性能和可定制性
- Font Awesome 仍然保留用于向后兼容

### 4. 响应式和可访问性改进

- 使用原生 HTML 元素（`<dialog>`, `<details>`）提升可访问性
- 改进移动端体验
- 优化键盘导航

## 文件变更清单

### 核心文件
- ✅ `web/static/css/admin.css` - 完全重写
- ✅ `web/views/layout.html` - 简化配置和脚本
- ✅ `web/views/components/header.html` - 新 Theme Controller
- ✅ `web/views/components/sidebar.html` - 使用 collapse
- ✅ `web/views/components/confirm-dialog.html` - 原生 dialog
- ✅ `web/views/components/pagination.html` - 新分页组件

### 视图文件
- ✅ `web/views/dashboard/index.html` - 优化布局
- ✅ `web/views/dashboard/stats.html` - 新 stats 组件
- ✅ `web/views/users/row.html` - 简化按钮
- ✅ `web/views/orders/row.html` - 新 dropdown 和 dialog
- ✅ `web/views/products/card.html` - 优化卡片样式

## 兼容性说明

### 浏览器支持
- Chrome/Edge 90+
- Firefox 88+
- Safari 14+
- 移动浏览器：iOS Safari 14+, Chrome Android 90+

### 后端接口
- 无需修改后端接口
- HTMX 行为保持不变
- 所有 API 端点保持兼容

## 测试清单

### 功能测试
- [x] 主题切换（浅色/深色模式）
- [x] 侧边栏展开/折叠
- [x] 用户管理二级菜单
- [x] 订单详情模态框
- [x] 确认删除对话框
- [x] 分页导航
- [x] 下拉菜单（通知、用户菜单、订单状态）
- [x] 响应式布局（移动端/桌面端）

### 样式测试
- [x] 统计卡片显示
- [x] 表格样式
- [x] 按钮状态（hover, active, disabled）
- [x] 徽章和标签
- [x] Toast 通知
- [x] 加载动画

## 性能提升

- **减少 CSS 体积**: 移除不必要的自定义样式
- **使用原生元素**: 减少 JavaScript 开销
- **优化加载**: CDN 加速加载 Tailwind 和 daisyUI
- **更好的缓存**: 静态资源更容易缓存

## 开发建议

### 使用新组件
```html
<!-- ✅ 推荐：使用 daisyUI 5 组件 -->
<button class="btn btn-primary">按钮</button>
<div class="badge badge-success">徽章</div>

<!-- ❌ 避免：自定义样式 -->
<button class="custom-button">按钮</button>
```

### 主题支持
```css
/* ✅ 推荐：使用 daisyUI 主题变量 */
.custom-element {
    background: oklch(var(--p));
    color: oklch(var(--pc));
}

/* ❌ 避免：硬编码颜色 */
.custom-element {
    background: #1890ff;
    color: white;
}
```

## 故障排除

### 主题不切换
- 检查 `localStorage` 中的 `theme` 值
- 确保 `data-theme` 属性正确设置在 `<html>` 元素上
- 检查浏览器控制台是否有 JavaScript 错误

### 模态框不显示
- 确认使用 `dialog.showModal()` 而不是旧的 checkbox 方法
- 检查 `z-index` 是否被其他元素覆盖

### 样式错乱
- 清除浏览器缓存
- 检查 CDN 是否正确加载
- 查看 admin.css 是否正确引入

## 未来改进

- [ ] 考虑使用 Tailwind CSS 构建工具代替 CDN
- [ ] 添加更多自定义主题
- [ ] 优化移动端侧边栏动画
- [ ] 添加暗色模式下的颜色微调

## 参考资料

- [daisyUI 5 文档](https://daisyui.com/)
- [Tailwind CSS 4 升级指南](https://tailwindcss.com/docs/upgrade-guide)
- [daisyUI 5 发布说明](https://github.com/saadeghi/daisyui/releases)

---

**升级完成日期**: 2025-01-18  
**测试状态**: ✅ 通过

