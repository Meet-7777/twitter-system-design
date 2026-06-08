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

func NewFollowHandler(s *services.FollowService) *FollowHandler {
	return &FollowHandler{services: s}
}

func (h *FollowHandler) FollowUser(
	w http.ResponseWriter, r *http.Request,
) {
	var follow models.Follow
	err := json.NewDecoder(r.Body).Decode(&follow)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.services.FollowUser(
		follow.FollowerID,
		follow.FolloweeID,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("followed successfully"))
}

func (h *FollowHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	user, err := h.services.GetFollowing(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
