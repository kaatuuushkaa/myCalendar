//передаем между слоями системы

package domain

import "time"

type Event struct {
	ID          string `gorm:"type:uuid;primaryKey"`
	UserID      string `gorm:"type:uuid;not null;index"`
	Title       string `gorm:"not null"`
	Description string
	StartAt     time.Time `gorm:"not null"`
	EndAt       time.Time `gorm:"not null"`
	EventDate   string
}
