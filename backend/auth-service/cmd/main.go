package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Mushka-pushka/flower-marketplace/backend/auth-service/internal/config"
	"github.com/Mushka-pushka/flower-marketplace/backend/auth-service/internal/handlers"
	"github.com/Mushka-pushka/flower-marketplace/backend/auth-service/internal/middleware"
	"github.com/Mushka-pushka/flower-marketplace/backend/auth-service/internal/repository"
	"github.com/Mushka-pushka/flower-marketplace/backend/auth-service/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	cfg := config.Load()

	log.Printf("🔍 DB_USER: %s", cfg.DBUser)
	log.Printf("🔍 DB_PASSWORD: %s", cfg.DBPassword)
	log.Printf("🔍 DB_HOST: %s", cfg.DBHost)
	log.Printf("🔍 DB_PORT: %s", cfg.DBPort)
	log.Printf("🔍 DB_NAME: %s", cfg.DBName)

	log.Printf("Auth Service starting on port %s", cfg.Port)

	// Подключаемся к PostgreSQL
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	log.Printf("Connection string: %s", connString)

		db, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	// Инициализируем репозитории, сервисы, хендлеры
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg)
	authHandler := handlers.NewAuthHandler(authService)
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Настраиваем роутер
	http.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	http.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	http.HandleFunc("GET /api/v1/auth/me", authMiddleware.JWT(authHandler.Me))

	// Создаём сервер
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      http.DefaultServeMux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запускаем сервер в горутине
	go func() {
		log.Printf("Server is running on http://localhost:%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Ждём сигнала для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}