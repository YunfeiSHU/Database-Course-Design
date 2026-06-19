package api

import (
	"log"
	"net/http"

	conversationapplication "chat-system/server/internal/conversation/application"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	conversations *conversationapplication.Service
}

func NewHandler(conversations *conversationapplication.Service) *Handler {
	return &Handler{conversations: conversations}
}

func (h *Handler) ListConversations(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	items, err := h.conversations.List(userID)
	if err != nil {
		log.Printf("list conversations failed for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}
