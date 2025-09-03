package blizzard

import (
	"auth-service/internal/entity"
	"context"
)

type BlizzardRepository interface {
	GetAuthURL(state string) string
	HandleCallback(ctx context.Context, code string) (*entity.BlizzardToken, error)
	GetUser(ctx context.Context, token string) (*entity.BlizzardUser, error)
	RefreshToken(ctx context.Context, refreshToken string) (*entity.BlizzardToken, error)
}
