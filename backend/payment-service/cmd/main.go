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

	"github.com/Mushka-pushka/flower-marketplace/backend/payment-service/internal/config"
	"github.com/Mushka-pushka/flower-marketplace/backend/payment-service/internal/handlers"
	"github.com/Mushka-pushka/flower-marketplace/backend/payment-service/internal/middleware"
	"github.com/Mushka-pushka/flower-marketplace/backend/payment-service/internal/repository"
	"github.com/Mushka-pushka/flower-marketplace/backend/payment-service/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	cfg := config.Load()

	log.Printf("Payment Service starting on port %s", cfg.Port)

	// Подключение к PostgreSQL
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

	// Подключение к RabbitMQ
	rabbitConn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()
	log.Println("Connected to RabbitMQ")

	rabbitCh, err := rabbitConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open RabbitMQ channel: %v", err)
	}
	defer rabbitCh.Close()

	// Объявляем очереди
	queues := []string{
		"payment.status_changed",
		"payment.created",
		"payment.confirmed",
		"payment.failed",
	}
	for _, queue := range queues {
		_, err = rabbitCh.QueueDeclare(
			queue,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Fatalf("Failed to declare queue %s: %v", queue, err)
		}
		log.Printf("RabbitMQ queue declared: %s", queue)
	}

	// Инициализация
	paymentRepo := repository.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepo, cfg, rabbitCh)
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg)

	// ============================================================
	// РОУТЫ
	// ============================================================

	// ----- ЗАЩИЩЕННЫЕ ЭНДПОИНТЫ (требуется JWT) -----
	// Создание платежа - только для авторизованных пользователей
	http.HandleFunc("POST /api/v1/payments", authMiddleware.AuthMiddleware(paymentHandler.CreatePayment))
	// Получение статуса платежа - только для авторизованных пользователей
	http.HandleFunc("GET /api/v1/payments", authMiddleware.AuthMiddleware(paymentHandler.GetPaymentStatus))
	// Получение платежа по ID заказа - только для авторизованных пользователей
	http.HandleFunc("GET /api/v1/payments/order", authMiddleware.AuthMiddleware(paymentHandler.GetPaymentByOrderID))

	// ----- SWAGGER -----
	http.HandleFunc("GET /swagger/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/swagger/doc.json" {
			data, err := os.ReadFile("./docs/swagger.json")
			if err != nil {
				http.Error(w, "Swagger docs not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
			return
		}

		html := `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Payment Service API</title>
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
						presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.SwaggerUIStandalonePreset],
						layout: "BaseLayout"
					});
				};
			</script>
		</body>
		</html>
		`
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	})

	// ============================================================
	// СЕРВЕР
	// ============================================================
	
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

	// Graceful shutdown
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