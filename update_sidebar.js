const fs = require('fs');
const file = '../ps-ecommerce-frontend/admin/js/common/sidebar.js';
let content = fs.readFileSync(file, 'utf8');

const newItem = `
      {
        label: "Delivery History",
        pageId: "delivery_history",
        action: "loadPage('delivery_history', 'Delivery History')",
        roles: ["chairman", "manager"],
        permissions: ["report.view"],
      },`;

content = content.replace(/(label:\s*"Online Orders",\s*pageId:\s*"online_orders_report",\s*action:\s*"loadPage\('online_orders_report', 'Online Orders Report'\)",\s*roles:\s*\["chairman", "manager"\],\s*permissions:\s*\["report\.view"\],\s*},)/g, "$1" + newItem);

fs.writeFileSync(file, content);
