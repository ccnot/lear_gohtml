/**
 * Admin Dashboard - 全局 JavaScript
 * 重构版本：基于 HTMX v2.0.7 + DaisyUI 5.3.7 + Alpine.js v3.15.0
 */

// ============================================
// Alpine.js 组件注册
// ============================================

document.addEventListener('alpine:init', () => {

    // 删除确认组件 - 简化版本
    Alpine.data('deleteConfirm', (options) => ({
        confirmDelete() {
            showConfirmDialog({
                title: options.title || '确认删除',
                message: options.message || '确定要删除这个项目吗？',
                confirmText: options.confirmText || '删除',
                onConfirm: () => {
                    // 使用 HTMX 发起请求
                    const element = document.querySelector(options.target);
                    if (element) {
                        htmx.trigger(element, 'delete');
                    }
                }
            });
        }
    }));

    // 通用表单验证器 - 简化版本
    function createValidator(rules) {
        return function (formData) {
            const errors = {};
            for (const [field, rule] of Object.entries(rules)) {
                const value = formData[field];
                if (rule.required && !value) {
                    errors[field] = rule.message || `${field}不能为空`;
                }
                if (rule.pattern && value && !rule.pattern.test(value)) {
                    errors[field] = rule.patternMessage || `${field}格式不正确`;
                }
                if (rule.validator && value) {
                    const customError = rule.validator(value);
                    if (customError) {
                        errors[field] = customError;
                    }
                }
            }
            return errors;
        };
    }

    // 基础验证规则
    const validationRules = {
        username: {
            required: true,
            pattern: /^[a-zA-Z0-9_]{3,20}$/,
            message: '用户名不能为空',
            patternMessage: '用户名只能包含字母、数字和下划线，长度3-20位'
        },
        email: {
            required: true,
            pattern: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
            message: '邮箱不能为空',
            patternMessage: '邮箱格式不正确'
        },
        realName: {
            required: true,
            message: '真实姓名不能为空'
        },
        name: {
            required: true,
            message: '商品名称不能为空'
        },
        price: {
            validator: (value) => {
                const price = parseFloat(value);
                if (isNaN(price) || price < 0) {
                    return '请输入有效的价格';
                }
            }
        }
    };

    // 通用表单组件
    function createFormComponent(defaultForm, rules) {
        return {
            loading: false,
            errors: {},
            form: { ...defaultForm },

            validate() {
                const validator = createValidator(rules);
                this.errors = validator(this.form);
                return Object.keys(this.errors).length === 0;
            },

            submitForm(event) {
                if (!this.validate()) {
                    event.preventDefault();
                    return false;
                }
                this.loading = true;
            },

            clearErrors() {
                this.errors = {};
            }
        };
    }

    // 用户表单组件
    Alpine.data('userForm', (isEdit = false) => {
        const rules = { ...validationRules };
        if (isEdit) delete rules.username;
        return createFormComponent({
            username: '', email: '', realName: '', phone: '', role: 'viewer', status: 'active'
        }, rules);
    });

    // 商品表单组件
    Alpine.data('productForm', () => createFormComponent({
        name: '', sku: '', category: '', price: '', stock: '', status: 'active', description: ''
    }, {
        name: validationRules.name,
        price: validationRules.price
    }));
});

// ============================================
// HTMX 事件监听 - 简化版本
// ============================================

// HTMX confirm 拦截 - 使用 DaisyUI dialog
document.addEventListener('htmx:confirm', function (event) {
    if (!event.detail.question) return;

    event.preventDefault();
    showConfirmDialog({
        title: '确认操作',
        message: event.detail.question,
        confirmText: '确定',
        onConfirm: () => event.detail.issueRequest(true)
    });
});

// Toast 消息处理
document.addEventListener('DOMContentLoaded', function () {
    // 监听 HTMX 请求完成后的 Toast 消息
    document.body.addEventListener('htmx:afterSwap', function (event) {
        const xhr = event.detail.xhr;
        if (!xhr) return;

        const toastMessage = xhr.getResponseHeader('X-Toast-Message');
        const toastType = xhr.getResponseHeader('X-Toast-Type') || 'success';

        console.log('Toast Debug - Message:', toastMessage, 'Type:', toastType); // 调试日志

        if (toastMessage) {
            let decodedMessage = toastMessage;
            try {
                decodedMessage = decodeURIComponent(toastMessage);
            } catch (e) {
                // 解码失败时使用原始消息
            }
            showToast(decodedMessage, toastType);
        }
    });

    // 监听 HTMX 错误
    document.body.addEventListener('htmx:responseError', function (event) {
        showToast('请求失败，请稍后重试', 'error');
    });
});

