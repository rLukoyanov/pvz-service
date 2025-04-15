package repositories

import (
	"context"
	"fmt"
	"pvz-service/internal/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type UserRepository struct {
	db   *pgxpool.Pool
	psql sq.StatementBuilderType
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db, psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar)}
}

func (r *UserRepository) CreateUser(ctx context.Context, user models.User) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query, args, err := r.psql.
		Insert("users").
		Columns("email", "password", "role").
		Values(user.Email, user.Password, user.Role).
		ToSql()
	if err != nil {
		return err
	}

	logrus.Debug(query)

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	return nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return models.User{}, err
	}

	defer conn.Release()

	query, args, err := r.psql.
		Select("id", "email", "password", "role").
		From("users").
		Where(sq.Eq{"email": email}).
		ToSql()
	if err != nil {
		return models.User{}, err
	}

	row := conn.QueryRow(ctx, query, args...)

	var user models.User
	err = row.Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
