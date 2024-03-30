package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/ShadyZiedan/gophermart/internal/logger"
)

type balanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := ctx.Value("user_id").(int)
	balance, err := h.balanceService.GetUserBalance(ctx, userId)
	if err != nil {
		logger.Log.Error("error getting user balance", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	withdrawnSum, err := h.balanceService.GetUserWithdrawalBalance(ctx, userId)
	if err != nil {
		logger.Log.Error("error getting user withdrawal balance", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := &balanceResponse{
		Current:   balance,
		Withdrawn: withdrawnSum,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
