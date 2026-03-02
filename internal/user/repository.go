package user

import (
	"context"
	"fmt"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
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
	return r.db.WithContext(ctx).
		Where("username = ?", username).
		Delete(&User{}).Error
}
