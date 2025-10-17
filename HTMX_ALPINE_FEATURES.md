# HTMX & Alpine.js 特性展示指南

本文档介绍系统设置页面中展示的 HTMX 和 Alpine.js 高级特性。

## 🎯 HTMX 特性展示

### 1. **实时自动保存** (`hx-trigger`)
```html
<input 
    name="site_name" 
    hx-post="/settings/autosave" 
    hx-trigger="input changed delay:500ms"
    hx-target="#save-status"
    hx-indicator="#save-indicator">
```
**特性**：
- `input changed` - 只在值改变时触发
- `delay:500ms` - 防抖，500ms 后发送请求
- 自动保存，无需点击按钮

**应用场景**：编辑器、表单草稿、用户偏好设置

---

### 2. **禁用状态** (`hx-disable-elt`)
```html
<button 
    hx-post="/settings" 
    hx-disable-elt="this"
    class="button is-primary">
    保存设置
</button>
```
**特性**：
- 请求期间自动禁用按钮
- 防止重复提交
- `this` 指自身，也可以是 CSS 选择器

**应用场景**：防止表单重复提交、按钮防抖

---

### 3. **内联事件** (`hx-on`)
```html
<div hx-get="/settings/stats" 
     hx-on::after-request="console.log('请求完成')">
</div>
```
**特性**：
- 直接在 HTML 中处理 HTMX 事件
- `::` 前缀表示 HTMX 事件
- 简化事件处理代码

**应用场景**：简单的事件处理、调试、埋点

---

### 4. **请求同步** (`hx-sync`)
```html
<input 
    hx-post="/settings/search" 
    hx-sync="this:replace">
```
**特性**：
- `replace` - 新请求替换旧请求
- `drop` - 有请求进行时丢弃新请求
- `abort` - 中止旧请求
- `queue` - 排队等待

**应用场景**：搜索框、自动补全、防抖

---

### 5. **验证触发** (`hx-validate`)
```html
<input 
    type="email" 
    required 
    hx-post="/settings" 
    hx-validate="true">
```
**特性**：
- 提交前触发 HTML5 验证
- 验证失败阻止请求

**应用场景**：表单验证、数据校验

---

### 6. **轮询** (`hx-trigger="every"`)
```html
<div 
    hx-get="/settings/stats" 
    hx-trigger="load, every 5s"
    hx-swap="innerHTML">
    统计数据...
</div>
```
**特性**：
- 定时刷新内容
- `load` 初次加载
- `every 5s` 每5秒轮询

**应用场景**：实时数据、监控面板、聊天消息

---

### 7. **选择性更新** (`hx-select`)
```html
<div 
    hx-get="/settings" 
    hx-select="#settings-content"
    hx-target="#main">
```
**特性**：
- 只选择响应中的特定部分
- 减少不必要的 DOM 更新

**应用场景**：部分页面更新、内容过滤

---

### 8. **历史管理** (`hx-push-url`)
```html
<a 
    hx-get="/settings/theme" 
    hx-push-url="true">
    主题设置
</a>
```
**特性**：
- 自动更新浏览器 URL
- 支持前进/后退
- `true` - 使用请求 URL
- `"/custom"` - 自定义 URL

**应用场景**：SPA 导航、标签页切换

---

## 🎨 Alpine.js 特性展示

### 1. **过渡动画** (`x-transition`)
```html
<div x-show="isOpen" x-transition>
    内容淡入淡出
</div>

<div 
    x-show="isOpen" 
    x-transition:enter="transition ease-out duration-300"
    x-transition:enter-start="opacity-0 transform scale-90"
    x-transition:enter-end="opacity-100 transform scale-100">
    自定义动画
</div>
```
**特性**：
- 默认淡入淡出
- 自定义 enter/leave 动画
- 支持 CSS 过渡类

**应用场景**：模态框、下拉菜单、通知、侧边栏

---

### 2. **防闪烁** (`x-cloak`)
```html
<style>
[x-cloak] { display: none !important; }
</style>

<div x-cloak x-data="{ count: 0 }">
    {{ count }}
</div>
```
**特性**：
- Alpine 初始化前隐藏元素
- 防止未渲染的模板闪现

**应用场景**：所有 Alpine 组件的根元素

---

### 3. **数据监听** (`$watch`)
```html
<div x-data="{ 
    color: '#1976d2',
    init() {
        this.$watch('color', value => {
            console.log('颜色变更:', value);
            // 实时预览
            document.body.style.setProperty('--primary-color', value);
        });
    }
}">
```
**特性**：
- 监听数据变化
- 执行副作用
- 类似 Vue 的 watch

**应用场景**：实时预览、数据同步、联动更新

---

### 4. **DOM 更新后** (`$nextTick`)
```html
<button @click="
    show = true;
    $nextTick(() => {
        $refs.input.focus();
    });
">
    显示输入框
</button>
<input x-show="show" x-ref="input">
```
**特性**：
- 等待 DOM 更新完成
- 确保元素已渲染

**应用场景**：聚焦元素、测量尺寸、滚动定位

---

