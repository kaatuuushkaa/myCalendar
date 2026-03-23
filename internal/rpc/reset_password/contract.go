package reset_password

import (
	"context"
	"myCalendar/internal/domain"
)

type repo interface {
	GetByUsername(ctx context.Context, username string) (domain.User, error)
	ResetPassword(ctx context.Context, username, password string) error
}
