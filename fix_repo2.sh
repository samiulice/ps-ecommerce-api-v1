#!/bin/bash

# Remove the broken `DBRepository` methods using awk/sed
sed -i '/func (r \*DBRepository) ListDeliveryMen/,$d' internal/repository/delivery_repository.go

cat << 'INNER_EOF' >> internal/repository/delivery_repository.go

func (r *DeliveryRepository) ListDeliveryMen(ctx context.Context) ([]model.DeliveryMan, error) {
        query := `
                SELECT 
                        dm.id, dm.customer_id, c.name, c.phone, dm.is_active, dm.is_online, dm.vehicle_type, dm.vehicle_number
                FROM delivery_men dm
                JOIN customers c ON dm.customer_id = c.id
                ORDER BY c.name ASC
        `
        rows, err := r.db.Query(ctx, query)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var men []model.DeliveryMan
        for rows.Next() {
                var d model.DeliveryMan
                if err := rows.Scan(&d.ID, &d.CustomerID, &d.CustomerName, &d.CustomerPhone, &d.IsActive, &d.IsOnline, &d.VehicleType, &d.VehicleNumber); err != nil {
                        return nil, err
                }
                men = append(men, d)
        }
        return men, nil
}

type OrderDeliveryHistory struct {
        model.OrderDelivery
        DeliveryManName  *string `json:"delivery_man_name,omitempty"`
        DeliveryManPhone *string `json:"delivery_man_phone,omitempty"`
}

func (r *DeliveryRepository) GetDeliveryHistory(ctx context.Context, limit, offset int) ([]OrderDeliveryHistory, error) {
        query := `
                SELECT 
                        od.id, od.order_id, od.delivery_man_id, c.name, c.phone, od.delivery_status, 
                        od.delivery_fee_collected, od.delivery_man_earning, od.assigned_at, od.delivered_at, od.created_at, od.updated_at
                FROM order_deliveries od
                LEFT JOIN delivery_men dm ON od.delivery_man_id = dm.id
                LEFT JOIN customers c ON dm.customer_id = c.id
                ORDER BY od.created_at DESC
                LIMIT $1 OFFSET $2
        `
        rows, err := r.db.Query(ctx, query, limit, offset)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var history []OrderDeliveryHistory
        for rows.Next() {
                var d OrderDeliveryHistory
                if err := rows.Scan(&d.ID, &d.OrderID, &d.DeliveryManID, &d.DeliveryManName, &d.DeliveryManPhone, &d.DeliveryStatus, 
                        &d.DeliveryFeeCollected, &d.DeliveryManEarning, &d.AssignedAt, &d.DeliveredAt, &d.CreatedAt, &d.UpdatedAt); err != nil {
                        return nil, err
                }
                history = append(history, d)
        }
        return history, nil
}
INNER_EOF
