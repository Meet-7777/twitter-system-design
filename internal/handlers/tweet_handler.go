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

// TweetsResponse is the structure for the get tweets response
type TweetsResponse struct {
	Tweets []models.Tweet `json:"tweets"`
}

// Tweets handles routing for tweet-related requests.
func (h *TweetHandler) Tweets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		h.getTweets(w, r)
	case http.MethodPost:
		h.createTweet(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
	}
}

func (h *TweetHandler) createTweet(w http.ResponseWriter, r *http.Request) {
	var tweet models.Tweet
	if err := json.NewDecoder(r.Body).Decode(&tweet); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	tweetID, err := h.services.CreateTweet(&tweet)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(
			map[string]string{
				"error": err.Error(),
			},
		)

		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"message": "Tweet created successfully", "tweetID": tweetID})
}

func (h *TweetHandler) getTweets(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "user_id is required"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid user_id"})
		return
	}

	tweets, err := h.services.GetTweetsByUserID(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Could not retrieve tweets"})
		return
	}

	response := TweetsResponse{Tweets: tweets}
	if tweets == nil {
		response.Tweets = []models.Tweet{}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
