package update_user

import (
	"context"
	"myCalendar/internal/domain"
)

type repo interface {
	GetByUsername(ctx context.Context, username string) (domain.User, error)
	Update(ctx context.Context, username, email, name, surname, birth string) (domain.User, error)
}
