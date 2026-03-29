package main

import (
	"cafeteria-delivery/internal/base/core/services"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"cafeteria-delivery/internal/base/adapters"
	"cafeteria-delivery/internal/base/config"
	customErrors "cafeteria-delivery/internal/base/errors"
	"cafeteria-delivery/internal/base/http/handlers"
	"cafeteria-delivery/internal/base/http/middlewares"
	"cafeteria-delivery/internal/base/http/routes"
	"cafeteria-delivery/internal/base/repositories"
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

	// repositories
	itemCategoryRepo := repositories.NewItemCategoryRepository(pool)
	itemRepo := repositories.NewItemRepository(pool)
	orderRepo := repositories.NewOrderRepository(pool)

	// adapters
	usersClient := adapters.NewUsersClient(os.Getenv("USERS_SERVICE_URL"))

	// services
	itemCategoryService := services.NewItemCategoryService(itemCategoryRepo)
	itemService := services.NewItemService(itemRepo)
	orderService := services.NewOrderService(orderRepo, itemRepo, usersClient)

	// handlers
	itemCategoryHandlers := handlers.NewItemCategoryHandlers(itemCategoryService, paginationConfig)
	itemHandlers := handlers.NewItemHandlers(itemService, paginationConfig)
	orderHandlers := handlers.NewOrderHandlers(orderService, paginationConfig)

	// fiber engine
	app := fiber.New(fiber.Config{
		AppName:      "cafeteria-delivery",
		ErrorHandler: errorHandler,
	})

	// middlewares
	middlewares.RouteLoggerMiddleware(app)

	// routes
	routes.HealthRoutes(app)
	routes.DocsRoutes(app)
	core := app.Group("/core")
	routes.ItemCategoryRoutes(core, itemCategoryHandlers)
	routes.ItemRoutes(core, itemHandlers)
	routes.OrderRoutes(core, orderHandlers)

	// print routes
	for _, route := range app.GetRoutes(true) {
		fmt.Printf("%-10s | %-50s | %s\n", route.Method, route.Path, route.Name)
	}

	// start server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
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
	v.SetConfigName("config")
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
