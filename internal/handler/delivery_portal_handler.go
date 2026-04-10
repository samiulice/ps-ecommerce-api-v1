package handler

import (
"errors"
"net/http"

"github.com/projuktisheba/pse-api-v1/internal/middleware"
"github.com/projuktisheba/pse-api-v1/pkg/utils"
)

// GetPortalOrders handles GET /delivery/portal/orders
func (h *DeliveryHandler) GetPortalOrders(w http.ResponseWriter, r *http.Request) {
employeeID, ok := middleware.AuthIDFromContext(r.Context())
if !ok || !middleware.IsEmployee(r.Context()) {
utils.Unauthorized(w, errors.New("unauthorized delivery user"))
return
}

orders, err := h.svc.GetPortalOrders(r.Context(), int64(employeeID))
if err != nil {
utils.ServerError(w, err)
return
}
utils.OK(w, "Assigned orders retrieved", orders)
}

// GetPortalWallet handles GET /delivery/portal/wallet
func (h *DeliveryHandler) GetPortalWallet(w http.ResponseWriter, r *http.Request) {
employeeID, ok := middleware.AuthIDFromContext(r.Context())
if !ok || !middleware.IsEmployee(r.Context()) {
utils.Unauthorized(w, errors.New("unauthorized delivery user"))
return
}

wallet, err := h.svc.GetPortalWallet(r.Context(), int64(employeeID))
if err != nil {
utils.ServerError(w, err)
return
}
utils.OK(w, "Wallet details retrieved", wallet)
}
