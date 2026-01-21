package handlers

import (
	"tuneslap/config"
	api "tuneslap/generated"
	"tuneslap/models"
	"tuneslap/repositories"
	"tuneslap/utils"
	"tuneslap/validation"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BoardHandler provides board-specific operations using the base handler
type BoardHandler struct {
	*BaseHandler[models.Board, api.CreateBoardRequest, api.UpdateBoardRequest]
	boardRepo *repositories.BoardRepository
	mediaRepo *repositories.MediaRepository
}

// NewBoardHandler creates a new board handler
func NewBoardHandler() *BoardHandler {
	boardRepo := repositories.NewBoardRepository()
	mediaRepo := repositories.NewMediaRepository()
	boardValidator := validation.NewBoardValidator()

	return &BoardHandler{
		BaseHandler: NewBaseHandler[models.Board, api.CreateBoardRequest, api.UpdateBoardRequest](
			boardRepo.Repository,
			boardValidator,
		),
		boardRepo: boardRepo,
		mediaRepo: mediaRepo,
	}
}

// BoardResponseMapper maps board model to response format with media URLs
func (h *BoardHandler) BoardResponseMapper(board models.Board) interface{} {
	// Convert board to API response
	boardResponse := utils.ToBoardResponse(board)

	// Collect all media IDs from keys
	var mediaIds []primitive.ObjectID
	for _, key := range board.Keys {
		if !key.AudioMediaId.IsZero() {
			mediaIds = append(mediaIds, key.AudioMediaId)
		}
		if !key.ImageMediaId.IsZero() {
			mediaIds = append(mediaIds, key.ImageMediaId)
		}
	}

	// Fetch media URLs
	mediaUrls := make(map[primitive.ObjectID]string)
	if len(mediaIds) > 0 {
		urls, err := h.mediaRepo.GetUrlsForMedia(mediaIds)
		if err == nil {
			mediaUrls = urls
		}
	}

	// Enrich keys with URLs
	if boardResponse.Keys != nil {
		for i, key := range boardResponse.Keys {
			// Find corresponding internal key to get media IDs
			if i < len(board.Keys) {
				internalKey := board.Keys[i]

				// Set audio URL if available
				if !internalKey.AudioMediaId.IsZero() {
					if url, exists := mediaUrls[internalKey.AudioMediaId]; exists {
						key.AudioUrl = &url
					}
				}

				// Set image URL if available
				if !internalKey.ImageMediaId.IsZero() {
					if url, exists := mediaUrls[internalKey.ImageMediaId]; exists {
						key.ImageUrl = &url
					}
				}

				boardResponse.Keys[i] = key
			}
		}
	}

	return boardResponse
}

// HandleGetAllBoards handles GET all boards with pagination
func (h *BoardHandler) HandleGetAllBoards(c *fiber.Ctx) error {
	return h.HandleGetAll(
		c,
		func(authorId primitive.ObjectID) bson.M {
			return bson.M{"authorId": authorId}
		},
		h.BoardResponseMapper,
		nil,
		"boards", // dataFieldName to match OpenAPI spec
	)
}

// HandleGetBoardById handles GET board by ID
func (h *BoardHandler) HandleGetBoardById(c *fiber.Ctx) error {
	authorId, err := GetAuthorId(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid author ID", err)
	}

	return h.HandleGetById(c, "boardId", authorId, h.BoardResponseMapper, nil)
}

// HandleCreateBoard handles CREATE board
func (h *BoardHandler) HandleCreateBoard(c *fiber.Ctx) error {
	authorId, err := GetAuthorId(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid author ID", err)
	}

	return h.HandleCreate(c, authorId, func(req api.CreateBoardRequest, authorId primitive.ObjectID) (models.Board, error) {
		return h.boardRepo.CreateBoard(&req, authorId)
	}, h.BoardResponseMapper)
}

// HandleUpdateBoard handles UPDATE board
func (h *BoardHandler) HandleUpdateBoard(c *fiber.Ctx) error {
	authorId, err := GetAuthorId(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid author ID", err)
	}

	return h.HandleUpdate(c, "boardId", authorId, func(id primitive.ObjectID, authorId primitive.ObjectID, req api.UpdateBoardRequest) (models.Board, error) {
		return h.boardRepo.UpdateBoard(id, authorId, &req)
	}, h.BoardResponseMapper)
}

// HandleDeleteBoard handles DELETE board
func (h *BoardHandler) HandleDeleteBoard(c *fiber.Ctx) error {
	authorId, err := GetAuthorId(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid author ID", err)
	}

	return h.HandleDelete(c, "boardId", authorId)
}

// HandleGetDemoBoard handles GET demo board (public, no auth required)
func (h *BoardHandler) HandleGetDemoBoard(c *fiber.Ctx) error {
	// Check if demo mode is enabled
	if !config.IsDemoMode() {
		return SendErrorResponse(c, fiber.StatusNotFound, "Demo mode is not enabled", nil)
	}

	// Fetch the demo board using the hardcoded ID
	board, err := h.boardRepo.FindByID(config.DemoBoardID)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusNotFound, "Demo board not found", err)
	}

	// Use the same response mapper for consistency
	response := h.BoardResponseMapper(board)

	return c.JSON(response)
}
