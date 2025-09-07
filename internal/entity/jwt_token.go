package entity

import (
	"time"
)

type JWTToken struct {
	UserID       string    `json:"user_id" db:"user_id"`
	BattleTag    string    `json:"battletag" db:"battletag"`
	AccessToken  string    `json:"access_token" db:"access_token"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	Expiry       time.Time `json:"expiry" db:"expiry"`
}
