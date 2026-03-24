package get_user_by_login

import (
	"context"
	"myCalendar/internal/domain"
)

type repo interface {
	GetByLogin(ctx context.Context, login string) (domain.User, error)
}
