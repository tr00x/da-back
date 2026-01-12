package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"dubai-auto/internal/config"
	"dubai-auto/internal/delivery/http"
	"dubai-auto/internal/repository"
	"dubai-auto/internal/service"
	"dubai-auto/pkg/auth"
)

func SetupMotorcycleRoutes(r fiber.Router, config *config.Config, db *pgxpool.Pool, validator *auth.Validator) {
	motorcycleRepository := repository.NewMotorcycleRepository(config, db)
	motorcycleService := service.NewMotorcycleService(motorcycleRepository)
	motorcycleHandler := http.NewMotorcycleHandler(motorcycleService, validator)

	{
		// get motorcycles categories
		r.Get("/categories", auth.TokenGuard, auth.LanguageChecker, motorcycleHandler.GetMotorcycleCategories)
		r.Get("/categories/:category_id/parameters", auth.TokenGuard, auth.LanguageChecker, motorcycleHandler.GetMotorcycleParameters)
		r.Get("/categories/:category_id/brands", auth.TokenGuard, auth.LanguageChecker, motorcycleHandler.GetMotorcycleBrands)
		r.Get("/categories/:category_id/brands/:brand_id/models", auth.TokenGuard, auth.LanguageChecker, motorcycleHandler.GetMotorcycleModelsByBrandID)

		// motorcycles
		r.Get("/", auth.TokenGuard, auth.LanguageChecker, motorcycleHandler.GetMotorcycles)
		r.Get("/:id", auth.TokenGuard, auth.LanguageChecker, motorcycleHandler.GetMotorcycleByID)
		r.Get("/:id/edit", auth.TokenGuard, auth.LanguageChecker, motorcycleHandler.GetEditMotorcycleByID)
		r.Post("/", auth.TokenGuard, motorcycleHandler.CreateMotorcycle)
		r.Post("/:id/images", auth.TokenGuard, motorcycleHandler.CreateMotorcycleImages)
		r.Post("/:id/videos", auth.TokenGuard, motorcycleHandler.CreateMotorcycleVideos)
		r.Post("/:id/buy", auth.TokenGuard, motorcycleHandler.BuyMotorcycle)
		r.Post("/:id/dont-sell", auth.TokenGuard, motorcycleHandler.DontSellMotorcycle)
		r.Post("/:id/sell", auth.TokenGuard, motorcycleHandler.SellMotorcycle)
		r.Delete("/:id/images", auth.TokenGuard, motorcycleHandler.DeleteMotorcycleImage)
		r.Delete("/:id/videos", auth.TokenGuard, motorcycleHandler.DeleteMotorcycleVideo)
		r.Delete("/:id", auth.TokenGuard, motorcycleHandler.DeleteMotorcycle)
	}
}
