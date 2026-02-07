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
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/broker"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

//AuthService defines the business logic for user authentication and registration.
type AuthService interface{
	Register(ctx context.Context, email, password string) (*models.User, error)
	Login(ctx context.Context, email, password string) (string, string, error)
	ForgotPassword(ctx context.Context, emailAddr string) (error)
	ResetPassword(ctx context.Context, email, token, newPassword string) (error)
	GetUserByID(ctx context.Context, UserID string) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	DeleteUser(ctx context.Context, UserID string) (error) 
	UpdateEmail(ctx context.Context, UserID string, newEmail string) (error)
	ChangePassword(ctx context.Context, userID, oldpassword, newpassword string) (error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (string, error)
	UpdateUserStatus(ctx context.Context, userID string, newStatus string) (error)
}

type authService struct {
	repo repository.UserRepository
	broker *broker.RabbitMQProvider
	jwtSecret string
}

//NewAuthService creates a new instance of AuthService with the provided UserRepository.
//This is our "Constructor" for Dependency Injection.
func NewAuthService(repo repository.UserRepository, broker *broker.RabbitMQProvider, secret string) AuthService{
	return  &authService{
		repo: repo,
		broker: broker,
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

	// ✨ ADD THIS: Tell RabbitMQ a user was created!
    if s.broker != nil {
        // We broadcast the ID and Email so the Profile Service can store them
        _ = s.broker.PublishUserCreated(ctx, user.ID, user.Email) 
    }

	return user, nil
}

//Login authenticates a user and returns a JWT token upon successful login.
func(s *authService) Login(ctx context.Context, email, password string) (string, string, error){
	//1. Get user from DB
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil{
		return "", "", errors.New("Invalid Credentials")
	}

	// Check account status
	if user.Status != "active" {
		return "", "", errors.New("account is" + user.Status + ". Please contact support")
	}

	//2. Compare Bcrypt hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil{
		return "", "", errors.New("Invalid Credentials")
	}

	//Update last login time (Asynchronous or ignore error)
	_ = s.repo.UpdateLastLogin(ctx, user.ID)

	//3. Generate JWT BOTH tokens
	accessToken, err := s.generateToken(user)
	if err != nil{
		return "", "", err
	}

	// 4. Generate Refresh Token
	refreshToken := s.generateRandomString(32)

	// 5. SAVE the refresh token to the database
	// We set it to expire in 7 days
	expiry := time.Now().Add(7 * 24 * time.Hour)
	err = s.repo.SaveRefreshToken(ctx, user.ID, refreshToken, expiry)
	if err != nil {
		return "", "", errors.New("Failed to create sessions")
	}

	return accessToken, refreshToken, nil
}

// Helper method to generate the refresh token
func (s *authService) generateRandomString(n int) string{
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return ""
	}

	return hex.EncodeToString(b)
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

	//Get the User by Email to find their ID
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil{
		return errors.New("user no longer exists")
	}

	// 4. Hash the new password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

	// 5. Update user password in DB (You might need to add UpdatePassword to your repo)
    err = s.repo.UpdatePassword(ctx, user.ID, string(hashedPassword))
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

func (s *authService) GetAllUsers(ctx context.Context) ([]*models.User, error){
	users, err := s.repo.GetAllUsers(ctx)
	if err != nil{
		return nil, err
	}
	return users, nil
}

func (s *authService) DeleteUser(ctx context.Context, UserID string) error {
	// Checke if user exsists
	_, err := s.repo.GetByID(ctx, UserID)
	if err != nil{
		return errors.New("User Not Found")
	}

	// Delete the user
	err = s.repo.DeleteUser(ctx, UserID)
	if err != nil{
		return errors.New("Failed to Delet user from database")
	}

	return nil
}

func (s *authService) UpdateEmail(ctx context.Context, UserID string, newEmail string) error{
	
	// 1. Update Users Email
	existingUser, err := s.repo.GetByEmail(ctx, newEmail)
	if err == nil && existingUser != nil{
		// If weq found a user with this email, and its's not the current user
		if existingUser.ID != UserID {
			return errors.New("Failed to Update user Email")
		}	
	}

	// 2. Call the repository to update
	return s.repo.UpdateEmail(ctx, UserID, newEmail)
}

func (s *authService) ChangePassword(ctx context.Context, userID, oldpassword, newpassword string) error {
	// 1. Fetch the user to get the CURRENT hash
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	// 2. Compare the "oldPassword" provided by the user with the hash in DB
	// We reuse existing Bcrypt comparison logic
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldpassword))
	if err != nil {
		return errors.New("current password is incorrect")
	}

	//Hash the NEW password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newpassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 4. Update the DB
	return s.repo.UpdatePassword(ctx, userID, string(hashedPassword))
}

//Refresh acces Token
func(s *authService) RefreshAccessToken(ctx context.Context, refreshToken string) (string, error){
	// 1. Get the full model from Db
	rt, err := s.repo.GetRefreshToken(ctx, refreshToken)
	if err != nil{
		return "", errors.New("invalid refresh token")
	}

	// 2. Check if token is expired
	if time.Now().After(rt.ExpiresAt){
		// Clean up the database
		_= s.repo.DeleteRefreshToken(ctx, refreshToken)
		return "", errors.New("refresh token expired")
	}

	// 3. Get user details to create a new JWT
	user, err := s.repo.GetByID(ctx, rt.UserID)

	// 4. Check if user is still active
	if user.Status != "active" {
		return "", errors.New("account is no longer activce contact help")
	}

	// 5. Generate and return the new short_lived Access Token
	return s.generateToken(user)
}

// Update user Status
func (s *authService) UpdateUserStatus(ctx context.Context, userID string, newStatus string) error {
	// 1. Validate status types
	validStatuses := map[string]bool{"active": true, "suspended": true, "deactivated": true}
	if !validStatuses[newStatus]{
		return errors.New("invalid status type")
	}

	// 2. Call repository
	return s.repo.UpdateUserStatus(ctx, userID, newStatus)
}