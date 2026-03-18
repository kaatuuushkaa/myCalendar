package pgrepo

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"myCalendar/internal/apperrors"
	"myCalendar/internal/domain"
	"time"
)

type EventRepo struct {
	db *gorm.DB
}

func NewEventRepo(db *gorm.DB) *EventRepo {
	return &EventRepo{db: db}
}

func (r *EventRepo) Create(ctx context.Context, e domain.Event) error {
	if err := r.db.WithContext(ctx).Create(e).Error; err != nil {
		return err
	}
	return nil
}

func (r *EventRepo) GetByID(ctx context.Context, id string) (domain.Event, error) {
	var e domain.Event
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&e).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Event{}, apperrors.ErrEventNotFound
		}
		return domain.Event{}, apperrors.ErrInternal
	}
	return e, nil
}

func (r *EventRepo) GetByUserID(ctx context.Context, userID string) ([]domain.Event, error) {
	var e []domain.Event
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("start_at ASC").First(&e).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrEventNotFound
		}
		return nil, apperrors.ErrInternal
	}
	return e, nil
}

func (r *EventRepo) Update(ctx context.Context, id, title, description, eventDate string, startAt, endAt time.Time) (domain.Event, error) {
	var e domain.Event
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&e).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Event{}, apperrors.ErrEventNotFound
		}
		return domain.Event{}, apperrors.ErrInternal
	}

	e.Title = title
	e.Description = description
	e.StartAt = startAt
	e.EndAt = endAt
	e.EventDate = eventDate

	if err = r.db.WithContext(ctx).Save(&e).Error; err != nil {
		return domain.Event{}, apperrors.ErrInternal
	}

	return e, nil
}

func (r *EventRepo) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&domain.Event{})
	if result.Error != nil {
		return apperrors.ErrInternal
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrEventNotFound
	}
	return nil
}
