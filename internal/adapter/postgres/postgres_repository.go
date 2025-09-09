package postgres

import (
	"auth-service/internal/entity"
	"context"
)

type PostgresRepository interface {
	SaveUser(ctx context.Context, userID, battletag string) error
	SaveBlizzardUser(ctx context.Context, userID string, battletag string) error
	SaveBlizzardToken(ctx context.Context, userID string, blizzardID string, token *entity.BlizzardToken) error
	SaveJWTToken(ctx context.Context, userID, battletag string, token *entity.JWTToken) error
	GetUser(ctx context.Context, battletag string) (*entity.User, error)
	GetBlizzardUser(ctx context.Context, battletag string) (*entity.BlizzardUser, error)
	GetBlizzardUserByID(ctx context.Context, blizzardID string) (*entity.BlizzardUser, error)
	GetBlizzardTokenByUserID(ctx context.Context, userID string) (*entity.BlizzardToken, error)
	GetJWTToken(ctx context.Context, userID string) (*entity.JWTToken, error)
	GetBlizzardTokenByBlizzID(ctx context.Context, blizzID string) (*entity.BlizzardToken, error)
}
