package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/models"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/repository"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/pkg/email"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

//AuthService defines the business logic for user authentication and registration.
type AuthService interface{
	Register(ctx context.Context, email, password string) (*models.User, error)
	Login(ctx context.Context, email, password string) (string, error)
	ForgotPassword(ctx context.Context, emailAddr string) (error)
	ResetPassword(ctx context.Context, email, token, newPassword string) (error)
	GetUserByID(ctx context.Context, UserID string) (*models.User, error)
}

type authService struct {
	repo repository.UserRepository
	jwtSecret string
}

//NewAuthService creates a new instance of AuthService with the provided UserRepository.
//This is our "Constructor" for Dependency Injection.
func NewAuthService(repo repository.UserRepository, secret string) AuthService{
	return  &authService{
		repo: repo,
		jwtSecret: secret,
	}
}

//Register registers a new user with the provided email and password.
func(s *authService) Register(ctx context.Context, email, password string)(*models.User, error){

	//1. Check if user with the email already exists
	existingUser, _ := s.repo.GetByEmail(ctx, email)
	if existingUser != nil{
		return  nil, errors.New("user with this email already exists")
	}

	//2. GenerateFromPassword handles salting and hashing password in one go.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		return nil, err
	}

	user := &models.User{
		Email:    email,
		PasswordHash: string(hashedPassword),
		Role:     "user",
		Status:   "active",
	}

	err = s.repo.Create(ctx, user)
	return user, err
}

//Login authenticates a user and returns a JWT token upon successful login.
func(s *authService) Login(ctx context.Context, email, password string) (string, error){
	//1. Get user from DB
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil{
		return "", errors.New("Invalid Credentials")
	}

	//2. Compare Bcrypt hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil{
		return "", errors.New("Invalid Credentials")
	}

	//Update last login time
	s.repo.UpdateLastLogin(ctx, user.ID)

	//3. Generate JWT Token
	return s.generateToken(user)
}

//generateToken creates a JWT token for the authenticated user.
func (s *authService) generateToken(user *models.User) (string, error){

	//1. Create the Claims (the data we want to "hide" inside the token)
	claims := jwt.MapClaims{
		"sub": user.ID,
		"role": user.Role,
		"exp": time.Now().Add(time.Hour * 24).Unix(), //Token expires in 24 hours
		"iat": time.Now().Unix(),
	}

	//2. Create the tocken object usiing HS256 algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//3. Sign the token with our secret key
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *authService) ForgotPassword(ctx context.Context, emailAddr string) error {
	//1. Check if the user exists (don't tell the user if they don't exist for security reasons)
	_, err := s.repo.GetByEmail(ctx, emailAddr)
	if err != nil{
		return nil// We return nil to prevent "Email Enumeration" attacks
	}

	//2. Generate a random secur token
	b := make([]byte, 4)// 8 Characters total
	rand.Read(b)
	token := hex.EncodeToString(b)

	//3. Save to DB with 15-minute expiry
	expiry := time.Now().Add(15 * time.Minute)
	err = s.repo.SaveResetToken(ctx, emailAddr, token, expiry)
	if err != nil{
		return err
	}

	//4. Send the actual email via our utiltiy
	return email.SendResetEmail(emailAddr, token)
}

func (s *authService) ResetPassword(ctx context.Context, email, token, newPassword string) error {

	//1. Get the token from DB
	storedToken, expiryAt, err := s.repo.GetResetToken(ctx, email)
	if err != nil{
		return errors.New("invalid or expired token")
	}

	// 2. Check if expired
    if time.Now().After(expiryAt) {
        s.repo.DeleteResetToken(ctx, email) // Clean up expired token
        return errors.New("token has expired")
    }

	// 3. Verify token match
    if storedToken != token {
        return errors.New("invalid token")
    }

	// 4. Hash the new password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

	// 5. Update user password in DB (You might need to add UpdatePassword to your repo)
    err = s.repo.UpdatePassword(ctx, email, string(hashedPassword))
    if err != nil {
        return err
    }

	// 6. Delete the token so it can't be used again
    return s.repo.DeleteResetToken(ctx, email)
}

func (s *authService) GetUserByID(ctx context.Context, UserID string) (*models.User, error){
	user, err := s.repo.GetByID(ctx, UserID)
	if err != nil{
		return nil, err
	}
	return user, nil
}
