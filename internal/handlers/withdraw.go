package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/ShadyZiedan/gophermart/internal/luhn"
	"github.com/ShadyZiedan/gophermart/internal/security"
	"github.com/ShadyZiedan/gophermart/internal/services"
)

type withdrawReq struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(security.UserIDKey{}).(int)
	var req withdrawReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderNumber, err := strconv.Atoi(req.Order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !luhn.Valid(orderNumber) {
		http.Error(w, "Order number is not valid", http.StatusUnprocessableEntity)
		return
	}

	if err = h.balanceService.WithdrawOrder(ctx, userID, orderNumber, req.Sum); err != nil {
		if errors.Is(err, services.ErrInsufficientBalance) {
			http.Error(w, err.Error(), http.StatusPaymentRequired)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
