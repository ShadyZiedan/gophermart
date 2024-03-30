package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/ShadyZiedan/gophermart/internal/logger"
	"github.com/ShadyZiedan/gophermart/internal/services"
)

func (h *Handler) UploadOrder(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		http.Error(w, "body empty", http.StatusBadRequest)
		return
	}
	orderNumber, err := strconv.Atoi(string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if userId, ok := r.Context().Value("user_id").(int); !ok {
		http.Error(w, "No user id found", http.StatusNotFound)
		return
	} else {
		_, err := h.orderService.CreateOrder(r.Context(), userId, orderNumber)
		if err != nil {
			if errors.Is(err, services.ErrOrderAlreadyExists) {
				w.WriteHeader(http.StatusOK)
				return
			}
			if errors.Is(err, services.ErrInvalidOrderNumber) {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			if errors.Is(err, services.ErrOrderBelongsToOtherUser) {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		go func() {
			err := h.accrualService.HandleOrder(context.Background(), strconv.Itoa(orderNumber))
			if err != nil {
				logger.Log.Error("error handling order by accrual service", zap.Error(err))
			}
		}()
	}

}
