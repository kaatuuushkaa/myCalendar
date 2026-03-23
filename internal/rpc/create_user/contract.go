package create_user

import (
	"context"
	"myCalendar/internal/domain"
)

type repo interface {
	Create(ctx context.Context, u domain.User) error
}
