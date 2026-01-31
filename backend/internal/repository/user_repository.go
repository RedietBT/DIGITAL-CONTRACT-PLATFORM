package repository

import (
	"context"
	"database/sql"

	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/models"
)

//UserRepository defines what actions we can take on the User table
type UserRepository interface{
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
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