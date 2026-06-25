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
	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/repository"
	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/service"
	"github.com/Mushka-pushka/flower-marketplace/backend/order-service/internal/worker"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
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

	// Объявляем очередь
	_, err = rabbitCh.QueueDeclare(
		"order.created", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}
	log.Println("RabbitMQ queue declared")

	// Инициализация
	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo, cfg, rabbitCh)
	orderHandler := handlers.NewOrderHandler(orderService)

	// Запуск воркера
	orderWorker := worker.NewOrderWorker(orderService, rabbitCh)
	go func() {
		if err := orderWorker.Start(context.Background()); err != nil {
			log.Printf("Worker error: %v", err)
		}
	}()

	// Настройка роутера
	http.HandleFunc("POST /api/v1/orders", orderHandler.CreateOrder)
	http.HandleFunc("GET /api/v1/orders", orderHandler.GetOrder)
	http.HandleFunc("GET /api/v1/orders/customer", orderHandler.GetOrdersByCustomer)
	http.HandleFunc("POST /api/v1/orders/cancel", orderHandler.CancelOrder)

	// Сервер
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