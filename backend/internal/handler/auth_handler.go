package handler

import (
	"encoding/json"
	"net/http"

	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/middleware"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/service"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

//AuthHandler coordinations HTTP requests from authentication.
type AuthHandler struct{
	svc service.AuthService
}

//NewAuthHandler initalizes the handler with the required service.
func NewAuthHandler(svc service.AuthService) *AuthHandler{
	return &AuthHandler{svc: svc}
}

//RegisterRequest represents the JSON body from registration endpoint.
type RegisterRequest struct{
	// 'validate' tags define our rules
	//required: must be present
	//email: must be a valid email format
	//min=8: password must be at least 8 characters
	Email string `json:"email" validate:"required,email" example:"dev@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"secret123"`
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account and returns the user object (excluding password).
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body      RegisterRequest  true  "Registration Info"
// @Success      201     {object}  models.User
// @Failure      400     {string}  string "Invalid Request"
// @Failure      500     {string}  string "Server Error"
// @Router       /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request){
	var req RegisterRequest

	//1. Decode the incoming JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	//2. Validate Struct (Content check)
	if err := validate.Struct(req); err != nil{
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	//3. Execute business logic
	user, err := h.svc.Register(r.Context(), req.Email, req.Password)
	if err != nil{
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	//Respond with the created user(excluding password)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)

}

//LoginRequest represents the input data from login endpoint.
type LoginRequest struct{
	Email string `json:"email" validate:"required,email" example:"dev@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"secret123"`
}

//LoginResponse represents the output data( the token)
type LoginResponse struct{
	Token string `json:"token"`
}

//Login godoc
// @Summary      User login
// @Description  Authenticate user and returns a JWT token upon successful login.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body      LoginRequest  true  "Login Info"
// @Success      200     {object}  LoginResponse
// @Failure      400     {string}  string "Invalid Request"
// @Failure      401     {string}  string "Unauthorized"
// @Failure      500     {string}  string "Server Error"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request){
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	//Validate the Struct
	if err := validate.Struct(req); err != nil{
		http.Error(w, "Invalid email or password format", http.StatusBadRequest)
		return
	}

	token, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil{
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})

}

type ForgotPasswordRequest struct{
	Email string `json:"email" validate:"required,email" example:"dev@example.com"`
}

//ForgotPassword godoc
// @Summary      Forgot Password
// @Description  Initiate password reset process by sending a reset link to the user's email.
// @Tags         auth
//@Accept        json
//@Produce       json
//@Param         request body      ForgotPasswordRequest  true  "Forgot Password Info"
//@Success       200     {string}  string "If the email exists, a reset code has been sent."
// @Failure      400     {string}  string "Invalid Request"
// @Failure      500     {string}  string "Internal Server Error"
//@Router        /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request){
	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	//Logic call
	_ = h.svc.ForgotPassword(r.Context(), req.Email)

	// We always return 200 so hackers can't "fish" for valid emails
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("If that email is in our system, a reset code has been sent."))
}

type ResetPasswordRequest struct {
    Email       string `json:"email" validate:"required,email"`
    Token       string `json:"token" validate:"required"`
    NewPassword string `json:"new_password" validate:"required,min=8"`
}

//ResetPassword godoc
// ResetPassword godoc
// @Summary      Reset password using token
// @Description  Verifies the reset token and updates the user's password in the database.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body      ResetPasswordRequest  true  "Reset Details"
// @Success      200     {string}  string "Password reset successful"
// @Failure      400     {string}  string "Invalid token or expired"
// @Failure      500     {string}  string "Internal Server Error"
// @Router       /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate the struct
	if err := validate.Struct(req); err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	//Call service logic 
	err := h.svc.ResetPassword(r.Context() , req.Email, req.Token, req.NewPassword)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("Password has been reset successfully. You can now login."))
}

