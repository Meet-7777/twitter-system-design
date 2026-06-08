package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"twitter-system-design/internal/models"
	"twitter-system-design/internal/services"
)

type TweetHandler struct {
	services *services.TweetService
}

func NewTweetHandler(s *services.TweetService) *TweetHandler {
	return &TweetHandler{services: s}
}

func (h *TweetHandler) CreateTweet(w http.ResponseWriter, r *http.Request) {
	var tweet models.Tweet
	err := json.NewDecoder(r.Body).Decode(&tweet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.services.CreateTweet(&tweet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Tweet created successfully"))
}

func (h *TweetHandler) GetTweetsByUserID(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(user)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}
	if userID <= 0 {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	tweets, err := h.services.GetTweetsByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tweets)
}
