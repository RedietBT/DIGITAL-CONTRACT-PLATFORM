package handler

import (
	"encoding/json"
	"net/http"

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
	Email string `json: "email" validate:"required,email" example:"dev@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"secret123"`
}

//Refiater godoc
//@Summary Register a new user
//@Description Create a new user account and returns the user abject(excluding password).
//@ Tags auth
//@Accept json
//@Produce json
//@Param request body RegisterRequest true "Registration Info"
//@Success 201 {object} model.User
//@Failure 400 {string} string "Invalide Request"
//@Failure 500 {string} string "Server Error"
//@Router /auth/register [post]
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