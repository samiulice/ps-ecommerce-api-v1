package routes

import (
	"net/http"

	"github.com/projuktisheba/pse-api-v1/internal/middleware"
)

func employeeAuth(secretKey string) func(http.Handler) http.Handler {
	return middleware.JWTAuth(secretKey)
}
