package entity

type User struct {
	ID        int    `json:"user_id" db:"user_id"`
	BattleTag string `json:"battle_tag" db:"battle_tag"`
}
