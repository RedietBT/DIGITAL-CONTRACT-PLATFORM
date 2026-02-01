package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/models"
)

//UserRepository defines what actions we can take on the User table
type UserRepository interface{
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	SaveResetToken(ctx context.Context, email, token string, expiresAt time.Time) error
	GetResetToken(ctx context.Context, email string) (string, time.Time, error)
	DeleteResetToken(ctx context.Context, email string) (error)
	UpdatePassword(ctx context.Context, email, hashedPassword string) error
}

type postgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) UserRepository{
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) Create(ctx context.Context, user *models.User) error{
	query := `
		INSERT INTO auth_schema.users (email, password_hash, role, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

		return  r.db.QueryRowContext(ctx, query,
		    user.Email,
			user.PasswordHash,
			user.Role,
			user.Status,
		).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *postgresUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error){
	query := `
		SELECT id, email, password_hash, role, status, created_at, updated_at, last_login_at
		FROM auth_schema.users
		WHERE email = $1`

		user := &models.User{}
		err := r.db.QueryRowContext(ctx, query, email).Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.Role,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.LastLoginAt,
		)

		if err !=nil{
			return  nil, err
		}
		return user, nil
}

func (r *postgresUserRepository) SaveResetToken(ctx context.Context, email, token string, expiresAt time.Time) error{
	query := `
		INSERT INTO auth_schema.password_resets (email, token, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (email) DO UPDATE SET token = $2, expires_at = $3`
	_, err := r.db.ExecContext(ctx, query, email, token, expiresAt)
	return err
}

func (r *postgresUserRepository) GetResetToken(ctx context.Context, email string) (string, time.Time, error){
	query := `
		SELECT token, expires_at
		FROM auth_schema.password_resets
		WHERE email = $1`
	var token string
	var expiresAt time.Time
	err := r.db.QueryRowContext(ctx, query, email).Scan(&token, &expiresAt)
	if err != nil{
		return "", time.Time{}, err
	}
	return token, expiresAt, nil
}

func (r *postgresUserRepository) DeleteResetToken(ctx context.Context, email string) error{
	query := `
		DELETE FROM auth_schema.password_resets
		WHERE email = $1`
	_, err := r.db.ExecContext(ctx, query, email)
	return err
}

func (r *postgresUserRepository) UpdatePassword(ctx context.Context, email, hashedPassword string) error {
    query := `UPDATE auth_schema.users SET password_hash = $1 WHERE email = $2`
    _, err := r.db.ExecContext(ctx, query, hashedPassword, email)
    return err
}