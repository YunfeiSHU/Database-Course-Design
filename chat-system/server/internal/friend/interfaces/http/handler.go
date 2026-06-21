package http

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	friendApplication "chat-system/server/internal/friend/application"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	friendService *friendApplication.Service
}

type addFriendRequest struct {
	Account string `json:"account" binding:"required"`
}

func NewHandler(friendService *friendApplication.Service) *Handler {
	return &Handler{friendService: friendService}
}

// GetFriendList 返回当前登录用户的已接受好友列表，并包含在线状态。
func (h *Handler) GetFriendList(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	friends, err := h.friendService.ListFriends(userID)
	if err != nil {
		log.Printf("list friends failed for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, friends)
}

// RequestFriend 由当前登录用户向目标账号发起好友请求。
func (h *Handler) RequestFriend(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var req addFriendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.friendService.AddFriend(userID, req.Account); err != nil {
		if errors.Is(err, friendApplication.ErrCannotAddSelf) {
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

// GetFriendRequestList 返回当前登录用户需要处理的待接受好友请求。
func (h *Handler) GetFriendRequestList(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	requests, err := h.friendService.ListFriendRequests(userID)
	if err != nil {
		log.Printf("list friend requests failed for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, requests)
}

// ApproveFriendRequest 根据路由中的请求 ID 同意一条待处理好友请求。
func (h *Handler) ApproveFriendRequest(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	requestID64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}
	if err := h.friendService.AcceptFriendRequest(uint(requestID64), userID); err != nil {
		log.Printf("accept friend request failed for user %d request %d: %v", userID, requestID64, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
