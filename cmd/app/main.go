package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"practice/config"
	"practice/env"
	"practice/internal"
	"practice/pkg/logger"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

func init() {
	config.LoadEnv()
	env.GetEnv()
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg := logger.DefaultConfig()
	logger, err := logger.NewLogger(cfg, "")
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer logger.Sync()

	app := fiber.New(fiber.Config{
		IdleTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,
	})

	port := fmt.Sprintf("0.0.0.0:%d", env.Port)

	db := config.NewDB(ctx, logger)
	defer db.Close()

	internal.MainRoutes(app, db, logger)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World!")
	})

	logger.Fatal(app.Listen(port).Error())
	logger.Info("Server running on port", port)

	<-ctx.Done()
	logger.Warn("Server is shutting down...")
	app.Shutdown()
}
