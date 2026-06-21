package main

import (
	"log"

	"chat-system/server/api"
	"chat-system/server/configs"
	conversationapplication "chat-system/server/internal/conversation/application"
	conversationhttp "chat-system/server/internal/conversation/interfaces/http"
	conversationrepository "chat-system/server/internal/conversation/repository"
	friendapplication "chat-system/server/internal/friend/application"
	friendhttp "chat-system/server/internal/friend/interfaces/http"
	friendrepository "chat-system/server/internal/friend/repository"
	"chat-system/server/internal/infrastructure"
	"chat-system/server/internal/infrastructure/websocket"
	messageapplication "chat-system/server/internal/message/application"
	messagehttp "chat-system/server/internal/message/interfaces/http"
	messagerepository "chat-system/server/internal/message/repository"
	presenceapplication "chat-system/server/internal/presence/application"
	presenceinfra "chat-system/server/internal/presence/infra"
	userapplication "chat-system/server/internal/user/application"
	userinfra "chat-system/server/internal/user/infra"
	userhttp "chat-system/server/internal/user/interfaces/http"
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
	presenceStore := presenceinfra.NewRedisPresenceStore(infrastructure.RedisClient)

	userService := userapplication.NewService(userRepository, sessionStore)
	presenceService := presenceapplication.NewService(presenceStore)
	friendService := friendapplication.NewService(friendRepository, userService, presenceService)
	conversationService := conversationapplication.NewService(conversationRepository, userService, nil)
	messageService := messageapplication.NewService(messageRepository, userService, friendService, conversationService, presenceService)
	conversationService = conversationapplication.NewService(conversationRepository, userService, messageService)

	hub := websocket.NewHub(userService, messageService, presenceService)
	go hub.Run()

	handlers := api.Handlers{
		Users:         userhttp.NewHandler(userService),
		Friends:       friendhttp.NewHandler(friendService),
		Conversations: conversationhttp.NewHandler(conversationService),
		Messages:      messagehttp.NewHandler(messageService),
	}
	router := api.NewRouter(handlers, hub)
	log.Printf("chat server listening on %s", cfg.HTTPAddr)
	if err := router.Run(cfg.HTTPAddr); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
