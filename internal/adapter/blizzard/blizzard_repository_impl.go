package blizzard

import (
	"auth-service/internal/entity"
	"auth-service/pkg/config"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type dtoBlizzUser struct {
	ID        int    `json:"id" db:"id"`
	BattleTag string `json:"battletag" db:"battletag"`
}

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
		Scopes:       []string{"openid", "wow.profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://oauth.battle.net/authorize",
			TokenURL: "https://oauth.battle.net/token",
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

func (br *blizzardRepository) HandleCallback(ctx context.Context, code string) (*entity.BlizzardToken, error) {
	if code == "" {
		return nil, fmt.Errorf("empty code err")
	}
	token, err := br.oauth.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	user, err := br.GetUser(ctx, token.AccessToken)
	if err != nil {
		return nil, err
	}

	bt := &entity.BlizzardToken{
		BlizzardID:   user.ID,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		Expiry:       token.Expiry,
	}

	return bt, nil
}

func (br *blizzardRepository) GetUser(ctx context.Context, token string) (*entity.BlizzardUser, error) {
	if token == "" {
		return nil, fmt.Errorf("empty token err")
	}

	reqUser, err := http.NewRequestWithContext(ctx, "GET", "https://oauth.battle.net/userinfo", nil)
	if err != nil {
		return nil, err
	}

	reqUser.Header.Set("Authorization", "Bearer "+token)

	respUser, err := br.client.Do(reqUser)
	if err != nil {
		return nil, err
	}
	defer respUser.Body.Close()

	if respUser.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed get user with status: %d", respUser.StatusCode)
	}

	body, err := io.ReadAll(respUser.Body)
	if err != nil {
		return nil, err
	}

	dto := &dtoBlizzUser{}

	if err := json.Unmarshal(body, dto); err != nil {
		return nil, err
	}

	bu := &entity.BlizzardUser{
		ID:        strconv.Itoa(dto.ID),
		BattleTag: dto.BattleTag,
	}

	return bu, nil
}

func (br *blizzardRepository) RefreshToken(ctx context.Context, refreshToken string) (*entity.BlizzardToken, error) {
	if refreshToken == "" {
		return nil, fmt.Errorf("empty token")
	}

	ts := br.oauth.TokenSource(ctx, &oauth2.Token{RefreshToken: refreshToken})
	newToken, err := ts.Token()
	if err != nil {
		return nil, err
	}

	return &entity.BlizzardToken{
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
		TokenType:    newToken.TokenType,
		Expiry:       newToken.Expiry,
	}, nil
}
