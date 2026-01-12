package route

import (
	"dubai-auto/internal/config"
	"dubai-auto/internal/delivery/http"
	"dubai-auto/internal/repository"
	"dubai-auto/internal/service"
	"dubai-auto/pkg/auth"
	"dubai-auto/pkg/firebase"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupWebSocketRoutes(app *fiber.App, db *pgxpool.Pool, firebaseService *firebase.FirebaseService, config *config.Config) {

	socketRepository := repository.NewSocketRepository(db, firebaseService, config)
	socketService := service.NewSocketService(socketRepository)
	socketHandler := http.NewSocketHandler(socketService)

	app.Get("/ws/conversations", auth.TokenGuard, socketHandler.GetConversations)
	app.Get("/ws/conversations/:conversation_id/messages", auth.TokenGuard, socketHandler.GetConversationMessages)
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Get("/ws", socketHandler.SetupWebSocketHandler())
}