// Toast 显示函数
let lastToast = null;
let lastToastTime = 0;

function showToast(message, type = 'info') {
    const container = document.getElementById('toast-container');
    if (!container) return;

    const now = Date.now();
    if (lastToast === message && (now - lastToastTime) < 1000) {
        return; // 防止重复消息
    }

    lastToast = message;
    lastToastTime = now;

    const alert = document.createElement('div');
    alert.className = `alert alert-${type} shadow-lg flex items-center justify-between`;

    const icons = {
        success: '✅',
        error: '❌',
        warning: '⚠️',
        info: 'ℹ️'
    };

    alert.innerHTML = `
        <div class="flex items-center gap-2">
            <span class="text-lg">${icons[type] || icons.info}</span>
            <span>${message}</span>
        </div>
        <button class="btn btn-ghost btn-xs" onclick="this.parentElement.remove()">✕</button>
    `;

    container.appendChild(alert);

    setTimeout(() => {
        if (alert.parentElement) {
            alert.remove();
        }
    }, 3000);
}

// 确认对话框组件
let isDialogShowing = false;

function showConfirmDialog(config) {
    if (isDialogShowing) return;

    createConfirmDialog();
    const dialog = document.getElementById('confirm-dialog');
    if (!dialog) return;

    isDialogShowing = true;

    document.getElementById('confirm-title').textContent = config.title || '确认操作';
    document.getElementById('confirm-message').textContent = config.message || '确定要执行此操作吗？';
    document.getElementById('confirm-text').textContent = config.confirmText || '确定';

    const confirmBtn = document.getElementById('confirm-button');
    confirmBtn.onclick = function () {
        if (typeof config?.onConfirm === 'function') {
            config.onConfirm();
        }
        closeConfirmDialog();
    };

    dialog.showModal();
}

function createConfirmDialog() {
    if (document.getElementById('confirm-dialog')) return;

    const dialogHTML = `
        <dialog id="confirm-dialog" class="modal">
            <div class="modal-box max-w-md">
                <div class="flex items-start gap-4 mb-4">
                    <div class="avatar">
                        <div class="w-12 rounded-full bg-warning/20 text-warning flex items-center justify-center">
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                            </svg>
                        </div>
                    </div>
                    <div class="flex-1">
                        <h3 class="font-bold text-lg" id="confirm-title"></h3>
                        <p class="py-2 text-sm opacity-70" id="confirm-message"></p>
                    </div>
                </div>
                <div class="modal-action">
                    <button class="btn btn-ghost" onclick="closeConfirmDialog()">取消</button>
                    <button class="btn btn-warning" id="confirm-button">
                        <span id="confirm-text"></span>
                    </button>
                </div>
            </div>
            <form method="dialog" class="modal-backdrop">
                <button onclick="closeConfirmDialog()">close</button>
            </form>
        </dialog>
    `;

    document.body.insertAdjacentHTML('beforeend', dialogHTML);
}

function closeConfirmDialog() {
    const dialog = document.getElementById('confirm-dialog');
    if (dialog) dialog.close();
    isDialogShowing = false;
}


// ============================================
// 工具函数 - 简化版本
// ============================================

function formatDate(date, format = 'YYYY-MM-DD HH:mm:ss') {
    const d = new Date(date);
    const year = d.getFullYear();
    const month = String(d.getMonth() + 1).padStart(2, '0');
    const day = String(d.getDate()).padStart(2, '0');
    const hours = String(d.getHours()).padStart(2, '0');
    const minutes = String(d.getMinutes()).padStart(2, '0');

    return format
        .replace('YYYY', year)
        .replace('MM', month)
        .replace('DD', day)
        .replace('HH', hours)
        .replace('mm', minutes);
}

function formatCurrency(amount, currency = '¥') {
    return `${currency}${parseFloat(amount).toFixed(2)}`;
}

function debounce(func, wait = 300) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

// 工具函数挂载
window.utils = { formatDate, formatCurrency, debounce };

