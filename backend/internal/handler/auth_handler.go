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
//@Accept       json
//@Produce      json
//@Param        request body      ForgotPasswordRequest  true  "Forgot Password Info"
//@Success      200     {string}  string "If the email exists, a reset code has been sent."
// @Failure      400     {string}  string "Invalid Request"
// @Failure      500     {string}  string "Internal Server Error"
//@Router       /auth/forgot-password [post]
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

//GetProfile godoc
// @Summary      Get user profile
// @Description  Retrieves the profile information of the authenticated user.
// @Tags         auth
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200     {object}  service.UserProfile
// @Failure      401     {string}  string "Unauthorized"
// @Router       /auth/me [get]
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request){
	// 1. Pull UserID out of the context (set by middleware)
	userID, ok := r.Context().Value(middleware.UserIDkey).(string)
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