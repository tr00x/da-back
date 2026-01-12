package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "dubai-auto/docs"
	"dubai-auto/internal/config"
	"dubai-auto/pkg/auth"
	"dubai-auto/pkg/firebase"
)

func Init(app *fiber.App, config *config.Config, db *pgxpool.Pool, firebaseService *firebase.FirebaseService, validator *auth.Validator) {

	userRoute := app.Group("/api/v1/users")
	SetupUserRoutes(userRoute, config, db, validator)

	authRoute := app.Group("/api/v1/auth")
	SetupAuthRoutes(authRoute, config, db, validator)

	motorcycleRoute := app.Group("/api/v1/motorcycles")
	SetupMotorcycleRoutes(motorcycleRoute, config, db, validator)

	comtransRoute := app.Group("/api/v1/comtrans")
	SetupComtranRoutes(comtransRoute, config, db, validator)

	adminRoute := app.Group("/api/v1/admin", auth.TokenGuard, auth.AdminGuard)
	SetupAdminRoutes(adminRoute, config, db, validator)

	thirdPartyRoute := app.Group("/api/v1/third-party")
	SetupThirdPartyRoutes(thirdPartyRoute, config, db, validator)

	SetupWebSocketRoutes(app, db, firebaseService, config)

}
