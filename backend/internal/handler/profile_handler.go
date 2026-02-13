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
	 userID, exists := c.Get("user_id")
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

// UpdateProfile updates the authenticated user's profile
// @Summary      Update Profile
// @Description  Allows the user to update their display name, bio, and other details.
// @Tags         profile
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        profile  body      models.Profile  true  "Profile data"
// @Success      200      {object}  map[string]string "message"
// @Failure      400      {object}  map[string]string "Invalid input"
// @Router       /profile/me [put]
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req models.Profile

	// 4. VALDATE: c.ShouldBindJSON checks the "json" tags and "binding" tags in  your model
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	} 

	// SECURITY: Ensure the user can only update their own profile by forcing the ID
	req.UserID = userID.(string)

	if err := h.svc.UpdateProfile(c.Request.Context(), &req); err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated successfully"})
}