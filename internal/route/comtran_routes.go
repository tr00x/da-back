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

func SetupComtranRoutes(r fiber.Router, config *config.Config, db *pgxpool.Pool, validator *auth.Validator) {
	comtransRepository := repository.NewComtransRepository(config, db)
	comtransService := service.NewComtransService(comtransRepository)
	comtransHandler := http.NewComtransHandler(comtransService, validator)

	{
		// get comtrans categories
		r.Get("/categories", auth.TokenGuard, auth.LanguageChecker, comtransHandler.GetComtransCategories)
		r.Get("/categories/:category_id/parameters", auth.TokenGuard, auth.LanguageChecker, comtransHandler.GetComtransParameters)
		r.Get("/categories/:category_id/brands", auth.TokenGuard, auth.LanguageChecker, comtransHandler.GetComtransBrands)
		r.Get("/categories/:category_id/brands/:brand_id/models", auth.TokenGuard, auth.LanguageChecker, comtransHandler.GetComtransModelsByBrandID)

		// comtrans
		r.Get("/", auth.TokenGuard, comtransHandler.GetComtrans)
		r.Get("/:id", auth.TokenGuard, comtransHandler.GetComtransByID)
		r.Get("/:id/edit", auth.TokenGuard, comtransHandler.GetEditComtransByID)
		r.Post("/", auth.TokenGuard, comtransHandler.CreateComtrans)
		r.Post("/:id/images", auth.TokenGuard, comtransHandler.CreateComtransImages)
		r.Post("/:id/videos", auth.TokenGuard, comtransHandler.CreateComtransVideos)
		r.Post("/:id/buy", auth.TokenGuard, comtransHandler.BuyComtrans)
		r.Post("/:id/dont-sell", auth.TokenGuard, comtransHandler.DontSellComtrans)
		r.Post("/:id/sell", auth.TokenGuard, comtransHandler.SellComtrans)
		r.Delete("/:id/images", auth.TokenGuard, comtransHandler.DeleteComtransImage)
		r.Delete("/:id/videos", auth.TokenGuard, comtransHandler.DeleteComtransVideo)
		r.Delete("/:id", auth.TokenGuard, comtransHandler.DeleteComtrans)
	}
}
