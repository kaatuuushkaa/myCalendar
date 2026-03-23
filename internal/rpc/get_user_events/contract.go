package get_user_events

import (
	"context"
	"myCalendar/internal/domain"
)

type repo interface {
	GetByUserID(ctx context.Context, userID string) ([]domain.Event, error)
}
