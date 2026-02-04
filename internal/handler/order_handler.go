package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/service"
	"github.com/projuktisheba/pse-api-v1/pkg/utils"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) handleErr(w http.ResponseWriter, err error) {
	fmt.Println("Order Error: ", err)
	if strings.Contains(err.Error(), "not found") {
		utils.NotFound(w, err)
	} else if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "required") {
		utils.BadRequest(w, err)
	} else {
		utils.ServerError(w, err)
	}
}

// PlaceOrder handles POST /orders/new - public endpoint for checkout
func (h *OrderHandler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var req model.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, err)
		return
	}
	fmt.Printf("%+v\n", req)
	response, err := h.orderService.PlaceOrder(r.Context(), req)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	utils.Created(w, "Order placed successfully", map[string]interface{}{
		"order":        response,
		"order_number": response.OrderNumber,
	})
}

// GetOrder handles GET /orders/{id}
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.BadRequest(w, fmt.Errorf("invalid order ID"))
		return
	}

	order, err := h.orderService.GetOrderByID(r.Context(), id)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	utils.OK(w, "Order retrieved successfully", map[string]interface{}{"order": order})
}

// GetOrderByNumber handles GET /orders/number/{orderNumber}
func (h *OrderHandler) GetOrderByNumber(w http.ResponseWriter, r *http.Request) {
	orderNumber := chi.URLParam(r, "orderNumber")
	if orderNumber == "" {
		utils.BadRequest(w, fmt.Errorf("order number is required"))
		return
	}

	order, err := h.orderService.GetOrderByOrderNumber(r.Context(), orderNumber)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	utils.OK(w, "Order retrieved successfully", map[string]interface{}{"order": order})
}

// ListOrders handles GET /orders
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	filter := model.OrderFilter{
		Status:        query.Get("status"),
		PaymentStatus: query.Get("payment_status"),
		PaymentMethod: query.Get("payment_method"),
		FromDate:      query.Get("from_date"),
		ToDate:        query.Get("to_date"),
		Search:        query.Get("search"),
	}

	if customerID := query.Get("customer_id"); customerID != "" {
		if id, err := strconv.ParseInt(customerID, 10, 64); err == nil {
			filter.CustomerID = id
		}
	}

	page, _ := strconv.Atoi(query.Get("page"))
	if page <= 0 {
		page = 1
	}
	filter.Page = page

	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit <= 0 {
		limit = 20
	}
	filter.Limit = limit

	orders, total, err := h.orderService.ListOrders(r.Context(), filter)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	totalPages := (int(total) + limit - 1) / limit

	utils.OK(w, "Orders retrieved successfully", map[string]interface{}{
		"orders":      orders,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	})
}

// GetCustomerOrders handles GET /orders/customer/{customerId}
func (h *OrderHandler) GetCustomerOrders(w http.ResponseWriter, r *http.Request) {
	customerIDStr := chi.URLParam(r, "customerId")
	customerID, err := strconv.ParseInt(customerIDStr, 10, 64)
	if err != nil {
		utils.BadRequest(w, fmt.Errorf("invalid customer ID"))
		return
	}

	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	if page <= 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit <= 0 {
		limit = 20
	}

	orders, total, err := h.orderService.GetCustomerOrders(r.Context(), customerID, page, limit)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	totalPages := (int(total) + limit - 1) / limit

	utils.OK(w, "Customer orders retrieved successfully", map[string]interface{}{
		"orders":      orders,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	})
}

// UpdateOrderStatus handles PUT /orders/{id}/status
func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.BadRequest(w, fmt.Errorf("invalid order ID"))
		return
	}

	var req model.UpdateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, err)
		return
	}

	if req.Status == "" {
		utils.BadRequest(w, fmt.Errorf("status is required"))
		return
	}

	if err := h.orderService.UpdateOrderStatus(r.Context(), id, req.Status, req.Reason); err != nil {
		h.handleErr(w, err)
		return
	}

	utils.OK(w, "Order status updated successfully", nil)
}

// UpdatePaymentStatus handles PUT /orders/{id}/payment-status
func (h *OrderHandler) UpdatePaymentStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.BadRequest(w, fmt.Errorf("invalid order ID"))
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, err)
		return
	}

	if err := h.orderService.UpdatePaymentStatus(r.Context(), id, req.Status); err != nil {
		h.handleErr(w, err)
		return
	}

	utils.OK(w, "Payment status updated successfully", nil)
}

// DeleteOrder handles DELETE /orders/{id}
func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.BadRequest(w, fmt.Errorf("invalid order ID"))
		return
	}

	if err := h.orderService.DeleteOrder(r.Context(), id); err != nil {
		h.handleErr(w, err)
		return
	}

	utils.OK(w, "Order deleted successfully", nil)
}

// GetOrderStats handles GET /orders/stats
func (h *OrderHandler) GetOrderStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.orderService.GetOrderStats(r.Context())
	if err != nil {
		h.handleErr(w, err)
		return
	}

	utils.OK(w, "Order stats retrieved successfully", map[string]interface{}{"stats": stats})
}
