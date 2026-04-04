sed -i 's/Customer ID/Employee ID/g' ../ps-ecommerce-frontend/admin/pages/delivery-men.html
sed -i "s/User's existing customer ID/User's existing employee ID/g" ../ps-ecommerce-frontend/admin/pages/delivery-men.html
sed -i 's/dm-customer-id/dm-employee-id/g' ../ps-ecommerce-frontend/admin/pages/delivery-men.html
sed -i 's/customer_id: parseInt(document.getElementById(.dm-customer-id.).value, 10)/employee_id: parseInt(document.getElementById("dm-employee-id").value, 10)/g' ../ps-ecommerce-frontend/admin/js/pages/delivery-men.js
sed -i "s/Customer successfully registered/Employee successfully registered/g" ../ps-ecommerce-frontend/admin/js/pages/delivery-men.js
sed -i 's/dm.customer_name/dm.employee_name/g' ../ps-ecommerce-frontend/admin/js/pages/order_list.js
sed -i 's/dm.customer_phone/dm.employee_mobile/g' ../ps-ecommerce-frontend/admin/js/pages/order_list.js
