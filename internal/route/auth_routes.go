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

func SetupAuthRoutes(r fiber.Router, config *config.Config, db *pgxpool.Pool, validator *auth.Validator) {
	authRepository := repository.NewAuthRepository(config, db)
	authService := service.NewAuthService(authRepository)
	authHandler := http.NewAuthHandler(authService, validator)

	{
		r.Post("/admin-login", authHandler.AdminLogin)
		r.Post("/send-application", authHandler.Application)
		r.Post("/send-application-document", auth.TokenGuard, authHandler.ApplicationDocuments)
		r.Post("/user-login-google", authHandler.UserLoginGoogle)
		r.Post("/user-login-email", authHandler.UserLoginEmail)
		r.Post("/user-forget-password", authHandler.UserForgetPassword)
		r.Post("/user-reset-password", authHandler.UserResetPassword)
		r.Post("/third-party-login", authHandler.ThirdPartyLogin)
		r.Post("/user-email-confirmation", authHandler.UserEmailConfirmation)
		r.Post("/user-login-phone", authHandler.UserLoginPhone)
		r.Post("/user-phone-confirmation", authHandler.UserPhoneConfirmation)
		r.Post("/user-register-device", auth.TokenGuard, authHandler.UserRegisterDevice)
		r.Delete("/account/:id", auth.TokenGuard, authHandler.DeleteAccount)
	}

}
