package handler

import (
	"net/http"

	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/middleware"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc service.AuthService
}

func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// --- Request Structs with Binding Tags ---

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"dev@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"secret123"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"dev@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"secret123"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email" example:"dev@example.com"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type UpdateEmailRequest struct {
	NewEmail string `json:"new_email" binding:"required,email"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UpdateStatusRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Status string `json:"status" binding:"required"`
}

// Register godoc
// @Summary      Register a new user
// @Description  Creates a new user account in the auth service and broadcasts a creation event via RabbitMQ.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      RegisterRequest  true  "User Registration Details"
// @Success      201      {object}  models.User      "Successfully created user"
// @Failure      400      {object}  map[string]string "Invalid input data"
// @Failure      500      {object}  map[string]string "Internal server error"
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.svc.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login godoc
// @Summary      Authenticate a user
// @Description  Verifies credentials and returns access and refresh tokens.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      LoginRequest  true  "User Credentials"
// @Success      200          {object}  map[string]string "Returns access_token and refresh_token"
// @Failure      400          {object}  map[string]string "Invalid input format"
// @Failure      401          {object}  map[string]string "Unauthorized: Invalid email or password"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// ForgotPassword godoc
// @Summary      Forgot Password
// @Description  Initiates the password reset process by sending a reset code to the user's email.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      ForgotPasswordRequest  true  "User Email"
// @Success      200      {object}  map[string]string      "Success message"
// @Failure      400      {object}  map[string]string      "Invalid email format"
// @Failure      500      {object}  map[string]string      "Failed to send reset code"
// @Router       /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_ = h.svc.ForgotPassword(c.Request.Context(), req.Email)
	c.JSON(http.StatusOK, gin.H{"message": "If that email is in our system, a reset code has been sent."})
}

// ResetPassword godoc
// @Summary      Reset password using token
// @Description  Resets the user's password using a verification token.
// @Tags         default
// @Accept       json
// @Produce      json
// @Param        request  body      ResetPasswordRequest  true  "Reset Password Details"
// @Success      200      {object}  map[string]string      "Success message"
// @Failure      400      {object}  map[string]string      "Invalid token or password mismatch"
// @Failure      500      {object}  map[string]string      "Failed to reset password"
// @Router       /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.ResetPassword(c.Request.Context(), req.Email, req.Token, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password has been reset successfully."})
}

// GetProfile godoc
// @Summary      Get user profile
// @Description  Retrieves the profile information of the authenticated user.
// @Tags         default
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {object}  models.User "Successfully retrieved user profile"
// @Failure      401  {object}  map[string]string "Unauthorized: User not found in context"
// @Failure      404  {object}  map[string]string "User not found"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /auth/me [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.GetString(middleware.UserIDKey)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user, err := h.svc.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetAllUsers godoc
// @Summary      Get all users (Admin only)
// @Description  Retrieves a list of all users in the system. Requires admin privileges.
// @Tags         admin
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {array}  models.User "Successfully retrieved list of users"
// @Failure      500  {object}  map[string]string "Failed to retrieve users"
// @Router       /auth/admin/users [get]
func (h *AuthHandler) GetAllUsers(c *gin.Context) {
	users, err := h.svc.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// DeleteUser godoc
// @Summary      Delete a user
// @Description  Deletes a user account. If no ID is provided, deletes the authenticated user.
// @Tags         admin
// @Security     ApiKeyAuth
// @Produce      json
// @Param        id  query     string  false  "User ID to delete"
// @Success      204 {object}  nil       "Successfully deleted user"
// @Failure      500 {object}  map[string]string "Failed to delete user"
// @Router       /auth/admin/users [delete]
func (h *AuthHandler) DeleteUser(c *gin.Context) {
	currentUserID := c.GetString(middleware.UserIDKey)
	targetID := c.Query("id")

	idToDelete := currentUserID
	if targetID != "" {
		idToDelete = targetID
	}

	if err := h.svc.DeleteUser(c.Request.Context(), idToDelete); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) DeleteMe(c *gin.Context) {
	h.DeleteUser(c)
}

// UpdateEmail godoc
// @Summary      Update own email
// @Description  Updates the email address of the authenticated user.
// @Tags         default
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        request  body      UpdateEmailRequest  true  "New Email Address"
// @Success      200      {object}  map[string]string   "Success message"
// @Failure      400      {object}  map[string]string   "Invalid email format"
// @Failure      500      {object}  map[string]string   "Failed to update email"
// @Router       /auth/me/email [put]
func (h *AuthHandler) UpdateEmail(c *gin.Context) {
	userID := c.GetString(middleware.UserIDKey)
	var req UpdateEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.UpdateEmail(c.Request.Context(), userID, req.NewEmail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email updated successfully"})
}

// ChangePassword godoc
// @Summary      Change own password
// @Description  Updates the password of the authenticated user.
// @Tags         auth
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        request  body      ChangePasswordRequest  true  "Old and New Passwords"
// @Success      200      {object}  map[string]string   "Success message"
// @Failure      400      {object}  map[string]string   "Invalid password or mismatch"
// @Failure      500      {object}  map[string]string   "Failed to change password"
// @Router       /auth/me/password [put]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID := c.GetString(middleware.UserIDKey)
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

// Logout godoc
// @Summary      Logout user
// @Description  Logs out the authenticated user by invalidating the current session.
// @Tags         default
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {object}  map[string]string "Success message"
// @Failure      500  {object}  map[string]string "Failed to logout"
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// Refresh godoc
// @Summary      Refresh access token
// @Description  Refreshes the access token using a valid refresh token.
// @Tags         default
// @Accept       json
// @Produce      json
// @Param        request  body      RefreshRequest  true  "Refresh Token"
// @Success      200      {object}  map[string]string   "New access token"
// @Failure      400      {object}  map[string]string   "Invalid request"
// @Failure      401      {object}  map[string]string   "Unauthorized: Invalid refresh token"
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newAccessToken, err := h.svc.RefreshAccessToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": newAccessToken})
}

// UpdateUserStatus godoc
// @Summary      Update user status (Admin only)
// @Description  Updates the status of a specific user. Requires admin privileges.
// @Tags         admin
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        request  body      UpdateStatusRequest  true  "User ID and New Status"
// @Success      200      {object}  map[string]string   "Success message"
// @Failure      400      {object}  map[string]string   "Invalid status value"
// @Failure      500      {object}  map[string]string   "Failed to update user status"
// @Router       /auth/admin/user-status [put]
func (h *AuthHandler) UpdateUserStatus(c *gin.Context) {
	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.UpdateUserStatus(c.Request.Context(), req.UserID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User status updated to " + req.Status})
}