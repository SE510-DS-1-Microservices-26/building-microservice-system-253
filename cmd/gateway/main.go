package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

func main() {
	// fetch and check service urls from env
	coreURL := os.Getenv("CORE_SERVICE_URL")
	usersURL := os.Getenv("USERS_SERVICE_URL")
	workflowURL := os.Getenv("WORKFLOW_SERVICE_URL")

	if coreURL == "" {
		log.Fatal("CORE_SERVICE_URL is required")
	}
	if usersURL == "" {
		log.Fatal("USERS_SERVICE_URL is required")
	}
	if workflowURL == "" {
		log.Fatal("WORKFLOW_SERVICE_URL is required")
	}

	// fiber engine
	app := fiber.New(fiber.Config{AppName: "gateway"})

	// logger middleware
	app.Use(logger.New(logger.Config{
		Format:     "[${ip}]:${port} ${method} ${path} ${status} ${latency}\n",
		TimeFormat: "15:04:05 02-01-2006",
	}))

	// swagger docs routes
	app.Get("/core/docs/*", func(c *fiber.Ctx) error {
		stripped := strings.TrimPrefix(c.OriginalURL(), "/core")
		return proxy.Do(c, coreURL+stripped)
	})

	app.Get("/users/docs/*", func(c *fiber.Ctx) error {
		stripped := strings.TrimPrefix(c.OriginalURL(), "/users")
		return proxy.Do(c, usersURL+stripped)
	})

	// routes
	app.All("/core/*", func(c *fiber.Ctx) error {
		return proxy.Do(c, coreURL+c.OriginalURL())
	})

	app.All("/users/*", func(c *fiber.Ctx) error {
		return proxy.Do(c, usersURL+c.OriginalURL())
	})

	app.All("/workflow/*", func(c *fiber.Ctx) error {
		return proxy.Do(c, workflowURL+c.OriginalURL())
	})

	// start server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8000"
	}

	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
			log.Fatalf("failed to start gateway: %v", err)
		}
	}()

	// graceful shutdown
	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-shutdownCtx.Done()
	fmt.Println("shutting down gateway...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("failed to shutdown gateway: %v", err)
	}
}
