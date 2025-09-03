package entity

import "time"

type BlizzardToken struct {
	UserID       int
	AccessToken  string
	RefreshToken string
	Expiry       time.Time
	TokenType    string
}
