package entity

type BlizzardUser struct {
	ID        string `json:"id" db:"id"`
	BattleTag string `json:"battletag" db:"battletag"`
}
