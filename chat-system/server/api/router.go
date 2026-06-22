package api

import (
	"net/http"

	websocket "chat-system/server/internal/infrastructure/websocket"

	"github.com/gin-gonic/gin"
)

// NewRouter 构建 HTTP 路由，统一挂载 REST API、鉴权中间件和 WebSocket 入口。
func NewRouter(handlers Handlers, hub *websocket.Hub) *gin.Engine {
	router := gin.Default()
	router.Use(corsMiddleware())

	api := router.Group("/api")
	{
		api.POST("/register", handlers.Users.Register)
		api.POST("/login", handlers.Users.Login)
		authenticated := api.Group("", handlers.Users.AuthMiddleware())
		{
			authenticated.GET("/friends", handlers.Friends.GetFriendList)
			authenticated.POST("/friends", handlers.Friends.RequestFriend)
			authenticated.GET("/friend-requests", handlers.Friends.GetFriendRequestList)
			authenticated.POST("/friend-requests/:id/accept", handlers.Friends.ApproveFriendRequest)
			authenticated.GET("/conversations", handlers.Conversations.GetConversationList)
			authenticated.GET("/messages", handlers.Messages.GetHistoryMessages)
			authenticated.GET("/history", handlers.Messages.GetConversationHistory)
			authenticated.POST("/messages/recall", handlers.Messages.RecallMessage)
		}
	}

	router.GET("/ws", websocket.Handler(hub))
	return router
}

// corsMiddleware 处理跨域请求，并允许前端携带认证信息访问接口。
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
