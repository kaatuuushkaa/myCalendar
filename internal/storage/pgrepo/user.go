package pgrepo

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"myCalendar/internal/apperrors"
	"myCalendar/internal/domain"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, u domain.User) error {
	if err := r.db.WithContext(ctx).Create(u).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	var u domain.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, apperrors.ErrUserNotFound
		}
		return domain.User{}, apperrors.ErrInternal
	}
	return u, nil
}

func (r *UserRepo) GetByLogin(ctx context.Context, login string) (domain.User, error) {
	var u domain.User
	err := r.db.WithContext(ctx).Where("username = ? OR email = ?", login, login).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, apperrors.ErrUserNotFound
		}
		return domain.User{}, apperrors.ErrInternal
	}
	return u, nil
}

func (r *UserRepo) Update(ctx context.Context, username, email, name, surname, birth string) (domain.User, error) {
	var u domain.User
	err := r.db.WithContext(ctx).
		Where("username = ?", username).
		First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, apperrors.ErrUserNotFound
		}
		return domain.User{}, apperrors.ErrInternal
	}

	u.Email = email
	u.Name = name
	u.Surname = surname
	u.Birth = birth

	if err = r.db.WithContext(ctx).Save(&u).Error; err != nil {
		return domain.User{}, apperrors.ErrInternal
	}

	return u, nil
}

func (r *UserRepo) Delete(ctx context.Context, username string) error {
	result := r.db.WithContext(ctx).
		Where("username = ?", username).
		Delete(&domain.User{})
	if result.Error != nil {
		return apperrors.ErrInternal
	}
	// RowsAffected == 0 означает что такого пользователя не было
	if result.RowsAffected == 0 {
		return apperrors.ErrUserNotFound
	}
	return nil
}
