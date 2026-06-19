package main

import (
	"log"

	"chat-system/server/api"
	"chat-system/server/api/websocket"
	"chat-system/server/configs"
	conversationapi "chat-system/server/internal/conversation/api"
	conversationapplication "chat-system/server/internal/conversation/application"
	conversationrepository "chat-system/server/internal/conversation/repository"
	friendapi "chat-system/server/internal/friend/api"
	friendapplication "chat-system/server/internal/friend/application"
	friendrepository "chat-system/server/internal/friend/repository"
	"chat-system/server/internal/infrastructure"
	messageapi "chat-system/server/internal/message/api"
	messageapplication "chat-system/server/internal/message/application"
	messagerepository "chat-system/server/internal/message/repository"
	notificationapplication "chat-system/server/internal/notification/application"
	notificationinfra "chat-system/server/internal/notification/infra"
	userapi "chat-system/server/internal/user/api"
	userapplication "chat-system/server/internal/user/application"
	userinfra "chat-system/server/internal/user/infra"
	userrepository "chat-system/server/internal/user/repository"
)

func main() {
	cfg := configs.Load()

	if err := infrastructure.InitMySQL(cfg.MySQLDSN); err != nil {
		log.Fatalf("init mysql failed: %v", err)
	}
	infrastructure.InitRedis(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)

	userRepository := userrepository.NewMySQLRepository()
	sessionStore := userinfra.NewRedisSessionStore(infrastructure.RedisClient)
	friendRepository := friendrepository.NewMySQLRepository()
	conversationRepository := conversationrepository.NewMySQLRepository()
	messageRepository := messagerepository.NewMySQLRepository()
	presenceStore := notificationinfra.NewRedisPresenceStore(infrastructure.RedisClient)

	userService := userapplication.NewService(userRepository, sessionStore)
	notificationService := notificationapplication.NewService(presenceStore)
	friendService := friendapplication.NewService(friendRepository, userService, notificationService)
	conversationService := conversationapplication.NewService(conversationRepository, userService, nil)
	messageService := messageapplication.NewService(messageRepository, userService, friendService, conversationService, notificationService)
	conversationService = conversationapplication.NewService(conversationRepository, userService, messageService)

	hub := websocket.NewHub(userService, messageService, notificationService)
	go hub.Run()

	handlers := api.Handlers{
		Users:         userapi.NewHandler(userService),
		Friends:       friendapi.NewHandler(friendService),
		Conversations: conversationapi.NewHandler(conversationService),
		Messages:      messageapi.NewHandler(messageService),
	}
	router := api.NewRouter(handlers, hub)
	log.Printf("chat server listening on %s", cfg.HTTPAddr)
	if err := router.Run(cfg.HTTPAddr); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
