const fs = require('fs');
const path = '../ps-ecommerce-frontend/index.html';
let content = fs.readFileSync(path, 'utf8');

// remove all
content = content.replace(/<li class="delivery-link" style="display: none;"><a href="#\/delivery-portal">DELIVERY PORTAL<\/a><\/li>/g, '');
// replace back
const deliveryLink = `<li class="delivery-link" style="display: none;"><a href="#/delivery-portal">DELIVERY PORTAL</a></li>`;
content = content.replace(/<li><a href="#\/shop">SHOP<\/a><\/li>\s*/g, '<li><a href="#/shop">SHOP</a></li>\n                            ' + deliveryLink + '\n');
fs.writeFileSync(path, content);
