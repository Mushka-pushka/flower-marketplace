// @title           Flower Marketplace Order Service API
// @version         1.0
// @description     Сервис управления заказами, доставкой и статусами
// @host      localhost:8083
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

	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/config"
	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/handlers"
	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/middleware"
	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/repository"
	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/service"
	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/worker"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	amqp "github.com/rabbitmq/amqp091-go"

	_ "github.com/Mushka-pushka/flower-marketplace/backend/order-service/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	
	cfg := config.Load()

	log.Printf("Order Service starting on port %s", cfg.Port)

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
		"order.created",
		"order.cancelled",
		"order.status_changed",
		"order.payment_completed",
		"notification.order_created",
		"notification.order_status_changed",
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

	// ИНИЦИАЛИЗАЦИЯ

	// Репозитории
	orderRepo := repository.NewOrderRepository(db)
	analyticsRepo := repository.NewAnalyticsRepository(db, cfg)

	// Сервисы
	analyticsService := service.NewAnalyticsService(analyticsRepo, cfg)
	orderService := service.NewOrderService(orderRepo, cfg, rabbitCh)
	notificationRepo := repository.NewNotificationRepository(db)
    notificationService := service.NewNotificationService(notificationRepo, cfg, rabbitCh)

	// Хендлеры
	orderHandler := handlers.NewOrderHandler(orderService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)
	http.HandleFunc("GET /api/v1/orders/items", orderHandler.GetOrderItemsByCustomer)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg)

	// Воркеры
	orderWorker := worker.NewOrderWorker(orderService, rabbitCh)
	notificationWorker := worker.NewNotificationWorker(notificationService, rabbitCh)

	// ============================================================
	// ЗАПУСК ВОРКЕРОВ
	// ============================================================
	
	// Запускаем несколько воркеров для обработки заказов
	numWorkers := 3 // Количество воркеров
	
	for i := 0; i < numWorkers; i++ {
		workerID := i
		go func() {
			log.Printf("Starting order worker #%d", workerID)
			if err := orderWorker.Start(context.Background()); err != nil {
				log.Printf("Order worker #%d error: %v", workerID, err)
			}
		}()
	}

	// Запускаем notification worker
	go func() {
		log.Println("Starting notification worker")
		if err := notificationWorker.Start(context.Background()); err != nil {
			log.Printf("Notification worker error: %v", err)
		}
	}()

	// ============================================================
	// РОУТЫ
	// ============================================================

	// ----- ПУБЛИЧНЫЕ ЭНДПОИНТЫ (без авторизации) -----
	// GetOrder - публичный, так как может использоваться для отслеживания статуса
	http.HandleFunc("GET /api/v1/orders", orderHandler.GetOrder)
	// GetOrdersByShop - публичный, но внутри проверяется принадлежность магазина
	http.HandleFunc("GET /api/v1/orders/shop", orderHandler.GetOrdersByShop)

	// ----- ЗАЩИЩЕННЫЕ ЭНДПОИНТЫ -----
	// Создание заказа - только для авторизованных пользователей
	http.HandleFunc("POST /api/v1/orders", authMiddleware.AuthMiddleware(orderHandler.CreateOrder))
	// Получение заказов покупателя - только для авторизованных пользователей
	http.HandleFunc("GET /api/v1/orders/customer", authMiddleware.AuthMiddleware(orderHandler.GetOrdersByCustomer))
	// Отмена заказа - только для авторизованных пользователей
	http.HandleFunc("POST /api/v1/orders/cancel", authMiddleware.AuthMiddleware(orderHandler.CancelOrder))
	// Обновление статуса заказа продавцом - только для авторизованных пользователей с ролью seller
	http.HandleFunc("PUT /api/v1/orders/status", authMiddleware.AuthMiddleware(orderHandler.UpdateOrderStatusBySeller))
	// Проверка возможности оставить отзыв - только для авторизованных пользователей
	http.HandleFunc("GET /api/v1/orders/can-review", authMiddleware.AuthMiddleware(orderHandler.CanReview))

	// ----- АНАЛИТИКА -----
	http.HandleFunc("GET /api/v1/analytics/seller", authMiddleware.AuthMiddleware(analyticsHandler.GetSellerAnalytics))
	http.HandleFunc("GET /api/v1/analytics/popular", authMiddleware.AuthMiddleware(analyticsHandler.GetPopularProducts))
	http.HandleFunc("GET /api/v1/analytics/statuses", authMiddleware.AuthMiddleware(analyticsHandler.GetOrderStatsByStatus))

	// ----- SWAGGER -----
	http.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
	http.HandleFunc("GET /swagger/*", httpSwagger.WrapHandler)

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