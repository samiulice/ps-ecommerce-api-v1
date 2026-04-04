const fs = require('fs');
const path = '../ps-ecommerce-frontend/admin/js/pages/order_list.js';
let content = fs.readFileSync(path, 'utf8');

const assignBtnHtml = `
        <button onclick='openAssignModal(\${order.id})'
          class="text-indigo-600 hover:text-indigo-700 hover:bg-indigo-100 transition-all p-1.5 rounded-lg"
          title="Assign Delivery Man">
          <i class="fa-solid fa-truck text-lg"></i>
        </button>
`;

content = content.replace(/<button onclick='printInvoice\(\$\{order\.id\}\)'/, assignBtnHtml + "\n        <button onclick='printInvoice(${order.id})'");

fs.writeFileSync(path, content);
