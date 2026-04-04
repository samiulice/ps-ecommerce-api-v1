sed -i 's/utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch delivery men")/utils.ServerError(w, err)/g' internal/handler/delivery_handler.go
sed -i 's/utils.SuccessResponse(w, http.StatusOK, "Delivery men retrieved", men)/utils.OK(w, "Delivery men retrieved", men)/g' internal/handler/delivery_handler.go
sed -i 's/utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch delivery history")/utils.ServerError(w, err)/g' internal/handler/delivery_handler.go
sed -i 's/utils.SuccessResponse(w, http.StatusOK, "Delivery history retrieved", history)/utils.OK(w, "Delivery history retrieved", history)/g' internal/handler/delivery_handler.go
