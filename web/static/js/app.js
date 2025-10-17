/**
 * Admin Dashboard - 全局 JavaScript
 * 包含 Alpine.js 组件和工具函数
 */

// ============================================
// Alpine.js 组件注册
// ============================================

document.addEventListener('alpine:init', () => {

    // ============================================
    // 删除确认组件（通用）
    // ============================================
    Alpine.data('deleteConfirm', (options) => ({
        confirmDelete() {
            window.dispatchEvent(new CustomEvent('confirm-dialog', {
                detail: {
                    title: options.title || '确认删除',
                    message: options.message,
                    confirmText: options.confirmText || '确定删除',
                    onConfirm: () => {
                        htmx.ajax('DELETE', options.url, {
                            target: options.target,
                            swap: 'outerHTML swap:300ms'
                        });
                    }
                }
            }));
        }
    }));

    // ============================================
    // 用户表单组件
    // ============================================
    Alpine.data('userForm', (isEdit = false) => ({
        loading: false,
        errors: {},
        form: {
            username: '',
            email: '',
            realName: '',
            phone: '',
            role: 'viewer',
            status: 'active'
        },

        validateForm() {
            this.errors = {};
            let isValid = true;

            if (!this.form.username && !isEdit) {
                this.errors.username = '用户名不能为空';
                isValid = false;
            } else if (this.form.username && !/^[a-zA-Z0-9_]{3,20}$/.test(this.form.username)) {
                this.errors.username = '用户名只能包含字母、数字和下划线，长度3-20位';
                isValid = false;
            }

            if (!this.form.email) {
                this.errors.email = '邮箱不能为空';
                isValid = false;
            } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(this.form.email)) {
                this.errors.email = '邮箱格式不正确';
                isValid = false;
            }

            if (!this.form.realName) {
                this.errors.realName = '真实姓名不能为空';
                isValid = false;
            }

            return isValid;
        },

        submitForm(event) {
            if (!this.validateForm()) {
                event.preventDefault();
                return false;
            }
            this.loading = true;
            // 验证通过，让 HTMX 处理表单提交
        }
    }));

    // ============================================
    // 商品表单组件
    // ============================================
    Alpine.data('productForm', () => ({
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

            if (!this.form.sku) {
                this.errors.sku = 'SKU不能为空';
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

            return isValid;
        },

        submitForm(event) {
            if (!this.validateForm()) {
                event.preventDefault();
                return false;
            }
            this.loading = true;
            // 验证通过，让 HTMX 处理表单提交
        }
    }));

    // ============================================
    // 设置表单组件
    // ============================================
    Alpine.data('settingsForm', () => ({
        loading: false,
        errors: {},
        form: {
            siteName: '',
            siteUrl: '',
            email: '',
            description: '',
            enableRegistration: true,
            enableComments: true,
            itemsPerPage: 10,
            maintenanceMode: false
        },

        validateForm() {
            this.errors = {};
            let isValid = true;

            if (!this.form.siteName) {
                this.errors.siteName = '网站名称不能为空';
                isValid = false;
            }

            if (!this.form.email) {
                this.errors.email = '管理员邮箱不能为空';
                isValid = false;
            } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(this.form.email)) {
                this.errors.email = '邮箱格式不正确';
                isValid = false;
            }

            return isValid;
        },

        submitForm(event) {
            if (this.validateForm()) {
                this.loading = true;
                htmx.trigger(event.target, 'submit');
            }
        }
    }));
});

// ============================================
// HTMX 事件监听
// ============================================

// 等待 DOM 加载完成后注册事件监听器
document.addEventListener('DOMContentLoaded', function () {
    console.log('Toast 监听器已注册'); // 调试日志

    // Toast 去重缓存（防止同一消息短时间内重复显示）
    let lastToastMessage = '';
    let lastToastTime = 0;

    // 监听 HTMX 交换完成后的 Toast 消息（使用 afterSwap，响应头此时还可用）
    document.body.addEventListener('htmx:afterSwap', function (event) {
        const xhr = event.detail.xhr;

        // 检查是否有 XHR 对象（有些 swap 可能没有）
        if (!xhr) {
            return;
        }

        // 检查响应头中的 Toast 消息
        const toastMessage = xhr.getResponseHeader('X-Toast-Message');
        const toastType = xhr.getResponseHeader('X-Toast-Type') || 'success';

        // 只有当真正有 Toast 消息时才处理
        if (toastMessage) {
            console.log('🔔 Toast 原始消息:', toastMessage, 'Type:', toastType);

            // 去重：如果是相同消息且在500毫秒内，忽略
            const now = Date.now();
            if (toastMessage === lastToastMessage && (now - lastToastTime) < 500) {
                console.log('⏭️ 跳过重复的 Toast 消息');
                return;
            }

            // 更新缓存
            lastToastMessage = toastMessage;
            lastToastTime = now;

            // 解码 URL 编码的中文消息
            let decodedMessage = toastMessage;
            try {
                decodedMessage = decodeURIComponent(toastMessage);
                console.log('✅ Toast 解码后:', decodedMessage);
            } catch (e) {
                console.error('❌ 解码消息失败:', e);
            }

            window.dispatchEvent(new CustomEvent('show-toast', {
                detail: {
                    message: decodedMessage,
                    type: toastType
                }
            }));
        }
    });

    // 监听 HTMX 错误
    document.body.addEventListener('htmx:responseError', function (event) {
        window.dispatchEvent(new CustomEvent('show-toast', {
            detail: {
                message: '请求失败，请稍后重试',
                type: 'error'
            }
        }));
    });
});

// ============================================
// 工具函数
// ============================================

/**
 * 格式化日期
 * @param {Date|string} date - 日期对象或字符串
 * @param {string} format - 格式，默认 'YYYY-MM-DD HH:mm:ss'
 * @returns {string} 格式化后的日期字符串
 */
function formatDate(date, format = 'YYYY-MM-DD HH:mm:ss') {
    const d = new Date(date);
    const year = d.getFullYear();
    const month = String(d.getMonth() + 1).padStart(2, '0');
    const day = String(d.getDate()).padStart(2, '0');
    const hours = String(d.getHours()).padStart(2, '0');
    const minutes = String(d.getMinutes()).padStart(2, '0');
    const seconds = String(d.getSeconds()).padStart(2, '0');

    return format
        .replace('YYYY', year)
        .replace('MM', month)
        .replace('DD', day)
        .replace('HH', hours)
        .replace('mm', minutes)
        .replace('ss', seconds);
}

/**
 * 格式化货币
 * @param {number} amount - 金额
 * @param {string} currency - 货币符号，默认 '¥'
 * @returns {string} 格式化后的货币字符串
 */
function formatCurrency(amount, currency = '¥') {
    return `${currency}${parseFloat(amount).toFixed(2)}`;
}

/**
 * 防抖函数
 * @param {Function} func - 要执行的函数
 * @param {number} wait - 等待时间（毫秒）
 * @returns {Function} 防抖后的函数
 */
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

// 将工具函数挂载到 window 对象（可选）
window.utils = {
    formatDate,
    formatCurrency,
    debounce
};

