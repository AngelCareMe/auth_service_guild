package entity

import (
	"time"
)

type BlizzardToken struct {
	UserID       string    `json:"user_id" db:"user_id"`
	BlizzardID   string    `json:"blizzard_id" db:"blizzard_id"`
	AccessToken  string    `json:"access_token" db:"access_token"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	Expiry       time.Time `json:"expiry" db:"expiry"`
	TokenType    string    `json:"token_type" db:"token_type"`
}
