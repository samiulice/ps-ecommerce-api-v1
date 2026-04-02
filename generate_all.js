const fs = require('fs');

function buildHtml(id, searchPlaceholder, filters, tableHeaders, includeAmount) {
    return `<div class="flex flex-col h-full fade-in space-y-4">
  <div class="bg-white p-4 rounded-xl shadow-sm border border-slate-200">
    <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
      <div class="flex flex-wrap items-center gap-2 lg:gap-3 w-full">
${filters}
        <div class="relative w-full md:flex-1">
            <i class="ph ph-magnifying-glass absolute left-3 top-1/2 -translate-y-1/2 text-slate-400"></i>
            <input type="text" id="${id}-search" placeholder="${searchPlaceholder}" class="w-full pl-10 pr-3 py-2.5 border border-slate-300 rounded-lg text-sm focus:ring-2 focus:ring-brand-500 outline-none bg-slate-50 focus:bg-white">
        </div>
      </div>
    </div>
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mt-6">
        <div class="bg-indigo-50 border border-indigo-100 rounded-xl p-4 flex items-center justify-between">
            <div>
                <p class="text-indigo-600 text-xs font-bold uppercase tracking-wider mb-1">Total Hits</p>
                <h4 class="text-2xl font-black text-indigo-900" id="summary-hits">0</h4>
            </div>
            <div class="w-12 h-12 rounded-full bg-indigo-100 text-indigo-600 flex items-center justify-center text-xl shadow-inner"><i class="ph ph-list-numbers"></i></div>
        </div>
${includeAmount ? `        <div class="bg-emerald-50 border border-emerald-100 rounded-xl p-4 flex items-center justify-between" id="summary-amount-card">
            <div>
                <p class="text-emerald-700 text-xs font-bold uppercase tracking-wider mb-1" id="summary-amount-label">Total Amount</p>
                <h4 class="text-2xl font-black text-emerald-900" id="summary-amount">$0.00</h4>
            </div>
            <div class="w-12 h-12 rounded-full bg-emerald-100 text-emerald-600 flex items-center justify-center text-xl shadow-inner"><i class="ph ph-currency-dollar"></i></div>
        </div>` : ''}
    </div>
  </div>
  <div class="flex-1 bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden relative flex flex-col">
    <div class="overflow-x-auto flex-1 p-0">
        <table class="w-full text-left border-collapse">
            <thead class="bg-slate-50 border-b border-slate-200 text-slate-500 text-xs uppercase sticky top-0 z-10">
                <tr>${tableHeaders}</tr>
            </thead>
            <tbody id="${id}-table-body" class="divide-y divide-slate-100 text-sm"></tbody>
        </table>
        <div id="${id}-loading" class="absolute inset-0 bg-white/80 backdrop-blur-sm z-10 flex flex-col items-center justify-center hidden">
            <i class="ph ph-spinner animate-spin text-4xl text-brand-600 mb-2"></i><p class="text-sm font-semibold text-slate-600 animate-pulse">Loading...</p>
        </div>
        <div id="${id}-empty" class="py-16 flex flex-col items-center justify-center text-slate-400 hidden">
            <i class="ph ph-folder-open text-6xl mb-3 text-slate-200"></i><p class="text-slate-500 font-medium">No records found matching criteria</p>
        </div>
    </div>
    <div class="p-4 border-t border-slate-200 bg-slate-50/50 flex flex-col sm:flex-row items-center justify-between gap-4 mt-auto">
        <p class="text-sm text-slate-500" id="${id}-pagination-info">Showing page 1</p>
        <div class="flex items-center gap-1">
            <button onclick="window.${id.replace(/-/g, '')}State.changePage(-1)" class="w-8 h-8 rounded border border-slate-200 bg-white text-slate-600 hover:text-brand-600 transition flex items-center justify-center"><i class="ph ph-caret-left"></i></button>
            <span id="${id}-current-page" class="w-8 h-8 rounded bg-brand-600 text-white font-bold flex items-center justify-center text-sm shadow-sm">1</span>
            <button onclick="window.${id.replace(/-/g, '')}State.changePage(1)" class="w-8 h-8 rounded border border-slate-200 bg-white text-slate-600 hover:text-brand-600 transition flex items-center justify-center"><i class="ph ph-caret-right"></i></button>
        </div>
    </div>
  </div>
</div>
<script src="js/pages/${id.replace(/-/g, '_')}_report.js"></script>`;
}

