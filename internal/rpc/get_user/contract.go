package get_user

import (
	"context"
	"myCalendar/internal/domain"
)

type repo interface {
	GetByUsername(ctx context.Context, username string) (domain.User, error)
}
