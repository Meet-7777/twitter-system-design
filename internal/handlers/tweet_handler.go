package handlers

import (
	"encoding/json"
	"net/http"

	"twitter-system-design/internal/models"
	"twitter-system-design/internal/services"
)

type TweetHandler struct {
	service *services.TweetService
}

func NewTweetHandler(
	s *services.TweetService,
) *TweetHandler {

	return &TweetHandler{
		service: s,
	}
}

func (h *TweetHandler) Tweets(
	w http.ResponseWriter,
	r *http.Request,
) {

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	switch r.Method {

	case http.MethodPost:

		h.createTweet(w, r)

	default:

		w.WriteHeader(
			http.StatusMethodNotAllowed,
		)

		json.NewEncoder(w).Encode(
			map[string]string{
				"error": "method not allowed",
			},
		)

	}

}

func (h *TweetHandler) createTweet(
	w http.ResponseWriter,
	r *http.Request,
) {

	var tweet models.Tweet

	err := json.NewDecoder(
		r.Body,
	).Decode(&tweet)

	if err != nil {

		w.WriteHeader(
			http.StatusBadRequest,
		)

		json.NewEncoder(w).Encode(
			map[string]string{
				"error": "invalid request body",
			},
		)

		return
	}

	id, err := h.service.CreateTweet(
		&tweet,
	)

	if err != nil {

		w.WriteHeader(
			http.StatusBadRequest,
		)

		json.NewEncoder(w).Encode(
			map[string]string{
				"error": err.Error(),
			},
		)

		return
	}

	w.WriteHeader(
		http.StatusCreated,
	)

	json.NewEncoder(w).Encode(
		map[string]any{
			"message":  "tweet created",
			"tweet_id": id,
		},
	)

}
