const fs = require('fs');
const path = '../ps-ecommerce-frontend/admin/js/common/sidebar.js';
let content = fs.readFileSync(path, 'utf8');

const deliveryMenu = `
  {
    type: "header",
    label: "Delivery Management",
  },
  {
    type: "parent",
    label: "Delivery",
    icon: "ph ph-truck",
    pageId: "delivery-menu",
    id: "delivery-menu",
    roles: ["chairman", "manager"],
    permissions: ["user.view"],
    children: [
      {
        label: "Delivery Men",
        pageId: "delivery-men",
        action: "loadPage('delivery-men','Manage Delivery Men')",
        roles: ["chairman", "manager"],
        permissions: ["user.view"],
      },
      {
        label: "Delivery Methods",
        pageId: "delivery-methods",
        action: "loadPage('delivery-methods','Manage Delivery Methods')",
        roles: ["chairman", "manager"],
        permissions: ["settings.view"],
      }
    ]
  },
`;

content = content.replace('label: "Product Management",\n  },', `label: "Product Management",\n  },${deliveryMenu}`);
fs.writeFileSync(path, content);
