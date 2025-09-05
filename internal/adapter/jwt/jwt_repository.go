package jwt

import "github.com/golang-jwt/jwt/v5"

type JWTRepository interface {
	GenerateAccessJWT(blizzardID string) (string, error)
	GenerateRefreshJWT(userID string) (string, error)
	ValidateJWT(tokenStr string) (*jwt.Token, error)
	GetUserIDByToken(token string) (string, error)
}
