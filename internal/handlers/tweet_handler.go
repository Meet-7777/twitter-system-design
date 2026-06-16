package handlers

import (
	"net/http"
	"twitter-system-design/internal/services"

	"github.com/gin-gonic/gin"
)

type TweetHandler struct {
	service *services.TweetService
}

func NewTweetHandler(s *services.TweetService) *TweetHandler {
	return &TweetHandler{
		service: s,
	}
}

func (h *TweetHandler) CreateTweet(c *gin.Context) {
	var body struct {
		UserID  int    `json:"user_id"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tweet, err := h.service.CreateTweet(body.UserID, body.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tweet)
}
