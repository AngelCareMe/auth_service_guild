package jwt

import "github.com/golang-jwt/jwt/v5"

type JWTRepository interface {
	GenerateAccessJWT(blizzardID int) (string, error)
	GenerateRefreshJWT(userID string) (string, error)
	ValidateJWT(tokenStr string) (*jwt.Token, error)
	GetUserIDByToken(token string) (any, error)
}
