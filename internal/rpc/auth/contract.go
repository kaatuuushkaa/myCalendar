package auth

import (
	"context"
	"myCalendar/internal/domain"
)

type repo interface {
	GetByLogin(ctx context.Context, login string) (domain.User, error)
}
