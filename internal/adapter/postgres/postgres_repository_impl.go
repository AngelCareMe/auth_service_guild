package postgres

import (
	"auth-service/internal/entity"
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type postgresRepository struct {
	pool *pgxpool.Pool
	log  *logrus.Logger
}

func NewPostgresRepository(pool *pgxpool.Pool, log *logrus.Logger) *postgresRepository {
	return &postgresRepository{
		pool: pool,
		log:  log,
	}
}

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func (pg *postgresRepository) SaveUser(ctx context.Context, userID, battletag string) error {
	query, args, err := psql.
		Insert("users").
		Columns("id", "battletag").
		Values(userID, battletag).
		ToSql()

	if err != nil {
		return err
	}

	_, err = pg.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (pg *postgresRepository) SaveBlizzardUser(ctx context.Context, userID string, battletag string) error {
	query, args, err := psql.
		Insert("blizzard_users").
		Columns("id", "battletag").
		Values(userID, battletag).
		Suffix("ON CONFLICT (id) DO UPDATE SET battletag = EXCLUDED.battletag").
		ToSql()

	if err != nil {
		return err
	}

	_, err = pg.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (pg *postgresRepository) SaveBlizzardToken(ctx context.Context, userID string, blizzardID string, token *entity.BlizzardToken) error {
	query, args, err := psql.
		Insert("blizzard_tokens").
		Columns(
			"user_id",
			"blizzard_id",
			"access_token",
			"refresh_token",
			"expiry",
			"token_type",
		).
		Values(
			userID,
			blizzardID,
			token.AccessToken,
			token.RefreshToken,
			token.Expiry,
			token.TokenType,
		).
		Suffix("ON CONFLICT (user_id) DO UPDATE SET access_token = EXCLUDED.access_token, refresh_token = EXCLUDED.refresh_token, expiry = EXCLUDED.expiry, token_type = EXCLUDED.token_type").
		ToSql()

	if err != nil {
		return err
	}

	_, err = pg.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (pg *postgresRepository) SaveJWTToken(ctx context.Context, userID, battletag string, token *entity.JWTToken) error {
	query, args, err := psql.
		Insert("jwt_tokens").
		Columns("user_id", "battletag", "refresh_token", "expiry").
		Values(userID, battletag, token.RefreshToken, token.Expiry).
		Suffix("ON CONFLICT (user_id) DO UPDATE SET refresh_token = EXCLUDED.refresh_token, expiry = EXCLUDED.expiry").
		ToSql()

	if err != nil {
		return err
	}

	_, err = pg.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (pg *postgresRepository) GetUser(ctx context.Context, battletag string) (*entity.User, error) {
	query, args, err := psql.
		Select("id", "battletag").
		From("users").
		Where(sq.Eq{"battletag": battletag}).
		ToSql()

	if err != nil {
		return nil, err
	}

	var user entity.User

	if err := pg.pool.QueryRow(ctx, query, args...).Scan(&user.ID, &user.BattleTag); err != nil {
		return nil, err
	}

	return &user, nil
}

func (pg *postgresRepository) GetBlizzardUser(ctx context.Context, battletag string) (*entity.BlizzardUser, error) {
	query, args, err := psql.
		Select("id", "battletag").
		From("blizzard_users").
		Where(sq.Eq{"battletag": battletag}).
		ToSql()

	if err != nil {
		return nil, err
	}

	var buser entity.BlizzardUser

	if err := pg.pool.QueryRow(ctx, query, args...).Scan(&buser.ID, &buser.BattleTag); err != nil {
		return nil, err
	}

	return &buser, nil
}

func (pg *postgresRepository) GetBlizzardUserByID(ctx context.Context, blizzardID string) (*entity.BlizzardUser, error) {
	query, args, err := psql.
		Select("id", "battletag").
		From("blizzard_users").
		Where(sq.Eq{"id": blizzardID}).
		ToSql()

	if err != nil {
		return nil, err
	}

	var buser entity.BlizzardUser

	if err := pg.pool.QueryRow(ctx, query, args...).Scan(&buser.ID, &buser.BattleTag); err != nil {
		return nil, err
	}

	return &buser, nil
}

func (pg *postgresRepository) GetBlizzardTokenByUserID(ctx context.Context, userID string) (*entity.BlizzardToken, error) {
	query, args, err := psql.
		Select(
			"user_id",
			"blizzard_id",
			"access_token",
			"refresh_token",
			"expiry",
			"token_type",
		).
		From("blizzard_tokens").
		Where(sq.Eq{"user_id": userID}).
		ToSql()

	if err != nil {
		return nil, err
	}

	var bt entity.BlizzardToken
	if err = pg.pool.QueryRow(ctx, query, args...).
		Scan(
			&bt.UserID,
			&bt.BlizzardID,
			&bt.AccessToken,
			&bt.RefreshToken,
			&bt.Expiry,
			&bt.TokenType,
		); err != nil {
		return nil, err
	}

	return &bt, nil
}

func (pg *postgresRepository) GetBlizzardTokenByBlizzID(ctx context.Context, blizzID string) (*entity.BlizzardToken, error) {
	query, args, err := psql.
		Select(
			"user_id",
			"blizzard_id",
			"access_token",
			"refresh_token",
			"expiry",
			"token_type",
		).
		From("blizzard_tokens").
		Where(sq.Eq{"blizzard_id": blizzID}).
		ToSql()

	if err != nil {
		return nil, err
	}

	var bt entity.BlizzardToken
	if err = pg.pool.QueryRow(ctx, query, args...).
		Scan(
			&bt.UserID,
			&bt.BlizzardID,
			&bt.AccessToken,
			&bt.RefreshToken,
			&bt.Expiry,
			&bt.TokenType,
		); err != nil {
		return nil, err
	}

	return &bt, nil
}

func (pg *postgresRepository) GetJWTToken(ctx context.Context, userID string) (*entity.JWTToken, error) {
	query, args, err := psql.
		Select(
			"user_id",
			"battletag",
			"refresh_token",
			"expiry",
		).
		From("jwt_tokens").
		Where(sq.Eq{"user_id": userID}).
		ToSql()

	if err != nil {
		return nil, err
	}

	var jwt entity.JWTToken
	if err = pg.pool.QueryRow(ctx, query, args...).
		Scan(
			&jwt.UserID,
			&jwt.BattleTag,
			&jwt.RefreshToken,
			&jwt.Expiry,
		); err != nil {
		return nil, err
	}

	return &jwt, nil
}
