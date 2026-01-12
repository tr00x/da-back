package internal

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	fb_logger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/jackc/pgx/v5/pgxpool"

	"dubai-auto/internal/config"
	"dubai-auto/internal/route"
	"dubai-auto/pkg/auth"
	"dubai-auto/pkg/firebase"
)

func InitApp(db *pgxpool.Pool, conf *config.Config) *fiber.App {
	firebaseService, err := firebase.InitFirebase(conf)

	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	// fbToken := "ei4Af4xzQgChqC-Rw6Emfa:APA91bGJ5L7PTh2KWR9atMbQANKJEoIKk4cVdT4VxCrc7grfgzzL1d72BaSZtq9uuAqGNZHiYI6xZZpcnQ45nyXX2PHPnxpqT4Y742Sw1eiHt-u1N32RfY4"
	validator := auth.NewValidator()
	appConfig := fiber.Config{
		BodyLimit: 50 * 1024 * 1024,
	}

	app := fiber.New(appConfig)
	app.Use(pprof.New())
	app.Use(auth.Cors)

	if config.ENV.APP_MODE != "release" {
		app.Use(fb_logger.New(fb_logger.Config{
			Format: "[${time}] ${ip} ${status} - ${method} ${path} ${latency}\n",
			Output: os.Stdout, // Or to a file: os.OpenFile("fiber.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		}))
	}

	app.Static("api/v1/images", "."+conf.STATIC_PATH)
	route.Init(app, conf, db, firebaseService, validator)
	return app
}
