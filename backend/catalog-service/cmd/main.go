// @title           Flower Marketplace Catalog Service API
// @version         1.0
// @description     Сервис управления товарами, категориями и поиском
// @host      localhost:8082
// @BasePath  /api/v1

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
	"github.com/Mushka-pushka/flower-marketplace/backend/catalog-service/internal/middleware"
	"github.com/Mushka-pushka/flower-marketplace/backend/catalog-service/internal/repository"
	"github.com/Mushka-pushka/flower-marketplace/backend/catalog-service/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"

	_ "github.com/Mushka-pushka/flower-marketplace/backend/catalog-service/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	
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

	// ИНИЦИАЛИЗАЦИЯ РЕПОЗИТОРИЕВ
	productRepo := repository.NewProductRepository(db)
	cartRepo := repository.NewCartRepository(db)
	favoriteRepo := repository.NewFavoriteRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	autocompleteRepo := repository.NewAutocompleteRepository(db)
	addressRepo := repository.NewAddressRepository(db)
	categoryAdminRepo := repository.NewCategoryAdminRepository(db)

	// ИНИЦИАЛИЗАЦИЯ СЕРВИСА
	catalogService := service.NewCatalogService(productRepo, cartRepo, favoriteRepo, reviewRepo, autocompleteRepo, addressRepo, categoryAdminRepo, cfg, valkeyClient)
	catalogHandler := handlers.NewCatalogHandler(catalogService)

	// РОУТЫ

	// ----- ТОВАРЫ (CRUD) - ПУБЛИЧНЫЕ -----
	http.HandleFunc("GET /api/v1/catalog/products", catalogHandler.GetProductByID)
	http.HandleFunc("GET /api/v1/catalog/products/{id}", catalogHandler.GetProductByIDPath)
	http.HandleFunc("GET /api/v1/catalog/products/slug/{slug}", catalogHandler.GetProductBySlug)

	// ----- ТОВАРЫ (CRUD) - ЗАЩИЩЕННЫЕ -----
	http.HandleFunc("POST /api/v1/catalog/products", middleware.AuthMiddleware(catalogHandler.CreateProduct))
	http.HandleFunc("PUT /api/v1/catalog/products/{id}", middleware.AuthMiddleware(catalogHandler.UpdateProduct))
	http.HandleFunc("DELETE /api/v1/catalog/products/{id}", middleware.AuthMiddleware(catalogHandler.DeleteProduct))
	http.HandleFunc("POST /api/v1/catalog/products/decrease-stock", catalogHandler.DecreaseStock)

	// ----- ПОИСК И КАТЕГОРИИ (ПУБЛИЧНЫЕ) -----
	http.HandleFunc("GET /api/v1/catalog/search", catalogHandler.SearchProducts)
	http.HandleFunc("GET /api/v1/catalog/categories", catalogHandler.GetCategories)

	// ----- КОРЗИНА (CART) - ЗАЩИЩЕННЫЕ -----
	http.HandleFunc("POST /api/v1/catalog/cart", middleware.AuthMiddleware(catalogHandler.AddToCart))
	http.HandleFunc("GET /api/v1/catalog/cart", middleware.AuthMiddleware(catalogHandler.GetCart))
	http.HandleFunc("PUT /api/v1/catalog/cart", middleware.AuthMiddleware(catalogHandler.UpdateCartItem))
	http.HandleFunc("DELETE /api/v1/catalog/cart", middleware.AuthMiddleware(catalogHandler.RemoveFromCart))

	// ----- ИЗБРАННОЕ (FAVORITES) - ЗАЩИЩЕННЫЕ -----
	http.HandleFunc("POST /api/v1/catalog/favorites", middleware.AuthMiddleware(catalogHandler.AddFavorite))
	http.HandleFunc("GET /api/v1/catalog/favorites", middleware.AuthMiddleware(catalogHandler.GetFavorites))
	http.HandleFunc("DELETE /api/v1/catalog/favorites", middleware.AuthMiddleware(catalogHandler.RemoveFavorite))
	http.HandleFunc("GET /api/v1/catalog/favorites/check", middleware.AuthMiddleware(catalogHandler.CheckFavorite))

	// ----- ОТЗЫВЫ (REVIEWS) - ЗАЩИЩЕННЫЕ -----
	http.HandleFunc("POST /api/v1/catalog/reviews", middleware.AuthMiddleware(catalogHandler.CreateReview))
	http.HandleFunc("GET /api/v1/catalog/reviews/me", middleware.AuthMiddleware(catalogHandler.GetMyReviews))
	http.HandleFunc("PUT /api/v1/catalog/reviews", middleware.AuthMiddleware(catalogHandler.UpdateReview))
	http.HandleFunc("DELETE /api/v1/catalog/reviews", middleware.AuthMiddleware(catalogHandler.DeleteReview))

	// ----- ОТЗЫВЫ (REVIEWS) - ПУБЛИЧНЫЕ -----
	http.HandleFunc("GET /api/v1/catalog/reviews", catalogHandler.GetProductReviews)

	// ----- АВТОДОПОЛНЕНИЕ (ПУБЛИЧНОЕ) -----
	http.HandleFunc("GET /api/v1/catalog/autocomplete", catalogHandler.GetAutocompleteSuggestions)

	// ----- АДРЕСА ДОСТАВКИ - ЗАЩИЩЕННЫЕ -----
	http.HandleFunc("POST /api/v1/catalog/addresses", middleware.AuthMiddleware(catalogHandler.CreateAddress))
	http.HandleFunc("GET /api/v1/catalog/addresses", middleware.AuthMiddleware(catalogHandler.GetAddresses))
	http.HandleFunc("PUT /api/v1/catalog/addresses", middleware.AuthMiddleware(catalogHandler.UpdateAddress))
	http.HandleFunc("DELETE /api/v1/catalog/addresses", middleware.AuthMiddleware(catalogHandler.DeleteAddress))
	http.HandleFunc("POST /api/v1/catalog/addresses/default", middleware.AuthMiddleware(catalogHandler.SetDefaultAddress))

	// ----- АДМИН: КАТЕГОРИИ - ЗАЩИЩЕННЫЕ -----
	http.HandleFunc("POST /api/v1/admin/categories", middleware.AuthMiddleware(catalogHandler.AdminCreateCategory))
	http.HandleFunc("GET /api/v1/admin/categories", middleware.AuthMiddleware(catalogHandler.AdminGetAllCategories))
	http.HandleFunc("GET /api/v1/admin/categories/id", middleware.AuthMiddleware(catalogHandler.AdminGetCategoryByID))
	http.HandleFunc("PUT /api/v1/admin/categories", middleware.AuthMiddleware(catalogHandler.AdminUpdateCategory))
	http.HandleFunc("DELETE /api/v1/admin/categories", middleware.AuthMiddleware(catalogHandler.AdminDeleteCategory))

	// ---- SWAGGER ----
	http.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
	http.HandleFunc("GET /swagger/*", httpSwagger.WrapHandler)

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