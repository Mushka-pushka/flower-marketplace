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

	"github.com/Mushka-pushka/flower-marketplace/backend/catalog-service/internal/config"
	"github.com/Mushka-pushka/flower-marketplace/backend/catalog-service/internal/handlers"
	"github.com/Mushka-pushka/flower-marketplace/backend/catalog-service/internal/repository"
	"github.com/Mushka-pushka/flower-marketplace/backend/catalog-service/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	log.Printf("Catalog Service starting on port %s", cfg.Port)

	// ПОДКЛЮЧЕНИЕ К POSTGRESQL
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	// ПОДКЛЮЧЕНИЕ К VALKEY (REDIS)
	valkeyClient := redis.NewClient(&redis.Options{
		Addr:     cfg.ValkeyHost + ":" + cfg.ValkeyPort,
		Password: "", // пароль не установлен
		DB:       cfg.ValkeyDB,
	})

	if err := valkeyClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Valkey: %v", err)
	}
	log.Printf("Connected to Valkey at %s:%s", cfg.ValkeyHost, cfg.ValkeyPort)

	// ИНИЦИАЛИЗАЦИЯ
	productRepo := repository.NewProductRepository(db)
	catalogService := service.NewCatalogService(productRepo, cfg, valkeyClient)
	catalogHandler := handlers.NewCatalogHandler(catalogService)

	// РОУТЫ
	http.HandleFunc("POST /api/v1/catalog/products", catalogHandler.CreateProduct)
	http.HandleFunc("GET /api/v1/catalog/products", catalogHandler.GetProductByID)
	http.HandleFunc("GET /api/v1/catalog/products/slug/{slug}", catalogHandler.GetProductBySlug)
	http.HandleFunc("GET /api/v1/catalog/search", catalogHandler.SearchProducts)
	http.HandleFunc("GET /api/v1/catalog/categories", catalogHandler.GetCategories)
	http.HandleFunc("PUT /api/v1/catalog/products/{id}", catalogHandler.UpdateProduct)
	http.HandleFunc("DELETE /api/v1/catalog/products/{id}", catalogHandler.DeleteProduct)

	// СЕРВЕР
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      http.DefaultServeMux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server is running on http://localhost:%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// GRACEFUL SHUTDOWN
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