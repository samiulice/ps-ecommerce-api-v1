const fs = require('fs');
const path = 'e:/ongoing-projects/projuktisheba/ps-ecommerce/ps-ecommerce-frontend/admin/js/common/sidebar.js';
let content = fs.readFileSync(path, 'utf8');

const regex = /\{\s*type:\s*"link",\s*label:\s*"Reports",[\s\S]*?permissions:\s*\["report\.view"\],\s*\},/;
if (regex.test(content)) {
    console.log('Match found!');
}

const replacement = `{
    type: "header",
    label: "Reports",
  },
  {
    type: "parent",
    label: "Reports",
    icon: "ph ph-chart-line-up",
    pageId: "reports-menu",
    id: "reports-menu",
    roles: ["chairman", "manager"],
    permissions: ["report.view"],
    children: [
      {
        label: "POS Sales",
        pageId: "reports",
        action: "window.pendingReportTab = 'pos-sales'; loadPage('reports','POS Sales Report')",
        roles: ["chairman", "manager"],
        permissions: ["report.view"],
      },
      {
        label: "Online Orders",
        pageId: "reports",
        action: "window.pendingReportTab = 'orders'; loadPage('reports','Online Orders Report')",
        roles: ["chairman", "manager"],
        permissions: ["report.view"],
      },
      {
        label: "Customer Due",
        pageId: "reports",
        action: "window.pendingReportTab = 'customer-dues'; loadPage('reports','Customer Due Report')",
        roles: ["chairman", "manager"],
        permissions: ["report.view"],
      },
      {
        label: "Supplier Due",
        pageId: "reports",
        action: "window.pendingReportTab = 'supplier-dues'; loadPage('reports','Supplier Due Report')",
        roles: ["chairman", "manager"],
        permissions: ["report.view"],
      },
      {
        label: "Low Stock",
        pageId: "reports",
        action: "window.pendingReportTab = 'low-stock'; loadPage('reports','Low Stock Report')",
        roles: ["chairman", "manager"],
        permissions: ["report.view"],
      }
    ]
  },`;

content = content.replace(regex, replacement);
fs.writeFileSync(path, content, 'utf8');
