package repositories

import (
	"context"
	"pvz-service/internal/models"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReceptionRepository struct {
	db   *pgxpool.Pool
	psql sq.StatementBuilderType
}

func NewReceptionRepository(db *pgxpool.Pool) *ReceptionRepository {
	return &ReceptionRepository{db: db, psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar)}
}

func (r *ReceptionRepository) CreateReception(ctx context.Context, Reception models.Reception) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query, args, err := r.psql.
		Insert("Reception").
		Columns("pvz_id", "status", "date_time").
		Values(Reception.PvzId, "in_progress", time.Now().Format(time.UnixDate)).
		Suffix("RETURNING id, date_time, pvz_id, status").
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

func (r *ReceptionRepository) GetActiveReceptionByPVZID(ctx context.Context, pvzID string) (*models.Reception, error) {
	query, args, err := r.psql.
		Select("id", "date_time", "pvz_id", "status").
		From("Reception").
		Where(sq.Eq{"pvz_id": pvzID, "status": "in_progress"}).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRow(ctx, query, args...)
	if row == nil {
		return nil, nil
	}

	var Reception models.Reception
	err = row.Scan(&Reception.ID, &Reception.DateTime, &Reception.PvzId, &Reception.Status)
	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &Reception, nil
}
