package main

import (
	"log"
	"net/http"
	"time"
	"url-shortener/internal/handler"
	"url-shortener/internal/repository"
	"url-shortener/internal/service"
	"url-shortener/pkg/config"
	"url-shortener/pkg/database"
)

func main() {
	// Load the Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to DB
	db, err := database.NewMySQLDB(&cfg.Database.MySQL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate DB
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 创建Repository实例
	repo := repository.NewMysqlShortURLRepository(db)

	// 连接Redis
	redisClient, err := database.NewRedisClient(&cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	svc := service.NewService(repo, redisClient)

	// 创建Handler实例（稍后会注入db依赖）
	handler := handler.NewHandler(svc)

	// 配置HTTP服务器
	server := &http.Server{
		Addr:         cfg.Server.Port,
		Handler:      handler,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	// 启动服务器
	log.Printf("Starting server on %s", server.Addr)
	log.Printf("Health:   http://localhost%s/health", server.Addr)
	log.Printf("Shorten:  http://localhost%s/shorten [POST]", server.Addr)
	log.Printf("Redirect: http://localhost%s/abc123 [GET]", server.Addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not start server: %v\n", err)
	}
}
