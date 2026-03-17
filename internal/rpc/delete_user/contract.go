package delete_user

import (
	"context"
	"myCalendar/internal/domain"
)

type repo interface {
	GetByUsername(ctx context.Context, username string) (domain.User, error)
	Delete(ctx context.Context, username string) error
}
