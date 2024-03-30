package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/ShadyZiedan/gophermart/internal/logger"
	"github.com/ShadyZiedan/gophermart/internal/security"
)

type balanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(security.UserIDKey{}).(int)
	balance, err := h.balanceService.GetUserBalance(ctx, userID)
	if err != nil {
		logger.Log.Error("error getting user balance", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	withdrawnSum, err := h.balanceService.GetUserWithdrawalBalance(ctx, userID)
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
