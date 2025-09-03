package entity

import "time"

type Token struct {
	UserID       int
	RefreshToken string
	ExpireAt     time.Time
}
