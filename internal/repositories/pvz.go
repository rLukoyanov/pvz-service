package repositories

import (
	"context"
	"fmt"
	"net/http"
	"pvz-service/internal/models"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type PVZRepository struct {
	db   *pgxpool.Pool
	psql sq.StatementBuilderType
}

func NewPVZRepository(db *pgxpool.Pool) *PVZRepository {
	return &PVZRepository{db: db, psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar)}
}

func (r *PVZRepository) GetAll(ctx context.Context, limit, offset int, from, to time.Time) ([]models.PVZ, error) {
	query := r.psql.
		Select("DISTINCT pvz.id", "pvz.city", "pvz.registration_date").
		From("pvz").
		Join("reception ON pvz.id = reception.pvz_id")

	if !from.IsZero() {
		query = query.Where(sq.GtOrEq{"reception.date_time": from})
	}
	if !to.IsZero() {
		query = query.Where(sq.LtOrEq{"reception.date_time": to})
	}

	query = query.OrderBy("pvz.registration_date DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var result []models.PVZ
	for rows.Next() {
		var pvz models.PVZ
		if err := rows.Scan(&pvz.ID, &pvz.City, &pvz.RegistrationDate); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		result = append(result, pvz)
	}

	return result, nil
}

func (r *PVZRepository) CreatePVZ(ctx context.Context, pvz models.PVZ) (models.PVZ, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return models.PVZ{}, err
	}
	defer tx.Rollback(ctx)
	pvz.RegistrationDate = time.Now()

	query, args, err := r.psql.
		Insert("pvz").
		Columns("city", "registration_date").
		Values(pvz.City, pvz.RegistrationDate).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return models.PVZ{}, err
	}

	var id string
	err = tx.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return models.PVZ{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return models.PVZ{}, err
	}

	pvz.ID = id

	return pvz, nil
}

func (r *PVZRepository) GetPVZByID(ctx context.Context, id string) (models.PVZ, error) {
	query, args, err := r.psql.
		Select("id", "city", "registration_date").
		From("pvz").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return models.PVZ{}, err
	}

	logrus.Debug(query)

	row := r.db.QueryRow(ctx, query, args...)

	var pvz models.PVZ
	if err := row.Scan(&pvz.ID, &pvz.City, &pvz.RegistrationDate); err != nil {
		return models.PVZ{}, err
	}

	return pvz, nil
}

func (r *PVZRepository) DeletePVZ(ctx context.Context, id string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query, args, err := r.psql.
		Delete("pvz").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	result, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "PVZ not found")
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
