package model

import "time"

type User struct {
	ID       int
	Username string
	Password string
	Profile  UserProfile
}

type UserProfile struct {
	FirstName string
	LastName  string
}

type Session struct {
	UserID    int
	CSRFToken string
	ExpiresAt time.Time
}
