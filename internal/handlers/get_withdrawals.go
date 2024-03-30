package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ShadyZiedan/gophermart/internal/security"
)

type withdrawalsResponseModel struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

func (h *Handler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(security.UserIDKey{}).(int)
	withdrawals, err := h.balanceService.GetWithdrawals(ctx, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var response []*withdrawalsResponseModel
	for _, withdrawal := range withdrawals {
		model := &withdrawalsResponseModel{
			Order:       strconv.Itoa(withdrawal.OrderNumber),
			Sum:         withdrawal.Sum,
			ProcessedAt: withdrawal.ProcessedAt.Format(time.RFC3339),
		}
		response = append(response, model)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