// 通用模态框管理函数
window.modal = {
    show: function (modalId) {
        const modal = document.getElementById(modalId);
        if (modal) {
            modal.showModal();
            return true;
        }
        console.warn(`模态框 #${modalId} 未找到`);
        return false;
    },

    close: function (modalId) {
        const modal = document.getElementById(modalId);
        if (modal) {
            modal.close();
            return true;
        }
        console.warn(`模态框 #${modalId} 未找到`);
        return false;
    },

    toggle: function (modalId) {
        const modal = document.getElementById(modalId);
        if (modal) {
            if (modal.open) {
                modal.close();
            } else {
                modal.showModal();
            }
            return true;
        }
        console.warn(`模态框 #${modalId} 未找到`);
        return false;
    }
};


// ============================================
// 从HTML文件迁移的脚本函数
// ============================================

// 从 web/views/components/header.html 迁移的主题控制器初始化
(function () {
    const themeController = document.querySelector('.theme-controller');
    if (!themeController) return;

    // 从 localStorage 读取主题
    const savedTheme = localStorage.getItem('theme') || 'light';
    const isDark = savedTheme === 'dark';

    // 设置初始状态
    document.documentElement.setAttribute('data-theme', savedTheme);
    themeController.checked = isDark;

    // 监听主题切换
    themeController.addEventListener('change', function (e) {
        const newTheme = e.target.checked ? 'dark' : 'light';
        document.documentElement.setAttribute('data-theme', newTheme);
        localStorage.setItem('theme', newTheme);
    });
})();

// 从 web/views/products/new.html 迁移的商品表单函数
function productForm(isEdit = false) {
    return {
        loading: false,
        errors: {},
        form: {
            name: '',
            sku: '',
            category: '',
            price: '',
            stock: '',
            status: 'active',
            description: ''
        },

        validateForm() {
            this.errors = {};
            let isValid = true;

            if (!this.form.name) {
                this.errors.name = '商品名称不能为空';
                isValid = false;
            }

            if (!this.form.sku && !isEdit) {
                this.errors.sku = 'SKU不能为空';
                isValid = false;
            }

            if (!this.form.category) {
                this.errors.category = '请选择商品分类';
                isValid = false;
            }

            if (!this.form.price || parseFloat(this.form.price) < 0) {
                this.errors.price = '请输入有效的价格';
                isValid = false;
            }

            if (!this.form.stock || parseInt(this.form.stock) < 0) {
                this.errors.stock = '请输入有效的库存数量';
                isValid = false;
            }

            if (!this.form.status) {
                this.errors.status = '请选择商品状态';
                isValid = false;
            }

            return isValid;
        },

        submitForm(event) {
            if (this.validateForm()) {
                this.loading = true;

                // 让 HTMX 处理表单提交
                htmx.trigger(event.target, 'submit');

                // 监听提交完成
                setTimeout(() => {
                    this.loading = false;
                }, 500);
            }
        }
    }
}

// 从 web/views/products/form.html 迁移的产品表单函数
function productFormForForm(isEdit) {
    return {
        loading: false,
        errors: {},
        form: {
            name: '',
            sku: '',
            category: '',
            price: '',
            stock: '',
            status: 'active',
            description: ''
        },

        validateForm() {
            this.errors = {};
            let isValid = true;

            if (!this.form.name) {
                this.errors.name = '商品名称不能为空';
                isValid = false;
            }

            if (!this.form.sku && !isEdit) {
                this.errors.sku = 'SKU 不能为空';
                isValid = false;
            }

            if (!this.form.price || parseFloat(this.form.price) <= 0) {
                this.errors.price = '请输入有效的价格';
                isValid = false;
            }

            if (!this.form.stock || parseInt(this.form.stock) < 0) {
                this.errors.stock = '请输入有效的库存';
                isValid = false;
            }

            return isValid;
        },

        submitForm(event) {
            if (this.validateForm()) {
                this.loading = true;
                // 让 HTMX 处理表单提交
                htmx.trigger(event.target, 'submit');

                // 监听提交完成
                setTimeout(() => {
                    this.loading = false;
                    // 关闭模态框
                    document.querySelector('.modal').classList.remove('is-active');
                }, 500);
            }
        }
    };
}

