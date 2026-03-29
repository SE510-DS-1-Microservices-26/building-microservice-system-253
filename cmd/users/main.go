package main

import (
	"cafeteria-delivery/internal/users/core/services"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"cafeteria-delivery/internal/users/config"
	customErrors "cafeteria-delivery/internal/users/errors"
	"cafeteria-delivery/internal/users/http/handlers"
	"cafeteria-delivery/internal/users/http/middlewares"
	"cafeteria-delivery/internal/users/http/routes"
	"cafeteria-delivery/internal/users/repositories"
	"cafeteria-delivery/pkg/postgres"

	"github.com/fsnotify/fsnotify"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
)

func main() {
	// config
	v, err := loadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	paginationConfig := config.NewPaginationConfig(v)

	// dynamic config reload
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("config file changed: %s\n", e.Name)
		paginationConfig.UpdatePaginationConfig(v)
	})

	// database
	ctx := context.Background()
	pool, err := postgres.New(ctx, postgres.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_DATABASE"),
	})
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	defer pool.Close()

	// repository
	userRepo := repositories.NewUserRepository(pool)
	// service
	userService := services.NewUserService(userRepo)
	// handlers
	userHandlers := handlers.NewUserHandlers(userService, paginationConfig)

	// fiber engine
	app := fiber.New(fiber.Config{
		AppName:      "users-service",
		ErrorHandler: errorHandler,
	})

	// middlewares
	middlewares.RouteLoggerMiddleware(app)

	// routes
	routes.HealthRoutes(app)
	routes.DocsRoutes(app)
	routes.UserRoutes(app, userHandlers)

	// print routes
	for _, route := range app.GetRoutes(true) {
		fmt.Printf("%-10s | %-50s | %s\n", route.Method, route.Path, route.Name)
	}

	// start server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8081"
	}

	go func() {
		if err = app.Listen(fmt.Sprintf(":%s", port)); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	// graceful shutdown
	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-shutdownCtx.Done()
	fmt.Println("shutting down server...")
	if err = app.Shutdown(); err != nil {
		log.Fatalf("failed to shutdown server: %v", err)
	}
}

// loadConfig - helper function to load the config file
func loadConfig() (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName("users")
	v.SetConfigType("yaml")
	v.AddConfigPath("configs")

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	v.WatchConfig()

	return v, nil
}

// errorHandler - helper function for global error handling in fiber (for general errors)
func errorHandler(ctx *fiber.Ctx, err error) error {
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return ctx.Status(fiberErr.Code).JSON(fiber.Map{"message": fiberErr.Message})
	}

	switch {
	case errors.Is(err, customErrors.ErrNotFound):
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
	case errors.Is(err, customErrors.ErrValidation):
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": err.Error()})
	case errors.Is(err, customErrors.ErrConflict):
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"message": err.Error()})
	case errors.Is(err, customErrors.ErrInvalidID):
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	case errors.Is(err, pgx.ErrNoRows):
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": customErrors.ErrNotFound.Error()})
	}

	log.Printf("unhandled error: %v", err)
	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "internal server error"})
}
