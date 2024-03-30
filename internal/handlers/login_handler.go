package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ShadyZiedan/gophermart/internal/services"
)

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest loginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if loginRequest.Login == "" || loginRequest.Password == "" {
		http.Error(w, "Login and Password are required", http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(r.Context(), loginRequest.Login, loginRequest.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", "Bearer "+token)
}
