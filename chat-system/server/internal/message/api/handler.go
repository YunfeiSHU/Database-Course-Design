package api

import (
	"log"
	"net/http"
	"strconv"

	messageapplication "chat-system/server/internal/message/application"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	messages *messageapplication.Service
}

func NewHandler(messages *messageapplication.Service) *Handler {
	return &Handler{messages: messages}
}

func (h *Handler) ListMessages(c *gin.Context) {
	h.History(c)
}

func (h *Handler) History(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	friendID64, err := strconv.ParseUint(c.Query("friend_id"), 10, 64)
	var messages interface{}
	if err != nil {
		account := c.Query("account")
		if account != "" {
			rows, findErr := h.messages.ListHistoryByAccount(userID, account, 50)
			if findErr != nil {
				log.Printf("list history by account failed for user %d account %q: %v", userID, account, findErr)
				c.JSON(http.StatusBadRequest, gin.H{"error": "friend_id or account is required"})
				return
			}
			messages = rows
		}
	}
	if err != nil && messages == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "friend_id or account is required"})
		return
	}
	if messages == nil {
		rows, err := h.messages.ListHistory(userID, uint(friendID64), 50)
		if err != nil {
			log.Printf("list history failed for user %d friend %d: %v", userID, friendID64, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		messages = rows
	}
	c.JSON(http.StatusOK, messages)
}