### 5. **全局状态** (`Alpine.store()`)
```javascript
// app.js
Alpine.store('settings', {
    unsavedChanges: false,
    lastSaveTime: null,
    
    markDirty() {
        this.unsavedChanges = true;
    },
    
    markSaved() {
        this.unsavedChanges = false;
        this.lastSaveTime = new Date();
    }
});
```

```html
<!-- 任何组件 -->
<div x-data>
    <span x-show="$store.settings.unsavedChanges">有未保存的更改</span>
    <span x-text="$store.settings.lastSaveTime"></span>
</div>
```
**特性**：
- 跨组件共享状态
- 响应式更新
- 类似 Vuex/Pinia

**应用场景**：全局状态、用户信息、主题设置

---

### 6. **副作用** (`x-effect`)
```html
<div x-data="{ 
    primary: '#1976d2',
    secondary: '#dc004e'
}" 
x-effect="
    document.documentElement.style.setProperty('--primary', primary);
    document.documentElement.style.setProperty('--secondary', secondary);
">
```
**特性**：
- 自动追踪依赖
- 数据变化时自动执行
- 类似 Vue 的 watchEffect

**应用场景**：主题切换、CSS 变量更新、localStorage 同步

---

### 7. **Magic Properties**

#### `$el` - 当前元素
```html
<button @click="$el.classList.add('clicked')">
    点击我
</button>
```

#### `$refs` - 元素引用
```html
<input x-ref="email" type="email">
<button @click="$refs.email.focus()">
    聚焦邮箱
</button>
```

#### `$dispatch` - 事件派发
```html
<div @settings-changed="console.log($event.detail)">
    <button @click="$dispatch('settings-changed', { key: 'theme' })">
        修改设置
    </button>
</div>
```

#### `$root` - 根元素
```html
<div x-data="{ count: 0 }">
    <button @click="$root.querySelector('.display').textContent = count++">
        增加
    </button>
    <span class="display">0</span>
</div>
```

---

### 8. **元素传送** (`x-teleport`)
```html
<div x-data="{ open: false }">
    <button @click="open = true">打开模态框</button>
    
    <template x-teleport="body">
        <div x-show="open" class="modal">
            <div class="modal-content">
                模态框内容
                <button @click="open = false">关闭</button>
            </div>
        </div>
    </template>
</div>
```
**特性**：
- 将元素渲染到其他位置
- 保持逻辑关联
- 类似 Vue 的 Teleport

**应用场景**：模态框、通知、工具提示

---

## 📦 实际应用示例

### 示例 1：带撤销的表单
```html
<div x-data="{
    current: { name: 'John', email: 'john@example.com' },
    history: [],
    
    change(field, value) {
        this.history.push({ ...this.current });
        this.current[field] = value;
        
        // 自动保存
        fetch('/settings/autosave', {
            method: 'POST',
            body: JSON.stringify(this.current)
        });
    },
    
    undo() {
        if (this.history.length > 0) {
            this.current = this.history.pop();
        }
    }
}">
    <input :value="current.name" @input="change('name', $event.target.value)">
    <button @click="undo()" :disabled="history.length === 0">
        撤销
    </button>
</div>
```

---

### 示例 2：实时主题预览
```html
<div x-data="{
    theme: {
        primary: '#1976d2',
        secondary: '#dc004e',
        background: '#ffffff'
    },
    
    applyTheme() {
        Object.entries(this.theme).forEach(([key, value]) => {
            document.documentElement.style.setProperty(`--${key}`, value);
        });
    }
}"
x-init="applyTheme()"
x-effect="applyTheme()">
    
    <input type="color" x-model="theme.primary">
    <input type="color" x-model="theme.secondary">
    
    <!-- 预览区域自动更新 -->
    <div class="preview-card" style="background: var(--primary)">
        预览
    </div>
</div>
```

---

### 示例 3：智能表单验证
```html
<form x-data="{
    form: { email: '', password: '' },
    errors: {},
    touched: {},
    
    validate(field) {
        this.touched[field] = true;
        
        if (field === 'email' && !this.form.email.includes('@')) {
            this.errors.email = '请输入有效的邮箱';
        } else {
            delete this.errors.email;
        }
    }
}"
@submit.prevent="console.log('提交', form)">
    
    <input 
        x-model="form.email"
        @blur="validate('email')"
        @input="touched.email && validate('email')"
        :class="{ 'is-danger': errors.email && touched.email }">
    
    <p x-show="errors.email && touched.email" 
       x-text="errors.email"
       x-transition
       class="help is-danger">
    </p>
</form>
```

---

## 🎓 学习建议

1. **HTMX** 专注于**服务端驱动**的交互
   - 减少 JavaScript 代码
   - 简化前后端交互
   - 适合传统服务端渲染

2. **Alpine.js** 处理**纯前端**的交互
   - 表单验证
   - UI 状态管理  
   - 动画和过渡

3. **结合使用**
   - HTMX 处理数据获取和提交
   - Alpine.js 处理 UI 逻辑
   - 互不干扰，各司其职

---

## 📚 参考资源

- HTMX 官方文档: https://htmx.org/
- Alpine.js 官方文档: https://alpinejs.dev/
- HTMX Examples: https://htmx.org/examples/
- Alpine.js Examples: https://alpinejs.dev/examples

