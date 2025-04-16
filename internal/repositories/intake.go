package repositories

import (
	"context"
	"pvz-service/internal/models"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IntakeRepository struct {
	db   *pgxpool.Pool
	psql sq.StatementBuilderType
}

func NewIntakeRepository(db *pgxpool.Pool) *IntakeRepository {
	return &IntakeRepository{db: db, psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar)}
}

func (r *IntakeRepository) CreateIntake(ctx context.Context, intake models.Intake) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query, args, err := r.psql.
		Insert("intake").
		Columns("pvz_id", "status", "date_time").
		Values(intake.PvzId, "in_progress", time.Now().Format(time.UnixDate)).
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

func (r *IntakeRepository) GetActiveIntakeByPVZID(ctx context.Context, pvzID string) (*models.Intake, error) {
	query, args, err := r.psql.
		Select("id", "date_time", "pvz_id", "status").
		From("intake").
		Where(sq.Eq{"pvz_id": pvzID, "status": "in_progress"}).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRow(ctx, query, args...)
	if row == nil {
		return nil, nil
	}

	var intake models.Intake
	err = row.Scan(&intake.ID, &intake.DateTime, &intake.PvzId, &intake.Status)
	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &intake, nil
}