// GetProfile godoc
// @Summary      Get user profile
// @Description  Retrieves the profile information of the authenticated user.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200 {object} models.User
// @Failure      401 {string} string "Unauthorized"
// @Failure      404 {string} string "User not found"
// @Router       /auth/me [get]
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request){
	// 1. Pull UserID out of the context (set by middleware)
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Could not find user in context", http.StatusUnauthorized)
		return
	}

	// 2. Featch user from DB using the ID
	user, err := h.svc.GetUserByID(r.Context(), userID)
	if err != nil{
		http.Error(w, "User not Found", http.StatusNotFound)
		return
	}

	// 3. Return the user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// GetAllUsers godoc
// @Summary      Get all users (Admin only)
// @Description  Retrieves a list of all registered users. Accessible only by admin users.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200 {array} models.User
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "Forbidden"
// @Router       /auth/admin/users [get]
func (h *AuthHandler) GetAllUsers(w http.ResponseWriter, r *http.Request){
	// 1. Call the service
	// Note: Since this is behind RoleMiddleware, we KNOW the user is an admin
	users, err := h.svc.GetAllUsers(r.Context())
	if err != nil{
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	// Set the header to JSON
	w.Header().Set("Content-Type", "application/json")

	// 3. Encode the slice of users into response
	if err := json.NewEncoder(w).Encode(users); err != nil{
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}


// DeleteUser godoc
// @Summary      Delete a user
// @Description  Removes a user from the system. Admin only.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   query      string  true  "User ID"
// @Success      204  {string}   string  "No Content"
// @Failure      400  {string}   string  "ID required"
// @Router       /auth/admin/users [delete]
func (h *AuthHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// 1. Get the ID of the person LOOGED IN (from JWT/Context)
	currentUserID,_ := r.Context().Value(middleware.UserIDKey).(string)

	// 2. Get the ID from the URL (if provided)
	targetID := r.URL.Query().Get("id")

	var idToDelete string

	if targetID != ""{
		// If an ID is in the URL, the user is trying to delete someone else.
        // Our RoleMiddleware (in main.go) ensures only Admins reach this logic.
		idToDelete = targetID
	} else {
		// If not ID in URL, the user is tring to delete THEMSELVES.
		idToDelete =currentUserID
	}

	// 3. Call service
	if err := h.svc.DeleteUser(r.Context(), idToDelete); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteMe godoc
// @Summary      Delete own account
// @Description  Deletes the currently authenticated user's account.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   query      string  true  "User ID"
// @Success      204  {string}  string  "No Content"
// @Failure      400  {string}   string  "ID required"
// @Router       /auth/me [delete]
func (h *AuthHandler) DeleteMe(w http.ResponseWriter, r *http.Request){
	h.DeleteUser(w, r)
}

type UpdateEmailRequest struct{
		NewEmail string `json:"new_email" validate:"required,email"`
	}

// UpdateEmail godoc
// @Summary      Update own email
// @Description  Updates the email address of the currently authenticated user.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request  body      UpdateEmailRequest  true  "New Email Details"
// @Success      200      {object}  map[string]string   "{"message": "Email updated successfully"}"
// @Failure      400      {string}  string              "Invalid request body"
// @Failure      401      {string}  string              "Unauthorized"
// @Router       /auth/me/email [put]
func (h *AuthHandler) UpdateEmail(w http.ResponseWriter, r *http.Request){
	// 1. Pull UserID out of the context (set by middleware)
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Could not find user in context", http.StatusUnauthorized)
		return
	}

	// 2. Decode the request body
	var req UpdateEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 3. Call the service
	if err := h.svc.UpdateEmail(r.Context(), userID, req.NewEmail); err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Email updaated successfully"})
}

type ChangePasswordRequest struct{
	OldPassword string `json:"old_password" validate:"required"`
    NewPassword string `json:"new_password" validate:"required,min=8"`
}

// ChangePassword godoc
// @Summary      Change own password
// @Description  Allows the user to update their password after verifying their old one.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request  body      ChangePasswordRequest  true  "Password Details"
// @Success      200      {object}  map[string]string      "Password changed successfully"
// @Failure      400      {string}  string                 "Incorrect current password"
// @Router       /auth/me/password [put]
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// 1. Pull UserID out of the context (set by middleware)
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return 
	}

	if err := h.svc.ChangePassword(r.Context(), userID, req.OldPassword, req.NewPassword); err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password updated successfully"})
	
}