package usecase

import "context"

type AuthUsecase interface {
	ValidateAccess(ctx context.Context, accessToken string) (int, error)
	RefreshTokens(ctx context.Context, refreshToken string) (string, string, error)
	HandleCallback(ctx context.Context, code string) (string, string, error)
}