function buildJs(id, endpoint, hasAmountCard, amountType, amountField, rowRenderHtml) {
    const jsStateName = `window.${id.replace(/-/g, '')}State`;
    return `${jsStateName} = {
    page: 1, limit: 15, totalPages: 1,
    init() { this.bindEvents(); this.fetchData(); },
    bindEvents() {
        const typeFilter = document.getElementById('${id}-filter-type');
        const orderFilter = document.getElementById('${id}-filter-order');
        const searchInput = document.getElementById('${id}-search');
        if(typeFilter) typeFilter.addEventListener('change', () => { this.page = 1; this.fetchData(); });
        if(orderFilter) orderFilter.addEventListener('change', () => { this.page = 1; this.fetchData(); });
        if(searchInput) {
            let timeout = null;
            searchInput.addEventListener('keyup', (e) => {
                clearTimeout(timeout);
                timeout = setTimeout(() => { this.page = 1; this.fetchData(); }, 500);
            });
        }
    },
    async fetchData() {
        const loading = document.getElementById('${id}-loading');
        const empty = document.getElementById('${id}-empty');
        const tbody = document.getElementById('${id}-table-body');
        if(loading) loading.classList.remove('hidden');
        if(empty) empty.classList.add('hidden');
        
        const type = document.getElementById('${id}-filter-type')?.value;
        const order = document.getElementById('${id}-filter-order')?.value;
        const search = document.getElementById('${id}-search')?.value;
        const params = new URLSearchParams({ page: this.page, limit: this.limit });
        if (type) params.append('sale_type', type);
        if (order) params.append('order_by', order);
        if (search) params.append('search', search);

        try {
            const response = await fetch(\`\${window.API_BASE_URL}/api/v1/reports/${endpoint}?\${params.toString()}\`, {
                headers: window.getAuthHeaders()
            });
            if (!response.ok) throw new Error('Failed to fetch data');
            const data = await response.json();
            this.totalPages = data.total_pages || 1;
            
            const pageEl = document.getElementById('${id}-current-page');
            const infoEl = document.getElementById('${id}-pagination-info');
            if(pageEl) pageEl.textContent = this.page;
            if(infoEl) infoEl.textContent = \`Showing page \${this.page} of \${this.totalPages} (Total: \${data.total_items})\`;

            if (!data.items || data.items.length === 0) {
                if(tbody) tbody.innerHTML = '';
                if(empty) empty.classList.remove('hidden');
            } else {
                this.renderData(data.items, tbody);
            }
            
            const hits = document.getElementById('summary-hits');
            if(hits) hits.textContent = data.total_items || 0;
            
            ${hasAmountCard ? `
            const amountCard = document.getElementById('summary-amount-card');
            const amountLabel = document.getElementById('summary-amount-label');
            const amount = document.getElementById('summary-amount');
            if(amountCard && amountLabel && amount) {
                amountLabel.textContent = '${amountType}';
                let total = (data.items || []).reduce((sum, item) => sum + (item.${amountField} || 0), 0);
                amount.className = '${amountType.includes('Due') ? 'text-2xl font-black text-red-600' : 'text-2xl font-black text-emerald-900'}';
                amount.textContent = \`$\${total.toFixed(2)}\`;
            }` : ''}

        } catch (error) {
            console.error('Error:', error);
            if(tbody) tbody.innerHTML = '';
            if(empty) empty.classList.remove('hidden');
        } finally {
            if(loading) loading.classList.add('hidden');
        }
    },
    renderData(items, tbody) {
        if(!tbody) return;
        let html = '';
        items.forEach(item => {
            ${rowRenderHtml}
        });
        tbody.innerHTML = html;
    },
    changePage(delta) {
        const newPage = this.page + delta;
        if (newPage >= 1 && newPage <= this.totalPages) {
            this.page = newPage;
            this.fetchData();
        }
    }
};
${jsStateName}.init();`;
}

const dirHtml = 'e:/ongoing-projects/projuktisheba/ps-ecommerce/ps-ecommerce-frontend/admin/pages';
const dirJs = 'e:/ongoing-projects/projuktisheba/ps-ecommerce/ps-ecommerce-frontend/admin/js/pages';

// 1. ONLINE ORDERS
fs.writeFileSync(`${dirHtml}/online_orders_report.html`, buildHtml(
    'online-orders', 'Search Orders...',
    `<select id="online-orders-filter-type" class="p-2.5 border border-slate-300 rounded-lg text-sm bg-slate-50 focus:ring-2 focus:ring-brand-500 outline-none w-full md:w-auto"><option value="">All Types</option><option value="retail">Retail Only</option><option value="wholesale">Wholesale Only</option></select>
    <select id="online-orders-filter-order" class="p-2.5 border border-slate-300 rounded-lg text-sm bg-slate-50 focus:ring-2 focus:ring-brand-500 outline-none w-full md:w-auto"><option value="date_desc">Latest First</option><option value="price_desc">Highest Amount</option><option value="price_asc">Lowest Amount</option></select>`,
    `<th class="p-4 font-bold">Order ID</th><th class="p-4 font-bold">Date</th><th class="p-4 font-bold">Customer</th><th class="p-4 font-bold">Type</th><th class="p-4 font-bold text-right">Total</th><th class="p-4 font-bold text-center">Status</th><th class="p-4 font-bold text-center">Payment</th>`,
    true
));
fs.writeFileSync(`${dirJs}/online_orders_report.js`, buildJs('online-orders', 'orders', true, 'Page Total Amount', 'total_amount', 
    `html += \`<tr class="hover:bg-slate-50 transition-colors">
        <td class="p-4 font-medium text-slate-800">#\${item.id}</td>
        <td class="p-4 text-slate-500">\${new Date(item.created_at).toLocaleDateString()}</td>
        <td class="p-4">\${item.customer_name}</td>
        <td class="p-4"><span class="px-2 py-1 rounded-full text-xs font-semibold \${item.sale_type === 'wholesale' ? 'bg-purple-100 text-purple-700' : 'bg-blue-100 text-blue-700'}">\${item.sale_type || 'retail'}</span></td>
        <td class="p-4 text-right font-medium text-slate-800">$\${item.total_amount.toFixed(2)}</td>
        <td class="p-4 text-center"><span class="px-2 py-1 rounded-full text-xs font-semibold bg-slate-100 text-slate-700">\${item.status}</span></td>
        <td class="p-4 text-center"><span class="px-2 py-1 rounded-full text-xs font-semibold \${item.payment_status === 'paid' ? 'bg-emerald-100 text-emerald-700' : 'bg-amber-100 text-amber-700'}">\${item.payment_status}</span></td>
    </tr>\`;`
));

