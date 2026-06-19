package api

import (
	"net/http"

	"chat-system/server/api/websocket"

	"github.com/gin-gonic/gin"
)

func NewRouter(handlers Handlers, hub *websocket.Hub) *gin.Engine {
	router := gin.Default()
	router.Use(corsMiddleware())

	api := router.Group("/api")
	{
		api.POST("/register", handlers.Users.Register)
		api.POST("/login", handlers.Users.Login)
		authenticated := api.Group("", handlers.Users.AuthMiddleware())
		{
			authenticated.GET("/friends", handlers.Friends.ListFriends)
			authenticated.POST("/friends", handlers.Friends.AddFriend)
			authenticated.GET("/friend-requests", handlers.Friends.ListFriendRequests)
			authenticated.POST("/friend-requests/:id/accept", handlers.Friends.AcceptFriendRequest)
			authenticated.GET("/conversations", handlers.Conversations.ListConversations)
			authenticated.GET("/messages", handlers.Messages.ListMessages)
			authenticated.GET("/history", handlers.Messages.History)
		}
	}

	router.GET("/ws", websocket.Handler(hub))
	return router
}

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
