package main

import (
	"cafeteria-delivery/internal/notification/adapters"
	"cafeteria-delivery/internal/notification/core/services"
	"cafeteria-delivery/internal/notification/repositories"
	"cafeteria-delivery/pkg/postgres"
	"cafeteria-delivery/pkg/rabbitmq"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()

	// database
	pool, err := postgres.New(ctx, postgres.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_DATABASE"),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// rabbitmq
	rabbitConn, err := rabbitmq.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatalf("failed to connect to rabbitmq: %v", err)
	}
	defer rabbitConn.Close()

	// repository
	repo := repositories.NewNotificationRepository(pool)
	// service
	service := services.NewNotificationService(repo)
	// consumer
	consumer, err := adapters.NewRabbitConsumer(rabbitConn, service)
	if err != nil {
		log.Fatalf("failed to initialize consumer: %v", err)
	}

	// graceful shutdown
	shutdownCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// start the consumer
	if err = consumer.Start(shutdownCtx); err != nil {
		log.Fatalf("consumer error: %v", err)
	}
}
