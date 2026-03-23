package delete_event

import (
	"context"
	"myCalendar/internal/domain"
)

type repo interface {
	GetByID(ctx context.Context, id string) (domain.Event, error)
	Delete(ctx context.Context, id string) error
}
