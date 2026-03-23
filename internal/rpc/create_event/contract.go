package create_event

import (
	"context"
	"myCalendar/internal/domain"
)

type repo interface {
	Create(ctx context.Context, e domain.Event) error
}
