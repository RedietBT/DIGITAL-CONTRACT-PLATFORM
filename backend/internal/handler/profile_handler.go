package handler

import "github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/service"

type ProfileHandler struct {
	service service.ProfileService
}

func NewProfileHandler(svc service.ProfileService) *ProfileHandler {
	return &ProfileHandler{service: svc}
}

// GetProfile godoc
// @Summary      Get user profile
// @Description  Retrieve profile details by User UUID
// @Tags         profiles
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  models.Profile
// @Failure      404  {object}  map[string]string
// @Router       /profiles/{id} [get]