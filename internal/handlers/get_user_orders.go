package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type orderResponseModel struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual,omitempty"`
	UploadedAt string  `json:"uploaded_at"`
}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(int)
	orders, err := h.orderService.GetOrders(r.Context(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var response []orderResponseModel
	for _, order := range orders {
		response = append(response, orderResponseModel{
			Number:     strconv.Itoa(order.Number),
			Status:     order.Status,
			Accrual:    order.Accrual,
			UploadedAt: order.UploadedAt.Format(time.RFC3339),
		})
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
