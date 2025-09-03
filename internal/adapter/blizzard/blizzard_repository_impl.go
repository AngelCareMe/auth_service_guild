package blizzard

import (
	"auth-service/internal/entity"
	"auth-service/pkg/config"
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type blizzardRepository struct {
	client *http.Client
	oauth  *oauth2.Config
	cfg    *config.Config
	log    *logrus.Logger
}

func NewBlizzardRepository(cfg *config.Config, log *logrus.Logger) *blizzardRepository {
	oauthConf := oauth2.Config{
		ClientID:     cfg.Blizzard.ClientID,
		ClientSecret: cfg.Blizzard.ClientSecret,
		RedirectURL:  cfg.Blizzard.RedirectURL,
		Scopes:       []string{"wow.profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://eu.battle.net/oauth/authorize",
			TokenURL: "https://eu.battle.net/oauth/token",
		},
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: false,
		},
	}

	return &blizzardRepository{
		client: client,
		oauth:  &oauthConf,
		cfg:    cfg,
		log:    log,
	}
}

func (br *blizzardRepository) GetAuthURL(state string) string {
	return br.oauth.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (br *blizzardRepository) HandleCallback(ctx context.Context, code string, userID int) (*entity.BlizzardToken, error) {
	token, err := br.oauth.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	bt := &entity.BlizzardToken{
		UserID:       userID,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
		TokenType:    token.TokenType,
	}

	return bt, nil
}
