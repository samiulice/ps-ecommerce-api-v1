#!/bin/bash

cat internal/service/delivery_service.go | sed '/func (s \*ServiceRepository) ListDeliveryMen/,$d' > internal/service/delivery_service.go.tmp
mv internal/service/delivery_service.go.tmp internal/service/delivery_service.go

cat << 'INNER_EOF' >> internal/service/delivery_service.go

func (s *ServiceRepository) ListDeliveryMen(ctx context.Context) ([]model.DeliveryMan, error) {
        return s.db.ListDeliveryMen(ctx)
}

func (s *ServiceRepository) ListDeliveryMethods(ctx context.Context) ([]model.DeliveryMethod, error) {
        return s.db.ListDeliveryMethods(ctx)
}

func (s *ServiceRepository) GetDeliveryHistory(ctx context.Context, page, limit int) ([]repository.OrderDeliveryHistory, error) {
        if page < 1 {
                page = 1
        }
        if limit < 1 || limit > 100 {
                limit = 20
        }
        offset := (page - 1) * limit
        return s.db.GetDeliveryHistory(ctx, limit, offset)
}
INNER_EOF
