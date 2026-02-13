package handler

import (
	"log"
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
// @Security     AuthKey
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
    DisplayName      string `json:"display_name" validate:"required,min=2,no_scripts"`
    Bio              string `json:"bio" validate:"max=500,no_scripts"`
	SkillLevel       int `json:"skill_level" validate:"required,min=2,no_scripts"`     
    IsTemplateSeller bool   `json:"is_template_seller" validate:"required,min=2,no_scripts"`
}

// UpdateProfile godoc
// @Summary      Update Profile
// @Description  Update your display name, bio, and seller settings.
// @Tags         profile
// @Security     AuthKey
// @Accept       json
// @Produce      json
// @Param        profile  body      UpdateProfileRequest  true  "Profile update data"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Router       /profile/me [put]
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
    // 1. Get UserID from Middleware
    userID, exists := c.Get("userID") 
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "user not identified"})
        return
    }

    // 2. Bind Request Body
    var req UpdateProfileRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format: check your fields"})
        return
    }

    // 3. Map Fields to Model (Crucial: Match your Repo's $1-$5)
    profile := &models.Profile{
        UserID:           userID.(string),
        DisplayName:      req.DisplayName,
        Bio:              req.Bio,
        SkillLevel:       req.SkillLevel,       
        IsTemplateSeller: req.IsTemplateSeller, 
    }

    // 4. Call Service
    if err := h.svc.UpdateProfile(c.Request.Context(), profile); err != nil {
        // Log the actual error to terminal so you can see if SQL fails
        log.Printf("UpdateProfile Error: %v", err) 
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "profile updated successfully"})
}