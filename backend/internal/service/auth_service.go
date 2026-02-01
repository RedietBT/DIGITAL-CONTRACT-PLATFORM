package service

import (
	"context"
	"errors"

	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/models"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

//AuthService defines the business logic for user authentication and registration.
type AuthService interface{
	Register(ctx context.Context, email, password string) (*models.User, error)
}

type authService struct {
	repo repository.UserRepository
}

//NewAuthService creates a new instance of AuthService with the provided UserRepository.
//This is our "Constructor" for Dependency Injection.
func NewAuthService(repo repository.UserRepository) AuthService{
	return  &authService{repo: repo}
}

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

