package jwt

import (
	"auth-service/pkg/config"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type jwtRepository struct {
	log *logrus.Logger
	cfg config.Config
}

func NewJWTRepository(log *logrus.Logger, cfg config.Config) *jwtRepository {
	return &jwtRepository{
		log: log,
		cfg: cfg,
	}
}

func (jr *jwtRepository) GenerateAccessJWT(blizzardID int) (string, error) {
	claims := jwt.MapClaims{
		"sub": blizzardID,
		"exp": jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
		"iat": jwt.NewNumericDate(time.Now()),
		"typ": "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(jr.cfg.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (jr *jwtRepository) GenerateRefreshJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
		"iat": jwt.NewNumericDate(time.Now()),
		"typ": "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(jr.cfg.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (jr *jwtRepository) ValidateJWT(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("bad method")
		}
		return []byte(jr.cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}
