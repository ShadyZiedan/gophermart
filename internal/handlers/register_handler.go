package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/ShadyZiedan/gophermart/internal/logger"
	"github.com/ShadyZiedan/gophermart/internal/services"
)

type registrationRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

var (
	ErrUserNameIsRequired = errors.New("username is required")
	ErrPasswordIsRequired = errors.New("password is required")
)

func (r *registrationRequest) validate() (bool, []error) {
	errs := make([]error, 0, 2)
	if r.Login == "" {
		errs = append(errs, ErrUserNameIsRequired)
	}
	if r.Password == "" {
		errs = append(errs, ErrPasswordIsRequired)
	}
	if len(errs) > 0 {
		return false, errs
	}
	return true, errs
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var registrationRequest registrationRequest
	err := json.NewDecoder(r.Body).Decode(&registrationRequest)
	if err != nil {
		logger.Log.Info("Failed to decode registration request", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if isValid, errs := registrationRequest.validate(); !isValid {
		err := errors.Join(errs...)
		logger.Log.Info("register validation error", zap.Error(err), zap.Any("registrationRequest", registrationRequest))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.authService.Register(ctx, registrationRequest.Login, registrationRequest.Password)
	if err != nil {
		if errors.Is(err, services.ErrUsernameAlreadyTaken) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token, err := h.authService.Login(ctx, registrationRequest.Login, registrationRequest.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w.WriteHeader(http.StatusOK)
}
