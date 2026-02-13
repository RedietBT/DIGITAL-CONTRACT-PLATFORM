package handler

import (
	"net/http"

	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/models"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	svc service.ProfileService
}

func NewProfileHandler(svc service.ProfileService) *ProfileHandler {
	return &ProfileHandler{svc: svc}
}

// GetProfile retrieves the authenticated user's profile
// @Summary      Get My Profile
// @Description  Fetches the profile for the user identified by the JWT token.
// @Tags         profile
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  models.Profile
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      404  {object}  map[string]string "Not Found"
// @Router       /profile/me [get]
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	// 1. EXTRACT: pull the ID from the "bucket" filled by profileAuthMiddleWare
	 userID, exists := c.Get("userID")
	 if !exists {
		// This should theorectically never happen if middleware is working
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id missing from context"})
		return
	 }

	 // 2. CALL SERVICE: We use c.Request.Context() to support cancellations/timeouts
	 profile, err := h.svc.GetProfile(c.Request.Context(), userID.(string))
	 if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "profile not found"})
		return
	 }

	 // 3. RESPOND: Gin automaticaly converts struct to JSON
	 c.JSON(http.StatusOK, profile)
}

// UpdateProfileRequest defines exactly what a user can change
type UpdateProfileRequest struct {
    DisplayName string `json:"display_name" validate:"required,min=2,no_scripts"`
    Bio         string `json:"bio" validate:"max=500,no_scripts"`
}

// UpdateProfile godoc
// @Summary      Update Profile
// @Description  Update your display name and bio.
// @Tags         profile
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        profile  body      UpdateProfileRequest  true  "Profile update data"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Router       /profile/me [put]
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
    userID, _ := c.Get("userID") // Get from standardized key

    var req UpdateProfileRequest
    // Bind to the DTO instead of the full Model
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format: check your fields"})
        return
    }

    // Map the safe fields to your profile model
    profile := &models.Profile{
        UserID:      userID.(string),
        DisplayName: req.DisplayName,
        Bio:         req.Bio,
    }

    if err := h.svc.UpdateProfile(c.Request.Context(), profile); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "profile updated successfully"})
}