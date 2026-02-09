package repository

import (
	"context"
	"database/sql"

	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/models"
)

// ProfileRepository defines the database operations for user profiles.
type ProfileRepository interface {
	CreateProfile(ctx context.Context, profile *models.Profile) error
	GetProfileByID(ctx context.Context, userID string) (*models.Profile, error)
	UpdateProfile(ctx context.Context, profile *models.Profile) error
	DeleteProfile(ctx context.Context, userID string) error
}

type profileRepository struct {
	db *sql.DB
}

// NewProfileRepository creates a new instance of the profile repository.
func NewProfileRepository(db *sql.DB) ProfileRepository {
	return &profileRepository{db: db}
}

// CreateProfile inserts a new profile record. 
// We use ON CONFLICT because RabbitMQ might retry and send the same message twice.
func (r *profileRepository) CreateProfile(ctx context.Context, p *models.Profile) error {
	query := `
	    INSERT INTO profile_schema.profiles (
		    user_id, display_name, bio, skill_level, is_template_seller
		) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id) DO NOTHING;`

		_, err := r.db.ExecContext(ctx, query,
		          p.UserID,
		          p.DisplayName,
		          p.Bio,
		          p.SkillLevel,
		          p.IsTemplateSeller,
		)
		return err
}

// GetProfileByID retrieves a profile by the user's UUID.
func (r *profileRepository) GetProfileByID(ctx context.Context, userID string) (*models.Profile, error){
	p := &models.Profile{}
	query := `SELECT user_id, display_name, bio, skill_level, is_template_seller, rating_avg, rating_count, created_at, updated_at 
		FROM profile_schema.profiles 
		WHERE user_id = $1`

		 err := r.db.QueryRowContext(ctx, query, userID).Scan(
			&p.UserID,
			&p.DisplayName,
			&p.Bio,
			&p.SkillLevel,
			&p.IsTemplateSeller,
			&p.RatingAvg,
			&p.RatingCount,
			&p.CreatedAt,
			&p.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		return p, nil
}

// UpdateProfile updates an existing profile details.
// Note: We don't update update_at here because it's automatically updated by the database. 
func (r *profileRepository) UpdateProfile(ctx context.Context, p *models.Profile) error {
	query := `
		UPDATE profile_schema.profiles
		SET display_name = $1, bio = $2, skill_level = $3, is_template_seller = $4
		WHERE user_id = $5
	`

	_, err := r.db.ExecContext(ctx, query,
		p.DisplayName,
		p.DisplayName,
		p.Bio,
		p.SkillLevel,
		p.IsTemplateSeller,
		p.UserID,
	)
	return err
}

// DeleteProfile removes profile.
// Note: If the Auth User is deleted, the profile is deleted automatically.
func (r *profileRepository) DeleteProfile(ctx context.Context, userID string) error {
	query := `
		DELETE FROM profile_schema.profiles
		WHERE user_id = $1
	`

	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}