const fs = require('fs');
const path = 'e:/ongoing-projects/projuktisheba/ps-ecommerce/ps-ecommerce-frontend/admin/js/pages/reports.js';
let content = fs.readFileSync(path, 'utf8');

const oldCode = `// Initial load
window.fetchReportData();`;

const newCode = `// Initial load
if (window.pendingReportTab) {
    window.loadReport(window.pendingReportTab);
    window.pendingReportTab = null;
} else {
    window.loadReport(window.reportState.currentTab);
}`;

content = content.replace(oldCode, newCode);
fs.writeFileSync(path, content, 'utf8');
