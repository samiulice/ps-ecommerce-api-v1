// Package middleware contains HTTP middleware components.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const (
	customerIDKey   ctxKey = "customerID"
	customerTypeKey ctxKey = "customerType"
	roleKey         ctxKey = "role"
	roleIDKey       ctxKey = "roleID"
	permissionsKey  ctxKey = "permissions"
)

// JWTAuth validates JWT access tokens and injects customer ID, type, and role into context.
func JWTAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			if h == "" {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(h, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid token format", http.StatusUnauthorized)
				return
			}

			token, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			uid, ok := claims["sub"].(float64)
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			// Build context with customer ID
			ctx := context.WithValue(r.Context(), customerIDKey, int(uid))

			// Extract customer type (employee or customer)
			if customerType, ok := claims["type"].(string); ok {
				ctx = context.WithValue(ctx, customerTypeKey, customerType)
			}

			// Extract role (for employees)
			if role, ok := claims["role"].(string); ok {
				ctx = context.WithValue(ctx, roleKey, role)
			}

			if roleID, ok := claims["role_id"].(float64); ok {
				ctx = context.WithValue(ctx, roleIDKey, int64(roleID))
			}

			if rawPermissions, ok := claims["permissions"].([]interface{}); ok {
				permissions := make([]string, 0, len(rawPermissions))
				for _, item := range rawPermissions {
					if key, ok := item.(string); ok && strings.TrimSpace(key) != "" {
						permissions = append(permissions, key)
					}
				}
				ctx = context.WithValue(ctx, permissionsKey, permissions)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// CustomerIDFromContext extracts authenticated customer ID.
func CustomerIDFromContext(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(customerIDKey).(int)
	return id, ok
}

// CustomerTypeFromContext extracts customer type (employee or customer).
func CustomerTypeFromContext(ctx context.Context) (string, bool) {
	customerType, ok := ctx.Value(customerTypeKey).(string)
	return customerType, ok
}

// RoleFromContext extracts customer role (for employees).
func RoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(roleKey).(string)
	return role, ok
}

func RoleIDFromContext(ctx context.Context) (int64, bool) {
	roleID, ok := ctx.Value(roleIDKey).(int64)
	return roleID, ok
}

func PermissionsFromContext(ctx context.Context) ([]string, bool) {
	permissions, ok := ctx.Value(permissionsKey).([]string)
	return permissions, ok
}

// IsEmployee checks if the authenticated customer is an employee.
func IsEmployee(ctx context.Context) bool {
	customerType, ok := CustomerTypeFromContext(ctx)
	return ok && customerType == "employee"
}

// IsCustomer checks if the authenticated customer is a customer.
func IsCustomer(ctx context.Context) bool {
	customerType, ok := CustomerTypeFromContext(ctx)
	return ok && customerType == "customer"
}

// RequireEmployee middleware ensures only employees can access the route.
func RequireEmployee(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsEmployee(r.Context()) {
			http.Error(w, "forbidden: employee access required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireCustomer middleware ensures only customers can access the route.
func RequireCustomer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsCustomer(r.Context()) {
			http.Error(w, "forbidden: customer access required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireRole middleware ensures the employee has a specific role.
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := RoleFromContext(r.Context())
			if !ok {
				http.Error(w, "forbidden: role required", http.StatusForbidden)
				return
			}

			for _, allowed := range roles {
				if role == allowed {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "forbidden: insufficient permissions", http.StatusForbidden)
		})
	}
}

func RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			permissions, ok := PermissionsFromContext(r.Context())
			if !ok {
				http.Error(w, "forbidden: permissions required", http.StatusForbidden)
				return
			}

			for _, key := range permissions {
				if key == permission {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "forbidden: insufficient permissions", http.StatusForbidden)
		})
	}
}

func RequireAnyPermission(permissionKeys ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			permissions, ok := PermissionsFromContext(r.Context())
			if !ok {
				http.Error(w, "forbidden: permissions required", http.StatusForbidden)
				return
			}

			for _, current := range permissions {
				for _, allowed := range permissionKeys {
					if current == allowed {
						next.ServeHTTP(w, r)
						return
					}
				}
			}

			http.Error(w, "forbidden: insufficient permissions", http.StatusForbidden)
		})
	}
}
