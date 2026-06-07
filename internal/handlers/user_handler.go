package handlers

import (
	"encoding/json"
	"net/http"
	"twitter-system-design/internal/models"
	"twitter-system-design/internal/services"
)

type UserHandler struct {
	Service *services.UserService
}

func NewUserHandler(s *services.UserService) *UserHandler {
	return &UserHandler{Service: s}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	var user models.User

	json.NewDecoder(r.Body).Decode(&user)

	err := h.Service.CreateUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}
