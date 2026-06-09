package handlers

import (
	"encoding/json"
	"net/http"

	"twitter-system-design/internal/models"
	"twitter-system-design/internal/services"
)

type UserHandler struct {
	services *services.UserService
}

func NewUserHandler(s *services.UserService) *UserHandler {
	return &UserHandler{services: s}
}

func (h *UserHandler) CreateUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)

		json.NewEncoder(w).Encode(
			map[string]string{
				"error": "method not allowed",
			},
		)

		return
	}

	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(
			map[string]string{
				"error": "invalid request body",
			},
		)

		return
	}

	if err := h.services.CreateUser(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(
			map[string]string{
				"error": err.Error(),
			},
		)

		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(user)
}
