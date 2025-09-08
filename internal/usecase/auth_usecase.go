package usecase

import (
	"auth-service/internal/entity"
	"context"
)

type AuthUsecase interface {
	ValidateAccess(ctx context.Context, accessToken string) (int, error)
	RefreshTokens(ctx context.Context, refreshToken string) (string, string, error)
	HandleCallback(ctx context.Context, code string) (string, string, error)
	GetBlizzardUser(ctx context.Context, jwtAccess string) (*entity.BlizzardUser, error)
	GetBlizzardToken(ctx context.Context, jwtAccess string) (string, error)
}
