package usecase

import (
	"auth-service/internal/adapter/blizzard"
	"auth-service/internal/adapter/jwt"
	"auth-service/internal/adapter/postgres"
	"auth-service/internal/entity"
	"context"
	"fmt"
	"strconv"
	"time"

	jwtPkg "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type authUsecase struct {
	jwtAd      jwt.JWTRepository
	blizzardAd blizzard.BlizzardRepository
	dbAd       postgres.PostgresRepository
	log        *logrus.Logger
}

func NewAuthUsecase(
	jwtAd jwt.JWTRepository,
	blizzardAd blizzard.BlizzardRepository,
	dbAd postgres.PostgresRepository,
	log *logrus.Logger,
) *authUsecase {
	return &authUsecase{
		jwtAd:      jwtAd,
		blizzardAd: blizzardAd,
		dbAd:       dbAd,
		log:        log,
	}
}

func (uc *authUsecase) HandleCallback(ctx context.Context, code string) (string, string, error) {
	if code == "" {
		return "", "", fmt.Errorf("empty code")
	}

	bt, err := uc.blizzardAd.HandleCallback(ctx, code)
	if err != nil {
		return "", "", err
	}

	bUser, err := uc.blizzardAd.GetUser(ctx, bt.AccessToken)
	if err != nil {
		return "", "", err
	}

	existingUser, err := uc.dbAd.GetUser(ctx, bUser.BattleTag)
	var userID string
	if err != nil {
		userID = uuid.NewString()
		if err := uc.dbAd.SaveUser(ctx, userID, bUser.BattleTag); err != nil {
			return "", "", err
		}
	} else {
		userID = existingUser.ID
	}

	if err := uc.dbAd.SaveBlizzardUser(ctx, bUser.ID, bUser.BattleTag); err != nil {
		return "", "", err
	}

	if err := uc.dbAd.SaveBlizzardToken(ctx, userID, bt.BlizzardID, bt); err != nil {
		return "", "", err
	}

	jwtAccess, err := uc.jwtAd.GenerateAccessJWT(bUser.ID)
	if err != nil {
		return "", "", err
	}

	jwtRefresh, err := uc.jwtAd.GenerateRefreshJWT(userID)
	if err != nil {
		return "", "", err
	}

	jwt := &entity.JWTToken{
		UserID:       userID,
		BattleTag:    bUser.BattleTag,
		AccessToken:  jwtAccess,
		RefreshToken: jwtRefresh,
		Expiry:       time.Now().Add(30 * 24 * time.Hour),
	}

	if err := uc.dbAd.SaveJWTToken(ctx, userID, bUser.BattleTag, jwt); err != nil {
		return "", "", err
	}

	return jwtAccess, jwtRefresh, nil
}

func (uc *authUsecase) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	if refreshToken == "" {
		return "", "", fmt.Errorf("empty token")
	}

	token, err := uc.jwtAd.ValidateJWT(refreshToken)
	if err != nil || token == nil || !token.Valid {
		return "", "", err
	}

	claims := token.Claims.(jwtPkg.MapClaims)
	userID, err := ExtractSub(claims, "refresh")
	if err != nil {
		return "", "", err
	}

	blizzardToken, err := uc.dbAd.GetBlizzardToken(ctx, userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to get blizzard token: %w", err)
	}

	if time.Now().After(blizzardToken.Expiry) {
		newBt, err := uc.blizzardAd.RefreshToken(ctx, blizzardToken.RefreshToken)
		if err != nil {
			return "", "", err
		}

		if err := uc.dbAd.SaveBlizzardToken(ctx, userID, newBt.BlizzardID, newBt); err != nil {
			return "", "", err
		}
		blizzardToken = newBt
	}

	bUser, err := uc.blizzardAd.GetUser(ctx, blizzardToken.AccessToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to get blizzard user: %w", err)
	}

	newAccess, err := uc.jwtAd.GenerateAccessJWT(bUser.ID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefresh, err := uc.jwtAd.GenerateRefreshJWT(userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	jwtTokenEntity := &entity.JWTToken{
		UserID:       userID,
		BattleTag:    bUser.BattleTag,
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
		Expiry:       time.Now().Add(30 * 24 * time.Hour),
	}

	if err := uc.dbAd.SaveJWTToken(ctx, userID, bUser.BattleTag, jwtTokenEntity); err != nil {
		return "", "", fmt.Errorf("failed to save jwt tokens: %w", err)
	}

	return newAccess, newRefresh, nil
}

func (uc *authUsecase) ValidateAccess(ctx context.Context, accessToken string) (int, error) {
	token, err := uc.jwtAd.ValidateJWT(accessToken)
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwtPkg.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid claims")
	}

	typ, ok := claims["typ"].(string)
	if !ok || typ != "access" {
		return 0, fmt.Errorf("invalid type")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return 0, fmt.Errorf("no sub in token")
	}

	blizzardID, err := strconv.Atoi(sub)
	if err != nil {
		return 0, err
	}

	return blizzardID, nil
}

func (uc *authUsecase) GetBlizzardToken(ctx context.Context, jwtAccess string) (string, error) {
	token, err := uc.jwtAd.ValidateJWT(jwtAccess)
	if err != nil {
		return "", err
	}

	claim := token.Claims.(jwtPkg.MapClaims)
	blizzardID, err := ExtractSub(claim, "access")
	if err != nil {
		return "", err
	}

	blizzToken, err := uc.dbAd.GetBlizzardToken(ctx, blizzardID)
	if err != nil {
		return "", err
	}

	if time.Now().After(blizzToken.Expiry) {
		newBt, err := uc.blizzardAd.RefreshToken(ctx, blizzToken.RefreshToken)
		if err != nil {
			return "", err
		}
		if err := uc.dbAd.SaveBlizzardToken(ctx, blizzToken.UserID, blizzardID, newBt); err != nil {
			return "", err
		}
		blizzToken = newBt
	}

	return blizzToken.AccessToken, nil
}

func (uc *authUsecase) GetBlizzardUser(ctx context.Context, jwtAccess string) (*entity.BlizzardUser, error) {
	if jwtAccess == "" {
		return nil, fmt.Errorf("empty token")
	}

	token, err := uc.jwtAd.ValidateJWT(jwtAccess)
	if err != nil || token == nil || !token.Valid {
		return nil, err
	}

	claims := token.Claims.(jwtPkg.MapClaims)
	blizzardIDStr, err := ExtractSub(claims, "access")
	if err != nil {
		return nil, err
	}
	return uc.dbAd.GetBlizzardUserByID(ctx, blizzardIDStr)
}

func ExtractSub(claims jwtPkg.MapClaims, expectedType string) (string, error) {
	typ, ok := claims["typ"].(string)
	if !ok || typ != expectedType {
		return "", fmt.Errorf("unexpected token type")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("invalid sub")
	}

	return sub, nil
}
