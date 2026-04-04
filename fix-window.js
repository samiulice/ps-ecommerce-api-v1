const fs = require('fs');
const file = '../ps-ecommerce-frontend/admin/js/pages/order_list.js';
let content = fs.readFileSync(file, 'utf8');
content = content.replace("function openAssignModal(orderId)", "window.openAssignModal = function(orderId)");
content = content.replace("function closeAssignModal()", "window.closeAssignModal = function()");
fs.writeFileSync(file, content);
