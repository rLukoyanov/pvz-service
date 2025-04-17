package repositories

import (
	"context"
	"pvz-service/internal/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type ProductRepository struct {
	db   *pgxpool.Pool
	psql sq.StatementBuilderType
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db, psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar)}
}

func (r *ProductRepository) AddProduct(ctx context.Context, product models.Product) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query, args, err := r.psql.Insert("products").
		Columns("date_time", "type", "reception_id").
		Values(product.DateTime, product.Type, product.ReceptionId).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil

}

func (r *ProductRepository) DeleteLastProduct(ctx context.Context, receptionId string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	subQuery := r.psql.
		Select("id").
		From("products").
		Where(sq.Eq{"reception_id": receptionId}).
		OrderBy("date_time DESC").
		Limit(1)

	query, args, err := r.psql.Delete("products").
		Where(sq.Expr("id = (?)", subQuery)).
		ToSql()

	logrus.Info(query)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
