package dto

type CreteUserInput struct {
	Username string
	Password string
	Email    string
	Name     string
	Surname  string
	Birth    string
}

type UpdateUserInput struct {
	Email   string
	Name    string
	Surname string
	Birth   string
}

type UserOutput struct {
	ID       string
	Username string
	Email    string
	Name     string
	Surname  string
	Birth    string
}

type Auth struct {
	Login    string
	Password string
}

type AuthOutput struct {
	AccessToken  string
	RefreshToken string
}
