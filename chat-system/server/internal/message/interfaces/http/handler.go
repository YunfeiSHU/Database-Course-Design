package http

import (
	"log"
	"net/http"
	"strconv"

	messageApplication "chat-system/server/internal/message/application"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	messageService *messageApplication.Service
}

type recallRequest struct {
	MessageID uint `json:"message_id"`
}

func NewHandler(messageService *messageApplication.Service) *Handler {
	return &Handler{messageService: messageService}
}

// GetHistoryMessages 是历史消息查询的直接入口，供前端打开会话时拉取记录。
func (h *Handler) GetHistoryMessages(c *gin.Context) {
	h.GetConversationHistory(c)
}

// GetConversationHistory 根据好友 ID 或账号查询最近聊天记录，供前端会话区域回显历史消息。
func (h *Handler) GetConversationHistory(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	friendID64, err := strconv.ParseUint(c.Query("friend_id"), 10, 64)
	var messages interface{}
	if err != nil {
		account := c.Query("account")
		if account != "" {
			rows, findErr := h.messageService.ListHistoryByAccount(userID, account, 50)
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
		rows, err := h.messageService.ListHistory(userID, uint(friendID64), 50)
		if err != nil {
			log.Printf("list history failed for user %d friend %d: %v", userID, friendID64, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		messages = rows
	}
	c.JSON(http.StatusOK, messages)
}

// RecallMessage 撤回指定消息，仅允许消息发送者执行。
func (h *Handler) RecallMessage(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var req recallRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.MessageID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "message_id is required"})
		return
	}
	account := c.MustGet("account").(string)
	message, err := h.messageService.Recall(userID, account, req.MessageID)
	if err != nil {
		status := http.StatusInternalServerError
		if err == messageApplication.ErrCannotRecall {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, message)
}
