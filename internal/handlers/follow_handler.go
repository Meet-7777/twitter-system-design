package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"twitter-system-design/internal/models"
	"twitter-system-design/internal/services"
)

type FollowHandler struct {
	services *services.FollowService
}

func NewFollowHandler(
	s *services.FollowService,
) *FollowHandler {
	return &FollowHandler{
		services: s,
	}
}

func (h *FollowHandler) FollowUser(
	w http.ResponseWriter,
	r *http.Request,
) {

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)

		json.NewEncoder(w).Encode(
			map[string]string{
				"error": "method not allowed",
			},
		)

		return
	}

	var follow models.Follow

	if err := json.NewDecoder(
		r.Body,
	).Decode(&follow); err != nil {

		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(
			map[string]string{
				"error": "invalid request body",
			},
		)

		return
	}

	if err := h.services.FollowUser(
		follow.FollowerID,
		follow.FolloweeID,
	); err != nil {

		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(
			map[string]string{
				"error": err.Error(),
			},
		)

		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(
		map[string]string{
			"message": "followed successfully",
		},
	)
}

func (h *FollowHandler) GetFollowing(
	w http.ResponseWriter,
	r *http.Request,
) {

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)

		json.NewEncoder(w).Encode(
			map[string]string{
				"error": "method not allowed",
			},
		)

		return
	}

	idStr := r.URL.Query().Get("user_id")

	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(
			map[string]string{
				"error": "user_id is required",
			},
		)

		return
	}

	id, err := strconv.Atoi(idStr)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(
			map[string]string{
				"error": "invalid user_id",
			},
		)

		return
	}

	users, err := h.services.GetFollowing(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(
			map[string]string{
				"error": err.Error(),
			},
		)

		return
	}

	type FollowingResponse struct {
		Following []models.User `json:"following"`
	}

	response := FollowingResponse{
		Following: users,
	}

	if users == nil {
		response.Following = []models.User{}
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}
