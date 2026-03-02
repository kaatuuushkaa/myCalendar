package user

type User struct {
	ID       string `gorm:"type:uuid;primaryKey"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Email    string
	Name     string
	Surname  string
	Birth    string
}