// 从 web/views/products/edit.html 迁移的编辑商品表单函数
function editProductForm(isEdit = false, productId = null) {
    return {
        loading: false,
        errors: {},
        form: {
            name: '',
            sku: '',
            category: '',
            price: '',
            stock: '',
            status: '',
            description: ''
        },

        handleFormInit(event, productData) {
            // 编辑模式：预填充表单数据
            if (isEdit && productData) {
                this.form = {
                    name: productData.name || '',
                    sku: productData.sku || '',
                    category: productData.category || '',
                    price: productData.price || '',
                    stock: productData.stock || '',
                    status: productData.status || '',
                    description: productData.description || ''
                };
            }
        },

        validateForm() {
            this.errors = {};
            let isValid = true;

            if (!this.form.name) {
                this.errors.name = '商品名称不能为空';
                isValid = false;
            }

            if (!this.form.category) {
                this.errors.category = '请选择商品分类';
                isValid = false;
            }

            if (!this.form.price || parseFloat(this.form.price) < 0) {
                this.errors.price = '请输入有效的价格';
                isValid = false;
            }

            if (!this.form.stock || parseInt(this.form.stock) < 0) {
                this.errors.stock = '请输入有效的库存数量';
                isValid = false;
            }

            if (!this.form.status) {
                this.errors.status = '请选择商品状态';
                isValid = false;
            }

            return isValid;
        },

        submitForm(event) {
            if (!this.validateForm()) {
                return;
            }

            this.loading = true;

            // 让 HTMX 处理表单提交
            htmx.trigger(event.target, 'submit');

            // 监听提交完成
            setTimeout(() => {
                this.loading = false;
            }, 500);
        }
    }
}

// 从 web/views/settings/form.html 迁移的设置表单函数
function settingsForm() {
    return {
        loading: false,

        submitForm(event) {
            this.loading = true;
            // 让 HTMX 处理表单提交
            htmx.trigger(event.target, 'submit');

            // 监听提交完成
            setTimeout(() => {
                this.loading = false;
            }, 500);
        }
    };
}

// 从 web/views/demo/features.html 迁移的演示页面函数
function demoPage() {
    return {
        activeTab: 'htmx',
        showTransition: false,
        showCustomTransition: false,

        init() {
            console.log('✅ 特性演示页面已初始化');
        }
    };
}

// 从 web/views/layout.html 迁移的全局应用状态和响应式处理
// Alpine.js 全局应用状态
function initApp() {
    return {
        sidebarOpen: window.innerWidth >= 1024,
        activeMenu: window.location.pathname,

        init() {
            this.$nextTick(() => {
                this.activeMenu = window.location.pathname;
            });
        },

        setActiveMenu(path) {
            this.activeMenu = path;
        },

        handleResize() {
            if (window.innerWidth >= 1024) {
                this.sidebarOpen = true;
            }
        }
    }
}

// 响应式处理
document.addEventListener('DOMContentLoaded', function () {
    window.addEventListener('resize', () => {
        if (window.app && window.app.handleResize) {
            window.app.handleResize();
        }
    });
});

// Tailwind配置（从layout.html迁移）
if (window.tailwind) {
    tailwind.config = {
        plugins: [require('daisyui')],
        daisyui: {
            themes: [
                "light",
                "dark",
                "cupcake",
                "bumblebee",
                "emerald",
                "corporate",
                "synthwave",
                {
                    "godash": {
                        "primary": "#3b82f6",
                        "primary-focus": "#2563eb",
                        "primary-content": "#ffffff",
                        "secondary": "#64748b",
                        "secondary-focus": "#475569",
                        "secondary-content": "#ffffff",
                        "accent": "#14b8a6",
                        "accent-focus": "#0d9488",
                        "accent-content": "#ffffff",
                        "neutral": "#374151",
                        "neutral-focus": "#1f2937",
                        "neutral-content": "#ffffff",
                        "base-100": "#ffffff",
                        "base-200": "#f8fafc",
                        "base-300": "#f1f5f9",
                        "base-content": "#1e293b",
                        "info": "#3b82f6",
                        "success": "#10b981",
                        "warning": "#f59e0b",
                        "error": "#ef4444",
                    }
                }
            ],
            darkTheme: "dark",
            base: true,
            styled: true,
            utils: true,
            prefix: "",
            logs: true,
            themeRoot: ":root"
        }
    }
}

