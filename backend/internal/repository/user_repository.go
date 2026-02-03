package repository

import (
	"context"
	"database/sql"
	"errors"
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
	UpdatePassword(ctx context.Context, email, hashedPassword string) (error)
	GetByID(ctx context.Context, UserID string) (*models.User, error)
	UpdateLastLogin(ctx context.Context, userID string) (error)
	GetAllUsers(ctx context.Context) ([]*models.User ,error)
	DeleteUser(ctx context.Context, UserID string) (error)

}

//postgresUserRepository is the Postgres implementation of UserRepository
type postgresUserRepository struct {
	db *sql.DB
}

//NewPostgresUserRepository creates a new instance of PostgresUserRepository
func NewPostgresUserRepository(db *sql.DB) UserRepository{
	return &postgresUserRepository{db: db}
}

//Create inserts a new user into the database
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

//GetByEmail retrieves a user by their email address
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

//SaveResetToken saves or updates a password reset token for a given email
func (r *postgresUserRepository) SaveResetToken(ctx context.Context, email, token string, expiresAt time.Time) error{
	query := `
		INSERT INTO auth_schema.password_resets (email, token, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (email) DO UPDATE SET token = $2, expires_at = $3`
	_, err := r.db.ExecContext(ctx, query, email, token, expiresAt)
	return err
}

//GetResetToken retrieves the reset token and its expiration time for a given email
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

//DeleteResetToken deletes the reset token for a given email
func (r *postgresUserRepository) DeleteResetToken(ctx context.Context, email string) error{
	query := `
		DELETE FROM auth_schema.password_resets
		WHERE email = $1`
	_, err := r.db.ExecContext(ctx, query, email)
	return err
}

//UpdatePassword updates the password hash for a given email
func (r *postgresUserRepository) UpdatePassword(ctx context.Context, email, hashedPassword string) error {
    query := `UPDATE auth_schema.users SET password_hash = $1 WHERE email = $2`
    _, err := r.db.ExecContext(ctx, query, hashedPassword, email)
    return err
}

//GetByID retrieves a user by their unique ID
func (r *postgresUserRepository) GetByID(ctx context.Context, UserID string) (*models.User, error){
	query := `
		SELECT id, email, password_hash, role, status, created_at, updated_at, last_login_at
		FROM auth_schema.users
		WHERE id = $1`

		var user models.User

		//We use QueryRowContext because we only expect ONE user back
		err := r.db.QueryRowContext(ctx, query, UserID).Scan(
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
		return &user, nil
}

// UpdateLastLogin 
func (r *postgresUserRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	query := `UPDATE auth_schema.users SET last_login_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

// GetAllUsers
func (r *postgresUserRepository) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	// 1. Define the qurey
	query := `SELECT id, email, role, status, created_at, updated_at, last_login_at 
        FROM auth_schema.users 
        ORDER BY created_at DESC`

		// 2. Use QureyContext to get rows
		rows, err := r.db.QueryContext(ctx, query)
		if err != nil{
			return nil, err
		}

		defer rows.Close() // Always close rows to prevent memory leaks!

		var users []*models.User

		// 3. Iterate through the rows
		for rows.Next(){
			var u models.User
			// Scan each column into the struct fields
			err := rows.Scan(
				&u.ID,
				&u.Email,
				&u.Role,
				&u.Status,
				&u.CreatedAt,
				&u.UpdatedAt,
				&u.LastLoginAt,
			)

			if err != nil{
				return nil, err
			}

			// Add the user to our slice
			users = append(users, &u)
		}

		if err = rows.Err(); err != nil{
			return nil, err
		}

		return users, nil 
}

//DeleteUser using userID
func (r *postgresUserRepository) DeleteUser(ctx context.Context, UserID string) error{
	query := `DELETE FROM auth_schema.users WHERE id=$1`
	result, err := r.db.ExecContext(ctx, query, UserID)
	if err != nil{
		return err
	} 

	//Check if a row was actually deleted
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no user with that ID")
	}

	return nil
}