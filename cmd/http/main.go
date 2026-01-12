package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/swagger"

	_ "dubai-auto/docs"
	app "dubai-auto/internal"
	"dubai-auto/internal/config"
	"dubai-auto/internal/storage/postgres"
	"dubai-auto/internal/utils"
	"dubai-auto/pkg/auth"
	"dubai-auto/pkg/logger"
)

// @title Project name
// @version 1.0
// @description Project Description
// @host api.mashynbazar.com
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	conf := config.Init()
	auth.Init(conf.ACCESS_KEY, conf.ACCESS_TIME, conf.REFRESH_KEY, conf.REFRESH_TIME)
	err := logger.InitLogger(conf.LOGGER_FOLDER_PATH, conf.LOGGER_FILENAME, conf.APP_MODE)

	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	db := postgres.Init(conf)
	app := app.InitApp(db, conf)

	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "/swagger/doc.json",
	}))

	if conf.MIGRATE == "true" {
		utils.MigrateV2(conf.MIGRATE_PATH, db)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Fiber server listening on %s", conf.PORT)

		if err := app.Listen(conf.PORT); err != nil {
			log.Fatalf("Fiber listen error: %v", err)
		}
	}()

	sig := <-quit
	log.Printf("Received signal: %s. Shutting down server...", sig)
	shutdownCtx := time.NewTicker(5 * time.Second)
	defer shutdownCtx.Stop()

	// Shutdown gracefully shuts down the server without interrupting any active connections.
	// err := app.Shutdown()

	// if err != nil {
	// 	log.Printf("Fiber graceful shutdown error: %v", err)
	// } else {
	// 	log.Println("Fiber server gracefully stopped.")
	// }

	log.Println("Application exited.")
}
