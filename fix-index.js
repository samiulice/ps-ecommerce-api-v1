const fs = require('fs');
const path = '../ps-ecommerce-frontend/index.html';
let content = fs.readFileSync(path, 'utf8');

const deliveryLink = `\n                            <li class="delivery-link" style="display: none;"><a href="#/delivery-portal">DELIVERY PORTAL</a></li>`;

content = content.replace('<li><a href="#/shop">SHOP</a></li>', `<li><a href="#/shop">SHOP</a></li>${deliveryLink}`);
fs.writeFileSync(path, content);
