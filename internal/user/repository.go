package user

import (
	"context"
	"fmt"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

type RepositoryInterface interface {
	Create(ctx context.Context, u *User) error
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByLogin(ctx context.Context, login string) (*User, error)
	Update(ctx context.Context, username, email, name, surname, birth string) (*User, error)
	Delete(ctx context.Context, username string) error
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, u *User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *Repository) GetByUsername(ctx context.Context, username string) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&u).Error
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &u, nil
}

func (r *Repository) GetByLogin(ctx context.Context, login string) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).Where("username = ? OR email = ?", login, login).First(&u).Error
	if err != nil {
		return nil, fmt.Errorf("User not found: %w", err)
	}
	return &u, nil
}

func (r *Repository) Update(ctx context.Context, username, email, name, surname, birth string) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).
		Where("username = ?", username).
		First(&u).Error
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	u.Email = email
	u.Name = name
	u.Surname = surname
	u.Birth = birth

	return &u, r.db.WithContext(ctx).Save(&u).Error
}

func (r *Repository) Delete(ctx context.Context, username string) error {
	if _, err := r.GetByUsername(ctx, username); err != nil {
		return err
	}
	return r.db.WithContext(ctx).
		Where("username = ?", username).
		Delete(&User{}).Error
}
