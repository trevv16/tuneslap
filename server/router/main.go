package router

import (
	"tuneslap/handlers"
	"tuneslap/services"

	"github.com/gofiber/fiber/v2"
)

// Handler instances
var (
	boardHandler        = handlers.NewBoardHandler()
	userHandler         = handlers.NewUserHandler()
	mediaHandler        = handlers.NewMediaHandler()
	collaboratorHandler = handlers.NewCollaboratorHandler()
	keyHandler          = handlers.NewKeyHandler()
)

func SetupRoutes(app *fiber.App) {
	app.Get("/health", handlers.HandleHealthCheck)

	// API routes
	apiV1 := app.Group("/api/v1")

	// auth
	auth := apiV1.Group("/auth")
	auth.Post("/signup", handlers.SignUpHandler)
	auth.Post("/signin", handlers.SignInHandler)
	auth.Post("/forgot", handlers.ForgotPasswordHandler)
	auth.Post("/reset", handlers.ResetPasswordHandler)

	// boards
	boards := apiV1.Group("/boards", services.JWTProtected())
	boards.Get("/", boardHandler.HandleGetAllBoards)
	boards.Post("/", boardHandler.HandleCreateBoard)
	boards.Get("/:boardId", boardHandler.HandleGetBoardById)
	boards.Patch("/:boardId", boardHandler.HandleUpdateBoard)
	boards.Delete("/:boardId", boardHandler.HandleDeleteBoard)

	// board collaborators
	boards.Get("/:boardId/collaborators", collaboratorHandler.HandleGetAllCollaborators)
	boards.Get("/:boardId/collaborators/:collaboratorId", collaboratorHandler.HandleGetCollaboratorById)
	boards.Post("/:boardId/collaborators", collaboratorHandler.HandleCreateCollaborator)
	boards.Patch("/:boardId/collaborators/:collaboratorId", collaboratorHandler.HandleUpdateCollaborator)
	boards.Delete("/:boardId/collaborators/:collaboratorId", collaboratorHandler.HandleDeleteCollaborator)

	// board keys
	boards.Get("/:boardId/keys", keyHandler.HandleGetAllKeys)
	boards.Get("/:boardId/keys/:keyId", keyHandler.HandleGetKeyById)
	boards.Post("/:boardId/keys", keyHandler.HandleCreateKey)
	boards.Patch("/:boardId/keys/:keyId", keyHandler.HandleUpdateKey)
	boards.Delete("/:boardId/keys/:keyId", keyHandler.HandleDeleteKey)

	// media
	media := apiV1.Group("/media", services.JWTProtected())
	media.Get("/stats", mediaHandler.HandleGetMyMediaStats)
	media.Get("/", mediaHandler.HandleGetAllMedia)
	media.Post("/upload-url", mediaHandler.HandleGenerateUploadURL)
	media.Post("/", mediaHandler.HandleCreateMedia)
	media.Post("/:mediaId/process", mediaHandler.HandleProcessMedia)
	media.Get("/:mediaId", mediaHandler.HandleGetMediaById)
	media.Patch("/:mediaId", mediaHandler.HandleUpdateMedia)
	media.Delete("/:mediaId", mediaHandler.HandleDeleteMedia)

	// users
	users := apiV1.Group("/users", services.JWTProtected())
	users.Get("/me", userHandler.HandleGetMe)
	users.Patch("/me", userHandler.HandleUpdateMe)

}
