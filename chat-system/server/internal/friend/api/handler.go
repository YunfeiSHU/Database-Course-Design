package api

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	friendapplication "chat-system/server/internal/friend/application"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	friends *friendapplication.Service
}

type addFriendRequest struct {
	Account string `json:"account" binding:"required"`
}

func NewHandler(friends *friendapplication.Service) *Handler {
	return &Handler{friends: friends}
}

func (h *Handler) ListFriends(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	friends, err := h.friends.ListFriends(userID)
	if err != nil {
		log.Printf("list friends failed for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, friends)
}

func (h *Handler) AddFriend(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var req addFriendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.friends.AddFriend(userID, req.Account); err != nil {
		if errors.Is(err, friendapplication.ErrCannotAddSelf) {
			log.Printf("add friend rejected for user %d account %q: %v", userID, req.Account, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("add friend failed for user %d account %q: %v", userID, req.Account, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "friend request sent"})
}

func (h *Handler) ListFriendRequests(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	requests, err := h.friends.ListFriendRequests(userID)
	if err != nil {
		log.Printf("list friend requests failed for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, requests)
}

func (h *Handler) AcceptFriendRequest(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	requestID64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}
	if err := h.friends.AcceptFriendRequest(uint(requestID64), userID); err != nil {
		log.Printf("accept friend request failed for user %d request %d: %v", userID, requestID64, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
