package entity

type BlizzardUser struct {
	ID        int    `json:"id" db:"id"`
	BattleTag string `json:"battletag" db:"battletag"`
}
