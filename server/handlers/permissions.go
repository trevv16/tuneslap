package handlers

import (
	"tuneslap/models"
	"tuneslap/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	RoleAuthor = "author"
	RoleOwner  = "owner"
	RoleEditor = "editor"
	RoleViewer = "viewer"
)

// CheckBoardAccess checks if a user has access to a board (as author or collaborator)
// Returns the board if access is granted, otherwise returns an error
func CheckBoardAccess(boardId primitive.ObjectID, userId primitive.ObjectID) (models.Board, error) {
	boardRepo := repositories.NewBoardRepository()
	board, err := boardRepo.FindByIDWithAccess(boardId, userId)
	if err != nil {
		return models.Board{}, err
	}
	return board, nil
}

// GetUserBoardRole returns the user's role for a board
// Returns "author" if user is the board author, or the collaborator role if they're a collaborator
// Returns empty string if user has no access
func GetUserBoardRole(board models.Board, userId primitive.ObjectID) string {
	// Check if user is the author
	if board.AuthorId == userId {
		return RoleAuthor
	}

	// Check if user is a collaborator
	for _, collaborator := range board.Collaborators {
		if collaborator.UserId == userId {
			return collaborator.Role
		}
	}

	return ""
}

// CanEditBoard checks if a user can edit a board
// Returns true if user is author, or has "owner" or "editor" role
func CanEditBoard(board models.Board, userId primitive.ObjectID) bool {
	role := GetUserBoardRole(board, userId)
	return role == RoleAuthor || role == RoleOwner || role == RoleEditor
}

// CanDeleteBoard checks if a user can delete a board
// Returns true if user is author, or has "owner" role
func CanDeleteBoard(board models.Board, userId primitive.ObjectID) bool {
	role := GetUserBoardRole(board, userId)
	return role == RoleAuthor || role == RoleOwner
}

// CanManageCollaborators checks if a user can manage collaborators
// Returns true if user is author, or has "owner" role
func CanManageCollaborators(board models.Board, userId primitive.ObjectID) bool {
	role := GetUserBoardRole(board, userId)
	return role == RoleAuthor || role == RoleOwner
}

// CanManageKeys checks if a user can manage keys
// Returns true if user is author, or has "owner" or "editor" role
func CanManageKeys(board models.Board, userId primitive.ObjectID) bool {
	role := GetUserBoardRole(board, userId)
	return role == RoleAuthor || role == RoleOwner || role == RoleEditor
}

