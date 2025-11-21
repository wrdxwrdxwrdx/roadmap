package user

import (
	"context"
	"fmt"

	userentity "roadmap/internal/domain/entities/user"
	"roadmap/internal/infrastructure/database"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type userRepository struct {
	db *database.Database
}

func NewUserRepository(db *database.Database) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *userentity.User) (*userentity.User, error) {
	query := `
		INSERT INTO users (id, email, password_hash, username, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, email, password_hash, username, created_at, updated_at
	`

	var createdUser userentity.User
	err := r.db.Pool.QueryRow(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Username,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(
		&createdUser.ID,
		&createdUser.Email,
		&createdUser.PasswordHash,
		&createdUser.Username,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &createdUser, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*userentity.User, error) {
	query := `
		SELECT id, email, password_hash, username, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user userentity.User
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Username,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*userentity.User, error) {
	query := `
		SELECT id, email, password_hash, username, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user userentity.User
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Username,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func (r *userRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}

func (r *userRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`

	var exists bool
	err := r.db.Pool.QueryRow(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}

	return exists, nil
}
