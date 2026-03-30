package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"cafeteria-delivery/internal/workflow/adapters"
	"cafeteria-delivery/internal/workflow/core/services"
	customErrors "cafeteria-delivery/internal/workflow/errors"
	"cafeteria-delivery/internal/workflow/http/handlers"
	"cafeteria-delivery/internal/workflow/http/middlewares"
	"cafeteria-delivery/internal/workflow/http/routes"
	"cafeteria-delivery/internal/workflow/repositories"
	"cafeteria-delivery/pkg/postgres"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func main() {
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

	usersClient := adapters.NewUsersClient(os.Getenv("USERS_SERVICE_URL"))
	coreClient := adapters.NewCoreClient(os.Getenv("CORE_SERVICE_URL"))

	workflowRepo := repositories.NewWorkflowRepository(pool)
	workflowService := services.NewWorkflowService(workflowRepo, usersClient, coreClient)
	workflowHandlers := handlers.NewWorkflowHandlers(workflowService)

	app := fiber.New(fiber.Config{
		AppName:      "workflow-service",
		ErrorHandler: errorHandler,
	})

	middlewares.RouteLoggerMiddleware(app)

	routes.HealthRoutes(app)
	routes.WorkflowRoutes(app, workflowHandlers)

	for _, route := range app.GetRoutes(true) {
		fmt.Printf("%-10s | %-50s | %s\n", route.Method, route.Path, route.Name)
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8082"
	}

	go func() {
		if err = app.Listen(fmt.Sprintf(":%s", port)); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-shutdownCtx.Done()
	fmt.Println("shutting down server...")
	if err = app.Shutdown(); err != nil {
		log.Fatalf("failed to shutdown server: %v", err)
	}
}

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
