package handlers

import (
	"tuneslap/config"
	api "tuneslap/generated"
	"tuneslap/models"
	"tuneslap/repositories"
	"tuneslap/utils"
	"tuneslap/validation"

	"github.com/gofiber/fiber/v2"
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
// Returns boards where user is author OR collaborator
func (h *BoardHandler) HandleGetAllBoards(c *fiber.Ctx) error {
	userId, err := GetAuthorId(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID", err)
	}

	// Validate pagination parameters
	pagination, err := validatePaginationParams(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid pagination parameters", err)
	}

	// Get all boards where user has access (author or collaborator)
	allBoards, err := h.boardRepo.FindAllWithAccess(userId)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Error retrieving boards", err)
	}

	// Apply pagination
	totalCount := int64(len(allBoards))
	start := pagination.Skip
	end := start + pagination.Limit
	if start > len(allBoards) {
		start = len(allBoards)
	}
	if end > len(allBoards) {
		end = len(allBoards)
	}

	var paginatedBoards []models.Board
	if start < len(allBoards) {
		paginatedBoards = allBoards[start:end]
	} else {
		paginatedBoards = []models.Board{}
	}

	// Map results
	mappedResults := make([]interface{}, len(paginatedBoards))
	for i, board := range paginatedBoards {
		mappedResults[i] = h.BoardResponseMapper(board)
	}

	// Create pagination metadata
	currentPage := 1
	if pagination.Limit > 0 {
		currentPage = (pagination.Skip / pagination.Limit) + 1
	}
	paginationMeta := PaginationMeta{
		CurrentPage:  currentPage,
		PageSize:     pagination.Limit,
		TotalResults: totalCount,
	}

	return SendPaginatedResponse(
		c,
		fiber.StatusOK,
		"Boards retrieved successfully",
		mappedResults,
		paginationMeta,
		"boards",
	)
}

// HandleGetBoardById handles GET board by ID
// Returns board if user has access (author or collaborator)
func (h *BoardHandler) HandleGetBoardById(c *fiber.Ctx) error {
	userId, err := GetAuthorId(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID", err)
	}

	// Parse board ID
	boardId, err := GetValidObjectId(c, "boardId")
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid board ID", err)
	}

	// Check board access (author or collaborator)
	board, err := CheckBoardAccess(boardId, userId)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusNotFound, "Board not found", err)
	}

	// Map and return board
	response := h.BoardResponseMapper(board)
	return c.Status(fiber.StatusOK).JSON(response)
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
// Requires user to be author, or collaborator with "owner" or "editor" role
func (h *BoardHandler) HandleUpdateBoard(c *fiber.Ctx) error {
	userId, err := GetAuthorId(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID", err)
	}

	// Parse board ID
	boardId, err := GetValidObjectId(c, "boardId")
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid board ID", err)
	}

	// Check board access
	board, err := CheckBoardAccess(boardId, userId)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusNotFound, "Board not found", err)
	}

	// Check edit permission
	if !CanEditBoard(board, userId) {
		return SendErrorResponse(c, fiber.StatusForbidden, "You do not have permission to edit this board", nil)
	}

	// Parse request body
	request := new(api.UpdateBoardRequest)
	if err := c.BodyParser(request); err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate request
	validationResult := h.validator.Validate(request)
	if !validationResult.IsValid {
		return SendValidationErrorResponse(c, validationResult)
	}

	// Update board
	updatedBoard, err := h.boardRepo.UpdateBoardWithAccess(boardId, request)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to update board", err)
	}

	// Map and return updated board
	response := h.BoardResponseMapper(updatedBoard)
	return SendSuccessResponse(c, fiber.StatusOK, "Board updated successfully", response)
}

// HandleDeleteBoard handles DELETE board
// Requires user to be author, or collaborator with "owner" role
func (h *BoardHandler) HandleDeleteBoard(c *fiber.Ctx) error {
	userId, err := GetAuthorId(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID", err)
	}

	// Parse board ID
	boardId, err := GetValidObjectId(c, "boardId")
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid board ID", err)
	}

	// Check board access
	board, err := CheckBoardAccess(boardId, userId)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusNotFound, "Board not found", err)
	}

	// Check delete permission
	if !CanDeleteBoard(board, userId) {
		return SendErrorResponse(c, fiber.StatusForbidden, "You do not have permission to delete this board", nil)
	}

	// Delete board
	err = h.boardRepo.DeleteBoardWithAccess(boardId)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete board", err)
	}

	return SendSuccessResponse(c, fiber.StatusOK, "Board deleted successfully", nil)
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
