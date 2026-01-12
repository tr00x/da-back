package route

import (
	"dubai-auto/internal/config"
	"dubai-auto/internal/delivery/http"
	"dubai-auto/internal/repository"
	"dubai-auto/internal/service"
	"dubai-auto/pkg/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupThirdPartyRoutes(r fiber.Router, config *config.Config, db *pgxpool.Pool, validator *auth.Validator) {
	thirdPartyRepository := repository.NewThirdPartyRepository(config, db)
	thirdPartyService := service.NewThirdPartyService(thirdPartyRepository)
	thirdPartyHandler := http.NewThirdPartyHandler(thirdPartyService, validator)

	{
		r.Get("/registration-data", auth.LanguageChecker, thirdPartyHandler.GetRegistrationData)
		r.Get("/profile", auth.TokenGuard, auth.LanguageChecker, thirdPartyHandler.GetProfile)
		r.Get("/profile/my-cars", auth.TokenGuard, auth.LanguageChecker, thirdPartyHandler.GetMyCars)
		r.Get("/profile/on-sale", auth.TokenGuard, auth.LanguageChecker, thirdPartyHandler.OnSale)
		r.Post("/first-login", auth.TokenGuard, thirdPartyHandler.FirstLogin)
		r.Post("/profile/banner", auth.TokenGuard, thirdPartyHandler.BannerImage)
		r.Delete("/profile/banner", auth.TokenGuard, thirdPartyHandler.DeleteBannerImage)
		r.Post("/profile/images", auth.TokenGuard, thirdPartyHandler.AvatarImages)
		r.Delete("/profile/images", auth.TokenGuard, thirdPartyHandler.DeleteAvatarImages)
		r.Post("/profile", auth.TokenGuard, thirdPartyHandler.Profile)

		// dealer routes
		r.Post("/dealer/car", auth.TokenGuard, auth.DealerGuard, thirdPartyHandler.CreateDealerCar)
		r.Get("/dealer/car/:id/edit", auth.TokenGuard, auth.DealerGuard, auth.LanguageChecker, thirdPartyHandler.GetEditCarByID)
		r.Post("/dealer/car/:id/sell", auth.TokenGuard, auth.DealerGuard, thirdPartyHandler.StatusDealer)
		r.Post("/dealer/car/:id/dont-sell", auth.TokenGuard, auth.DealerGuard, thirdPartyHandler.StatusDealer)
		r.Post("/dealer/car/:id", auth.TokenGuard, auth.DealerGuard, thirdPartyHandler.UpdateDealerCar)
		r.Delete("/dealer/car/:id", auth.TokenGuard, auth.DealerGuard, thirdPartyHandler.DeleteDealerCar)
		r.Post("/dealer/car/:id/images", auth.TokenGuard, auth.DealerGuard, thirdPartyHandler.CreateDealerCarImages)
		r.Post("/dealer/car/:id/videos", auth.TokenGuard, auth.DealerGuard, thirdPartyHandler.CreateDealerCarVideos)
		r.Delete("/dealer/car/:id/images", auth.TokenGuard, auth.DealerGuard, thirdPartyHandler.DeleteDealerCarImage)
		r.Delete("/dealer/car/:id/videos", auth.TokenGuard, auth.DealerGuard, thirdPartyHandler.DeleteDealerCarVideo)

		// logist routes
		r.Get("/logist/destinations", auth.TokenGuard, auth.LogistGuard, auth.LanguageChecker, thirdPartyHandler.GetLogistDestinations)
		r.Post("/logist/destinations", auth.TokenGuard, auth.LogistGuard, thirdPartyHandler.CreateLogistDestination)
		r.Delete("/logist/destinations/:id", auth.TokenGuard, auth.LogistGuard, thirdPartyHandler.DeleteLogistDestination)
		// broker routes
		// car service routes
	}
}
