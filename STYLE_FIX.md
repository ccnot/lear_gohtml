# 样式修复说明 - Bulma 覆盖

## 问题描述

用户反馈：
1. 表格表头是灰底灰字，看不清
2. 输入框的黑色边框很奇怪
3. 字体颜色不够清晰

**根本原因**：Bulma CSS 框架的默认样式覆盖了我们的自定义样式

## 修复策略

使用 `!important` 强制覆盖 Bulma 的默认样式，确保企业级应用的可读性和专业性。

## 修复内容

### 1. 表格样式强化

#### 修复前的问题
- Bulma 的表格默认有灰色文字
- 表头颜色对比度不够
- 边框颜色不清晰

#### 修复后
```css
/* 表格 */
.data-table th {
    font-weight: 700 !important;        /* 加粗 */
    color: #212121 !important;          /* 深黑色文字 */
    background: #f5f5f5 !important;     /* 浅灰背景 */
    border-color: #dbdbdb !important;   /* 清晰边框 */
}

.data-table td {
    color: #212121 !important;          /* 深黑色文字 */
    border-color: #e0e0e0 !important;   /* 统一边框 */
}
```

**效果**：
- ✅ 表头：浅灰背景 + 深黑色加粗文字 = 高对比度
- ✅ 表格内容：深黑色文字，清晰易读
- ✅ 边框：统一的浅灰色，不抢眼但清晰

### 2. 输入框优化

#### 修复前的问题
- Bulma 的输入框有深色边框和阴影
- 焦点状态的阴影太重
- 文字颜色不够清晰

#### 修复后
```css
.input, .textarea, .select select {
    border: 1px solid #dbdbdb !important;     /* 浅灰边框 */
    color: #363636 !important;                /* 深灰文字 */
    background-color: white !important;        /* 白色背景 */
    box-shadow: none !important;              /* 无阴影 */
}

.input:focus, .textarea:focus, .select select:focus {
    border-color: #1976d2 !important;         /* 蓝色边框 */
    box-shadow: 0 0 0 0.125em rgba(25, 118, 210, 0.25) !important;  /* 轻微光晕 */
}
```

**效果**：
- ✅ 普通状态：浅灰色边框，清爽简洁
- ✅ 悬停状态：稍深的灰色边框
- ✅ 焦点状态：蓝色边框 + 淡蓝色光晕，清晰的视觉反馈
- ✅ 文字：深灰色，清晰易读

### 3. 按钮样式统一

```css
.button.is-primary {
    background-color: #1976d2 !important;  /* 专业蓝 */
    border-color: #1976d2 !important;
    color: white !important;
}

.button.is-primary:hover {
    background-color: #1565c0 !important;  /* 深一点的蓝 */
}
```

**效果**：
- ✅ 主按钮：清晰的蓝色，符合企业级应用风格
- ✅ 悬停效果：颜色稍深，明确的交互反馈

### 4. 文字颜色层次

```css
body {
    color: #363636 !important;           /* 全局文字：深灰 */
}

strong, b {
    color: #212121 !important;           /* 强调文字：更深 */
    font-weight: 600 !important;
}

small {
    color: #757575 !important;           /* 次要文字：中灰 */
}
```

**文字层次**：
- **#212121** - 重要文字（标题、表头、强调）
- **#363636** - 普通文字（正文）
- **#757575** - 次要文字（提示、说明）
- **#b5b5b5** - 占位符文字

### 5. 分页组件优化

```css
.pagination-link, .pagination-previous, .pagination-next {
    color: #363636 !important;           /* 深色文字 */
    border-color: #dbdbdb !important;    /* 浅灰边框 */
}

.pagination-link.is-current {
    background-color: #1976d2 !important; /* 当前页蓝色 */
    color: white !important;
}
```

**效果**：
- ✅ 分页按钮：清晰的边框和文字
- ✅ 当前页：蓝色背景，白色文字，一目了然

### 6. 卡片和容器

```css
.content-card {
    background: white !important;
}

.card-header {
    background: white !important;
    border-bottom: 1px solid #e0e0e0 !important;
}

.card-title {
    color: #212121 !important;
}
```

**效果**：
- ✅ 白色背景，干净整洁
- ✅ 清晰的边框分隔
- ✅ 深色标题，突出重点

## 颜色系统

### 主题色
- **主色**：`#1976d2` - 专业蓝
- **成功**：`#4caf50` - Material 绿
- **警告**：`#ff9800` - Material 橙
- **危险**：`#f44336` - Material 红

### 中性色
- **深黑**：`#212121` - 标题、重要文字
- **深灰**：`#363636` - 正文文字
- **中灰**：`#757575` - 次要文字
- **浅灰**：`#b5b5b5` - 占位符、分隔
- **边框**：`#dbdbdb` / `#e0e0e0` - 边框颜色

### 背景色
- **白色**：`#ffffff` - 卡片、表格
- **浅灰**：`#f5f5f5` - 表头、页面背景
- **极浅灰**：`#fafafa` / `#f9f9f9` - 悬停效果

## 使用 !important 的原因

Bulma 是一个强大的 CSS 框架，但它的样式优先级很高。为了确保我们的企业级定制样式生效，必须使用 `!important`。

**使用场景**：
1. ✅ 覆盖第三方框架的默认样式
2. ✅ 确保关键视觉元素的一致性
3. ✅ 企业级应用的品牌统一

**注意事项**：
- 只在必要时使用 `!important`
- 保持样式的可维护性
- 文档化所有覆盖的原因

## 测试清单

启动项目后，检查以下内容：

### 表格
- [ ] 表头背景：浅灰色
- [ ] 表头文字：深黑色，加粗，清晰可读
- [ ] 表格内容文字：深黑色，易读
- [ ] 边框：清晰但不抢眼
- [ ] 悬停效果：轻微的灰色背景

### 输入框
- [ ] 边框：浅灰色，1px
- [ ] 文字：深灰色，清晰
- [ ] 背景：纯白
- [ ] 焦点状态：蓝色边框 + 淡蓝光晕
- [ ] 占位符：中灰色，不太突出

### 按钮
- [ ] 主按钮：蓝色背景，白色文字
- [ ] 悬停：颜色稍深
- [ ] 文字清晰易读

### 整体视觉
- [ ] 文字层次分明
- [ ] 颜色对比度足够
- [ ] 符合企业级应用审美
- [ ] 没有过度的阴影或渐变

## 对比效果

### 修复前
```
表头：灰底灰字 ❌ 看不清
输入框：黑边框 + 重阴影 ❌ 太重
文字：浅色 ❌ 不够清晰
```

### 修复后
```
表头：浅灰底 + 深黑字 ✅ 高对比度
输入框：浅灰边框 + 无阴影 ✅ 清爽
文字：深色层次分明 ✅ 清晰易读
```

## 浏览器兼容性

所有样式使用标准 CSS 属性，兼容：
- ✅ Chrome/Edge (Chromium)
- ✅ Firefox
- ✅ Safari
- ✅ Opera

## 维护建议

1. **保持一致性**：新增组件也要遵循相同的颜色系统
2. **避免滥用 !important**：只在覆盖 Bulma 时使用
3. **定期检查**：确保 Bulma 更新后样式仍然正常
4. **文档化**：记录所有重要的样式覆盖

---

**修复完成时间**：2025-10-17  
**修复状态**：✅ 完成  
**测试状态**：✅ 编译通过  

