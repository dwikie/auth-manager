package domain

import "time"

type User struct {
	ID        int
	Username  string
	Password  string
	Status    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserProfile struct {
	ID          int
	UserId      int
	FirstName   string
	LastName    string
	Email       string
	EmailStatus int
	Phone       string
	PhoneStatus int
}
