package dto

type CreateEventInput struct {
	UserID      string
	Title       string
	Description string
	StartAt     string
	EndAt       string
	EventDate   string
}

type UpdateEventInput struct {
	Title       string
	Description string
	StartAt     string
	EndAt       string
	EventDate   string
}

type EventOutput struct {
	ID          string
	UserID      string
	Title       string
	Description string
	StartAt     string
	EndAt       string
	EventDate   string
}
