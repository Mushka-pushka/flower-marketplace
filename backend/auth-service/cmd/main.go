// @title           Flower Marketplace Auth Service API
// @version         1.0
// @description     Сервис аутентификации и управления пользователями
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@flowermarketplace.com

// @host      localhost:8081
// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Введите токен в формате: Bearer <token>

package main

import (
	"context"
	"fmt"
	"log"
	"io/ioutil"  
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

// @title           Flower Marketplace Auth Service API
// @version         1.0
// @description     Сервис аутентификации и управления пользователями
// @host      localhost:8081
// @BasePath  /api/v1

func main() {
	cfg := config.Load()

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
	adminRepo := repository.NewAdminRepository(db)
	adminService := service.NewAdminService(adminRepo, cfg)
	adminHandler := handlers.NewAdminHandler(adminService)
	statsRepo := repository.NewAdminStatsRepository(db)
	statsService := service.NewAdminStatsService(statsRepo, cfg)
	statsHandler := handlers.NewAdminStatsHandler(statsService)

	// Настраиваем роутер
	// ----- ПУБЛИЧНЫЕ ЭНДПОИНТЫ (без авторизации) -----
	http.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	http.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	http.HandleFunc("POST /api/v1/auth/refresh", authHandler.RefreshToken)
	http.HandleFunc("POST /api/v1/auth/validate", authHandler.ValidateToken)

	// ----- ЗАЩИЩЕННЫЕ ЭНДПОИНТЫ (требуется JWT) -----
	http.HandleFunc("GET /api/v1/auth/me", authMiddleware.JWT(authHandler.Me))
	http.HandleFunc("PUT /api/v1/auth/profile", authMiddleware.JWT(authHandler.UpdateProfile))
	http.HandleFunc("PUT /api/v1/auth/password", authMiddleware.JWT(authHandler.ChangePassword))

	// ----- АДМИНИСТРИРОВАНИЕ (требуется JWT + роль admin) -----
	http.HandleFunc("GET /api/v1/admin/sellers", authMiddleware.JWT(adminHandler.GetSellers))
	http.HandleFunc("PUT /api/v1/admin/sellers/verify", authMiddleware.JWT(adminHandler.VerifySeller))
	http.HandleFunc("PUT /api/v1/admin/users/status", authMiddleware.JWT(adminHandler.UpdateUserStatus))
	http.HandleFunc("GET /api/v1/admin/users", authMiddleware.JWT(adminHandler.GetUsersList))
	http.HandleFunc("GET /api/v1/admin/users/list", authMiddleware.JWT(adminHandler.GetUsersListWithFilters))
	http.HandleFunc("GET /api/v1/admin/users/details", authMiddleware.JWT(adminHandler.GetUserByIDForAdmin))
	http.HandleFunc("GET /api/v1/admin/stats", authMiddleware.JWT(statsHandler.GetAdminStats))

	// ----- SWAGGER -----
	http.HandleFunc("GET /swagger/", func(w http.ResponseWriter, r *http.Request) {
		// Если запрос на doc.json
		if r.URL.Path == "/swagger/doc.json" {
			data, err := ioutil.ReadFile("./docs/swagger.json")
			if err != nil {
				http.Error(w, "Swagger docs not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
			return
		}

		// Если запрос на index.html или статику
		if r.URL.Path == "/swagger/" || r.URL.Path == "/swagger/index.html" {
			// Отдаём стандартный HTML-интерфейс Swagger UI
			html := `
			<!DOCTYPE html>
			<html>
			<head>
				<meta charset="UTF-8">
				<title>Swagger UI</title>
				<link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.10.5/swagger-ui.min.css" />
			</head>
			<body>
				<div id="swagger-ui"></div>
				<script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.10.5/swagger-ui-bundle.min.js"></script>
				<script>
					window.onload = function() {
						SwaggerUIBundle({
							url: "/swagger/doc.json",
							dom_id: '#swagger-ui',
							presets: [
								SwaggerUIBundle.presets.apis,
								SwaggerUIBundle.SwaggerUIStandalonePreset
							],
							layout: "BaseLayout"
						});
					};
				</script>
			</body>
			</html>
			`
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(html))
			return
		}

		http.NotFound(w, r)
	})

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