package handler

import (
	"github.com/projuktisheba/pse-api-v1/internal/service"
)

type HandlerRepository struct {
	AuthHandler     *AuthHandler
	CategoryHandler *CategoryHandler
	ProductHandler *ProductHandler
}

func NewHandlerRepository(svc *service.ServiceRepository) *HandlerRepository {
	return &HandlerRepository{
		AuthHandler:     NewAuthHandler(svc.AuthService),
		CategoryHandler: NewCategoryHandler(svc.CategoryService),
		ProductHandler: NewProductHandler(svc.ProductService),
	}
}
