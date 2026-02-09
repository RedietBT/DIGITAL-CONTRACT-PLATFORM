package service

import (
	"context"
	"errors"

	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/models"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/repository"
)

type ProfileService interface{
	// Use the Event struct instead of individual strings
	HandleUserCreatedEvent(ctx context.Context, event models.UserCreatedEvent) error
}

type profileService struct {
	repo repository.ProfileRepository
}

func NewProfileService(repo repository.ProfileRepository) ProfileService {
	return &profileService{repo: repo}
}

// HandleUserCreatedEvent is triggered by RabbitMQ to ensure a profile exists for every user.
func (s *profileService) HandleUserCreatedEvent(ctx context.Context, event models.UserCreatedEvent) error {
	newProfile := &models.Profile{
		UserID: event.UserID,
		DisplayName: event.Email,
		Bio : "Welcome to my Profile!",
		SkillLevel: 1,
		IsTemplateSeller: false,
	}

	return s.repo.CreateProfile(ctx, newProfile)
}

func (s *profileService) GetProfile(ctx context.Context, userID string) (*models.Profile, error) {
	profile, err := s.repo.GetProfileByID(ctx, userID)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	return profile, nil
}

func (s *profileService) UpdateProfile(ctx context.Context, p *models.Profile) error {
	// Prevent skill level from going below 1
	if p.SkillLevel < 1 {
		p.SkillLevel = 1
	}

	return s.repo.UpdateProfile(ctx, p)
}