const fs = require('fs');
const file = '../ps-ecommerce-frontend/admin/js/common/sidebar.js';
let content = fs.readFileSync(file, 'utf8');

content = content.replace(/(\{\s*label:\s*"Delivery History"[\s\S]*?\},)\s*\{\s*label:\s*"Delivery History"[\s\S]*?\},\s*/g, "$1");

fs.writeFileSync(file, content);
