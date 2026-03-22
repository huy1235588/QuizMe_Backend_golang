package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/huy/quizme-backend/internal/dto/response"
	"github.com/huy/quizme-backend/internal/middleware"
	"github.com/huy/quizme-backend/internal/service"
	"github.com/huy/quizme-backend/internal/service/storage"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService       service.UserService
	cloudinaryService *storage.CloudinaryService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService service.UserService, cloudinaryService *storage.CloudinaryService) *UserHandler {
	return &UserHandler{
		userService:       userService,
		cloudinaryService: cloudinaryService,
	}
}

// GetUserByID handles getting a user by ID
// GET /api/users/:id
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid user ID"))
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, response.Error("User not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error("Failed to get user"))
		return
	}

	c.JSON(http.StatusOK, response.Success(user, "User retrieved successfully"))
}

// GetTopUsers handles getting top users by quiz plays
// GET /api/users/top
func (h *UserHandler) GetTopUsers(c *gin.Context) {
	limit := 10 // default
	if l, err := strconv.Atoi(c.DefaultQuery("limit", "10")); err == nil && l > 0 {
		limit = l
	}

	users, err := h.userService.GetTopUsers(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to get top users"))
		return
	}

	c.JSON(http.StatusOK, response.Success(users, "Top users retrieved successfully"))
}

// GetUserCount handles getting total user count
// GET /api/users/count
func (h *UserHandler) GetUserCount(c *gin.Context) {
	count, err := h.userService.GetUserCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to get user count"))
		return
	}

	c.JSON(http.StatusOK, response.Success(map[string]int64{"count": count}, "User count retrieved successfully"))
}

// GetPagedUsers handles getting paginated users
// GET /api/users/paged
func (h *UserHandler) GetPagedUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	search := c.Query("search")
	sortBy := c.DefaultQuery("sortBy", "created_at")
	sortDir := c.DefaultQuery("sortDir", "desc")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, total, err := h.userService.GetPagedUsers(page, pageSize, search, sortBy, sortDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to get users"))
		return
	}

	pageResponse := response.PageResponse[*response.UserResponse]{
		Content:       users,
		Page:          page,
		Size:          pageSize,
		TotalElements: total,
		TotalPages:    (total + int64(pageSize) - 1) / int64(pageSize),
	}

	c.JSON(http.StatusOK, response.Success(pageResponse, "Users retrieved successfully"))
}

// GetUserProfile handles getting a user's profile by ID
// GET /api/users/profile/:id
func (h *UserHandler) GetUserProfile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid user ID"))
		return
	}

	user, err := h.userService.GetUserProfile(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, response.Error("User not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error("Failed to get user profile"))
		return
	}

	c.JSON(http.StatusOK, response.Success(response.FromUser(user), "User profile retrieved successfully"))
}

// GetCurrentUserProfile handles getting the current user's profile
// GET /api/users/profile
func (h *UserHandler) GetCurrentUserProfile(c *gin.Context) {
	currentUser := middleware.GetCurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, response.Error("Authentication required"))
		return
	}

	user, err := h.userService.GetUserProfile(currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to get profile"))
		return
	}

	c.JSON(http.StatusOK, response.Success(response.FromUser(user), "Profile retrieved successfully"))
}

// UploadAvatar handles avatar upload for current user
// POST /api/users/avatar/upload
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	currentUser := middleware.GetCurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, response.Error("Authentication required"))
		return
	}

	// Get file from form data
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("No file uploaded"))
		return
	}
	defer file.Close()

	// Validate file size (max 5MB)
	if fileHeader.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, response.Error("File size must be less than 5MB"))
		return
	}

	// Validate file type
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/jpg" && contentType != "image/webp" {
		c.JSON(http.StatusBadRequest, response.Error("Only JPEG, PNG, and WebP images are allowed"))
		return
	}

	// Get the user profile to use the profile ID
	user, err := h.userService.GetUserProfile(currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to get user profile"))
		return
	}

	// Get profile ID - if UserProfile is nil, use user ID as fallback
	profileID := currentUser.ID
	if user.UserProfile != nil {
		profileID = user.UserProfile.ID
	}

	// Upload to Cloudinary
	filename, err := h.cloudinaryService.UploadProfileImage(c.Request.Context(), file, fileHeader.Filename, profileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to upload avatar"))
		return
	}

	// Get the full URL
	avatarURL := h.cloudinaryService.GetProfileImageURL(filename)

	// Update user avatar in database
	if err := h.userService.UpdateUserAvatar(currentUser.ID, avatarURL); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to update avatar"))
		return
	}

	c.JSON(http.StatusOK, response.Success(map[string]string{"avatarUrl": avatarURL}, "Avatar uploaded successfully"))
}

// RemoveAvatar handles avatar removal for current user
// DELETE /api/users/avatar
func (h *UserHandler) RemoveAvatar(c *gin.Context) {
	currentUser := middleware.GetCurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, response.Error("Authentication required"))
		return
	}

	// Get current user profile to get the avatar filename
	user, err := h.userService.GetUserProfile(currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to get user profile"))
		return
	}

	// Delete from Cloudinary if avatar exists
	if user.ProfileImage != nil && *user.ProfileImage != "" {
		// Extract filename from URL
		// This is a simple extraction, you might need to adjust based on your URL format
		// For now, we'll just attempt to delete and ignore errors if the file doesn't exist
		_ = h.cloudinaryService.DeleteProfileImage(c.Request.Context(), *user.ProfileImage)
	}

	// Remove from database
	if err := h.userService.RemoveUserAvatar(currentUser.ID); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to remove avatar"))
		return
	}

	c.JSON(http.StatusOK, response.Success[any](nil, "Avatar removed successfully"))
}

// CreateUser handles user creation (admin only)
// POST /api/users/create
func (h *UserHandler) CreateUser(c *gin.Context) {
	// TODO: Implement user creation
	c.JSON(http.StatusNotImplemented, response.Error("User creation not implemented yet"))
}

// UpdateUser handles user update (admin only)
// PUT /api/users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	// TODO: Implement user update
	c.JSON(http.StatusNotImplemented, response.Error("User update not implemented yet"))
}

// DeleteUser handles user deletion (admin only)
// DELETE /api/users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid user ID"))
		return
	}

	if err := h.userService.DeleteUser(uint(id)); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, response.Error("User not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error("Failed to delete user"))
		return
	}

	c.JSON(http.StatusOK, response.Success[any](nil, "User deleted successfully"))
}

// ToggleUserActiveStatus handles locking/unlocking a user (admin only)
// PUT /api/users/:id/lock
func (h *UserHandler) ToggleUserActiveStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid user ID"))
		return
	}

	var req struct {
		IsActive bool `json:"isActive"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid request body"))
		return
	}

	if err := h.userService.ToggleUserActiveStatus(uint(id), req.IsActive); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, response.Error("User not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error("Failed to update user status"))
		return
	}

	c.JSON(http.StatusOK, response.Success[any](nil, "User status updated successfully"))
}
