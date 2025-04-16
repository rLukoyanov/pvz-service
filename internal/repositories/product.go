package repositories

import (
	"context"
	"errors"
	"pvz-service/internal/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	db   *pgxpool.Pool
	psql sq.StatementBuilderType
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db, psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar)}
}

func (r *ProductRepository) AddProduct(ctx context.Context, product models.Product, pvzID string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query, args, err := r.psql.
		Select("id").
		From("reception").
		Where(sq.Eq{"pvz_id": pvzID, "status": "in_progress"}).
		ToSql()
	if err != nil {
		return err
	}

	row := r.db.QueryRow(ctx, query, args...)

	var receptionId string
	err = row.Scan(&receptionId)

	if err == pgx.ErrNoRows {
		return errors.New("cant find opened receprions")
	}

	product.ReceptionId = receptionId

	query, args, err = r.psql.Insert("products").
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
