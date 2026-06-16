package handlers

import (
	"net/http"
	"strconv"
	"twitter-system-design/internal/services"

	"github.com/gin-gonic/gin"
)

// handlers/timeline_handler.go
type TimelineHandler struct {
	service *services.TimelineService
}

func NewTimelineHandler(s *services.TimelineService) *TimelineHandler {
	return &TimelineHandler{service: s}
}

func (h *TimelineHandler) GetFeed(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	feed, err := h.service.GetFeed(id, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, feed)
}