// 2. CUSTOMER DUE
fs.writeFileSync(`${dirHtml}/customer_due_report.html`, buildHtml(
    'customer-due', 'Search Customer...',
    `<select id="customer-due-filter-order" class="p-2.5 border border-slate-300 rounded-lg text-sm bg-slate-50 focus:ring-2 focus:ring-brand-500 outline-none w-full md:w-auto"><option value="date_desc">Date Desc</option><option value="price_desc">Highest Due</option><option value="price_asc">Lowest Due</option></select>`,
    `<th class="p-4 font-bold">Customer Name</th><th class="p-4 font-bold">Phone</th><th class="p-4 font-bold text-center">Total Sales</th><th class="p-4 font-bold text-center">Total Orders</th><th class="p-4 font-bold text-right">Total Due</th><th class="p-4 font-bold text-center">Action</th>`,
    true
));
fs.writeFileSync(`${dirJs}/customer_due_report.js`, buildJs('customer-due', 'customer-dues', true, 'Page Total Due', 'total_due', 
    `html += \`<tr class="hover:bg-slate-50 transition-colors">
        <td class="p-4 font-medium text-slate-800"><div class="flex items-center gap-3"><div class="w-8 h-8 rounded-full bg-brand-100 text-brand-600 flex items-center justify-center font-bold">\${item.name.charAt(0)}</div>\${item.name}</div></td>
        <td class="p-4 text-slate-500">\${item.phone || '-'}</td>
        <td class="p-4 text-center text-slate-600 font-medium">\${item.total_pos_sales || 0}</td>
        <td class="p-4 text-center text-slate-600 font-medium">\${item.total_orders || 0}</td>
        <td class="p-4 text-right font-bold text-red-600">$\${(item.total_due || 0).toFixed(2)}</td>
        <td class="p-4 text-center"><button class="text-brand-600 hover:text-brand-800 p-2 rounded hover:bg-brand-50 transition" title="View Customer Details"><i class="ph ph-eye text-lg"></i></button></td>
    </tr>\`;`
));

