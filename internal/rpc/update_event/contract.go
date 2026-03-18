package update_event

import (
	"context"
	"myCalendar/internal/domain"
	"time"
)

type repo interface {
	GetByID(ctx context.Context, id string) (domain.Event, error)
	Update(ctx context.Context, uuid, title, description, eventDate string, startAt, endAt time.Time) (domain.Event, error)
}
