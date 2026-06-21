package http

import (
	"log"
	"net/http"

	conversationApplication "chat-system/server/internal/conversation/application"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	conversationService *conversationApplication.Service
}

func NewHandler(conversationService *conversationApplication.Service) *Handler {
	return &Handler{conversationService: conversationService}
}

// GetConversationList 返回当前登录用户的会话列表，供客户端侧边栏展示。
func (h *Handler) GetConversationList(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	items, err := h.conversationService.List(userID)
	if err != nil {
		log.Printf("list conversations failed for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}
