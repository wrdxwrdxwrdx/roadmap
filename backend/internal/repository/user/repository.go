package user

import (
	"context"

	userentity "roadmap/internal/domain/entities/user"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *userentity.User) (*userentity.User, error)

	GetByID(ctx context.Context, id uuid.UUID) (*userentity.User, error)

	GetByEmail(ctx context.Context, email string) (*userentity.User, error)

	EmailExists(ctx context.Context, email string) (bool, error)

	UsernameExists(ctx context.Context, username string) (bool, error)
}
