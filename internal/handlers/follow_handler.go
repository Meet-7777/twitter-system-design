package handlers

import (
	"net/http"
	"twitter-system-design/internal/services"

	"github.com/gin-gonic/gin"
)

type FollowHandler struct {
	service *services.FollowService
}

func NewFollowHandler(s *services.FollowService) *FollowHandler {
	return &FollowHandler{service: s}
}

func (h *FollowHandler) FollowUser(c *gin.Context) {
	var body struct {
		FollowerID int `json:"follower_id"`
		FolloweeID int `json:"followee_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.FollowUser(body.FollowerID, body.FolloweeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "followed successfully"})
}