// 3. SUPPLIER DUE
fs.writeFileSync(`${dirHtml}/supplier_due_report.html`, buildHtml(
    'supplier-due', 'Search Supplier...',
    `<select id="supplier-due-filter-order" class="p-2.5 border border-slate-300 rounded-lg text-sm bg-slate-50 focus:ring-2 focus:ring-brand-500 outline-none w-full md:w-auto"><option value="date_desc">Date Desc</option><option value="price_desc">Highest Due</option><option value="price_asc">Lowest Due</option></select>`,
    `<th class="p-4 font-bold">Supplier Name</th><th class="p-4 font-bold">Contact Name</th><th class="p-4 font-bold">Phone</th><th class="p-4 font-bold text-right">Total Purchases</th><th class="p-4 font-bold text-right">Total Due</th><th class="p-4 font-bold text-center">Action</th>`,
    true
));
fs.writeFileSync(`${dirJs}/supplier_due_report.js`, buildJs('supplier-due', 'supplier-dues', true, 'Page Total Due', 'total_due', 
    `html += \`<tr class="hover:bg-slate-50 transition-colors">
        <td class="p-4 font-medium text-slate-800"><div class="flex items-center gap-3"><div class="w-8 h-8 rounded-full bg-slate-100 text-slate-600 flex items-center justify-center font-bold"><i class="ph ph-buildings"></i></div>\${item.name}</div></td>
        <td class="p-4 text-slate-500">\${item.contact_name || '-'}</td>
        <td class="p-4 text-slate-500">\${item.phone || '-'}</td>
        <td class="p-4 text-right text-slate-600 font-medium">\${item.total_purchases || 0}</td>
        <td class="p-4 text-right font-bold text-red-600">$\${(item.total_due || 0).toFixed(2)}</td>
        <td class="p-4 text-center"><button class="text-amber-600 hover:text-amber-800 p-2 rounded hover:bg-amber-50 transition" title="Make Payment"><i class="ph ph-money text-lg"></i></button></td>
    </tr>\`;`
));

// 4. LOW STOCK
fs.writeFileSync(`${dirHtml}/low_stock_report.html`, buildHtml(
    'low-stock', 'Search Product...',
    `<select id="low-stock-filter-order" class="p-2.5 border border-slate-300 rounded-lg text-sm bg-slate-50 focus:ring-2 focus:ring-brand-500 outline-none w-full md:w-auto"><option value="price_asc">Stock Asc</option><option value="price_desc">Stock Desc</option><option value="date_desc">Name A-Z</option></select>`,
    `<th class="p-4 font-bold">Product Name</th><th class="p-4 font-bold">SKU</th><th class="p-4 font-bold text-right">Current Stock</th><th class="p-4 font-bold text-right">Price</th><th class="p-4 font-bold text-center">Status</th>`,
    false
));
fs.writeFileSync(`${dirJs}/low_stock_report.js`, buildJs('low-stock', 'low-stock', false, '', '', 
    `html += \`<tr class="hover:bg-slate-50 transition-colors">
        <td class="p-4 font-medium text-slate-800 flex items-center gap-3">\${item.image_url ? \`<img src="\${item.image_url}" class="w-8 h-8 rounded object-cover border" alt="\${item.name}">\` : \`<div class="w-8 h-8 rounded bg-slate-100 flex items-center justify-center text-slate-400"><i class="ph ph-image"></i></div>\`}\${item.name}</td>
        <td class="p-4 text-slate-500">\${item.sku || '-'}</td>
        <td class="p-4 text-right font-bold \${item.stock_quantity === 0 ? 'text-red-600' : 'text-amber-600'}">\${item.stock_quantity}</td>
        <td class="p-4 text-right text-slate-600 font-medium">$\${(item.price || 0).toFixed(2)}</td>
        <td class="p-4 text-center"><span class="px-2 py-1 rounded-full text-xs font-semibold \${item.stock_quantity === 0 ? 'bg-red-100 text-red-700' : 'bg-amber-100 text-amber-700'}">\${item.stock_quantity === 0 ? 'Out of Stock' : 'Low Stock'}</span></td>
    </tr>\`;`
));
