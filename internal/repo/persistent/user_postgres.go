package persistent

import (
	"NEW_LOE_bot/internal/entity"
	"NEW_LOE_bot/pkg/postgres"
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) Store(ctx context.Context, user *entity.User) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("UserRepo - Store - r.Pool.Begin: %w", err)
	}
	defer tx.Rollback(ctx)

	err = r.storeUser(ctx, tx, user)
	if err != nil {
		return fmt.Errorf("UserRepo - Store: %w", err)
	}

	for _, group := range user.Groups {
		err = r.storeSubscription(ctx, tx, user, group)
		if err != nil {
			return fmt.Errorf("UserRepo - Store: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("UserRepo - Store - tx.Commit: %w", err)
	}

	return nil
}

func (r *UserRepo) storeUser(ctx context.Context, tx pgx.Tx, user *entity.User) error {
	userSql, args, err := r.Builder.
		Insert("users").
		Columns("id, telegramID, username, firstname, is_active, state, created_at").
		Values(user.ID, user.TelegramID, user.Username, user.FirstName, user.IsActive, user.State, user.CreatedAt).
		ToSql()

	if err != nil {
		return fmt.Errorf("storeUser - r.Builder: %w", err)
	}

	_, err = tx.Exec(ctx, userSql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return entity.ErrUserAlreadyExists
		}

		return fmt.Errorf("storeUser - tx.Exec: %w", err)
	}

	return nil
}

func (r *UserRepo) storeSubscription(ctx context.Context, tx pgx.Tx, user *entity.User, group entity.GroupID) error {
	subsSql, args, err := r.Builder.
		Insert("users_subscriptions").
		Columns("user_id, group_id").
		Values(user.ID, group).
		ToSql()

	if err != nil {
		return fmt.Errorf("storeSubscription - r.Builder: %w", err)
	}

	_, err = tx.Exec(ctx, subsSql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return entity.ErrSubscriptionAlreadyExists
		}

		return fmt.Errorf("storeSubscription - tx.Exec: %w", err)
	}
	return nil
}

func (r *UserRepo) GetByTelegramID(ctx context.Context, telegramID int64) (entity.User, error) {
	return r.getUser(ctx, "telegramID", telegramID)
}

func (r *UserRepo) getUser(ctx context.Context, column string, value any) (entity.User, error) {
	sql, args, err := r.Builder.
		Select("id, telegramID, username, firstname, is_active, state, created_at").
		From("users").
		Where(squirrel.Eq{column: value}).
		ToSql()

	if err != nil {
		return entity.User{}, fmt.Errorf("UserRepo - getUser - r.Builder: %w", err)
	}

	var user entity.User

	err = r.Pool.QueryRow(ctx, sql, args...).
		Scan(&user.ID, &user.TelegramID, &user.Username, &user.FirstName, &user.IsActive, &user.State, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, entity.ErrUserNotFound
		}

		return entity.User{}, fmt.Errorf("UserRepo - getUser - r.Pool.QueryRow: %w", err)
	}

	err = r.getUserGroups(ctx, &user)
	if err != nil {
		return entity.User{}, fmt.Errorf("UserRepo - getUser - getUserGroups: %w", err)
	}

	return user, nil
}

func (r *UserRepo) getUserGroups(ctx context.Context, user *entity.User) error {
	sql, args, err := r.Builder.
		Select("group_id").
		From("users_subscriptions").
		Where(squirrel.Eq{"user_id": user.ID}).
		OrderBy("group_id").
		ToSql()

	if err != nil {
		return fmt.Errorf("UserRepo - getUserGroups - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)

	if err != nil {
		return fmt.Errorf("UserRepo - getUserGroups - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var group entity.GroupID
		err = rows.Scan(&group)
		if err != nil {
			return fmt.Errorf("UserRepo - getUserGroups - rows.Scan: %w", err)
		}
		user.Groups = append(user.Groups, group)
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("UserRepo - getUserGroups - rows.Err: %w", err)
	}

	return nil
}

func (r *UserRepo) Update(ctx context.Context, user *entity.User) error {
	sql, args, err := r.Builder.
		Update("users").
		Set("username", user.Username).
		Set("firstname", user.FirstName).
		Set("is_active", user.IsActive).
		Set("state", user.State).
		Where(squirrel.Eq{"telegramID": user.TelegramID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("UserRepo - Update - r.Builder: %w", err)
	}
	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UserRepo - Update - r.Pool.Exec: %w", err)
	}

	return nil
}
