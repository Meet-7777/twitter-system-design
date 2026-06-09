package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"twitter-system-design/internal/services"
)

type TimelineHandler struct {
	services *services.TimelineService
}

func NewTimelineHandler(
	s *services.TimelineService,
) *TimelineHandler {
	return &TimelineHandler{
		services: s,
	}
}

func (h *TimelineHandler) GetTimeline(
	w http.ResponseWriter,
	r *http.Request,
) {

	userIDStr := r.URL.Query().Get("user_id")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(
			w,
			"invalid user_id",
			http.StatusBadRequest,
		)
		return
	}

	tweets, err := h.services.GetTimeline(
		userID,
		20,
	)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		map[string]any{
			"tweets": tweets,
		},
	)
}
