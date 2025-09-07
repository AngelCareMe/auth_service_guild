package entity

type User struct {
	ID        string `json:"id" db:"id"`
	BattleTag string `json:"battletag" db:"battletag"`
}
