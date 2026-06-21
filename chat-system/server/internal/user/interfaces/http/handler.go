package http

import (
	"errors"
	"log"
	"net/http"
	"strings"

	userapplication "chat-system/server/internal/user/application"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	users *userapplication.Service
}

type registerRequest struct {
	Nickname        string `json:"nickname" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type loginRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewHandler(users *userapplication.Service) *Handler {
	return &Handler{users: users}
}

// Register 处理注册请求，校验表单后创建账号并返回新账号信息。
func (h *Handler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "passwords do not match"})
		return
	}
	user, err := h.users.Register(req.Nickname, req.Password)
	if err != nil {
		log.Printf("register failed for nickname %q: %v", req.Nickname, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"account": user.Account, "nickname": user.Nickname, "message": "register_success"})
}

// Login 校验登录凭证，生成会话令牌，并返回客户端建立连接所需的数据。
func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	session, err := h.users.Login(req.Account, req.Password)
	if err != nil {
		if errors.Is(err, userapplication.ErrInvalidCredentials) {
			log.Printf("login rejected for account %q: %v", req.Account, err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		log.Printf("login failed for account %q: %v", req.Account, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login service unavailable"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"type":            "login_success",
		"token":           session.Token,
		"account":         session.Account,
		"user_id":         session.UserID,
		"nickname":        session.Nickname,
		"last_login_time": session.LastLoginTime,
	})
}

// AuthMiddleware 校验 Bearer Token，并把用户身份写入 Gin 上下文，供后续接口使用。
func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		token := strings.TrimPrefix(auth, "Bearer ")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		session, err := h.users.ParseSession(token)
		if err != nil {
			log.Printf("auth failed: parse session: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set("user_id", session.UserID)
		c.Set("account", session.Account)
		c.Set("nickname", session.Nickname)
		c.Next()
	}
}
