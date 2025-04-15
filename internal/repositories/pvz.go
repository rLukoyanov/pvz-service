package repositories

import (
	"context"
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

func (r *PVZRepository) CreatePVZ(ctx context.Context, pvz models.PVZ) (models.PVZ, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return models.PVZ{}, err
	}
	defer tx.Rollback(ctx)
	if pvz.RegistrationDate == "" {
		pvz.RegistrationDate = time.Now().Format(time.UnixDate)
	}

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
		Select("id", "city", "registratio_date").
		From("pvz").
		ToSql()
	if err != nil {
		return models.PVZ{}, err
	}

	logrus.Debug(query)

	row := r.db.QueryRow(ctx, query, args...)
	if err != nil {
		return models.PVZ{}, err
	}

	var pvz models.PVZ
	if err := row.Scan(&pvz.ID, &pvz.City, &pvz.RegistrationDate); err != nil {
		return models.PVZ{}, err
	}

	return pvz, nil
}

func (r *PVZRepository) UpdatePVZ(ctx context.Context, id string, pvz models.PVZ) (models.PVZ, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return models.PVZ{}, err
	}
	defer tx.Rollback(ctx)

	query, args, err := r.psql.
		Update("pvz").
		Set("city", pvz.City).
		Set("registration_date", pvz.RegistrationDate).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id, city, registration_date").
		ToSql()
	if err != nil {
		return models.PVZ{}, err
	}

	var updated models.PVZ
	err = tx.QueryRow(ctx, query, args...).Scan(&updated.ID, &updated.City, &updated.RegistrationDate)
	if err != nil {
		return models.PVZ{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return models.PVZ{}, err
	}

	return updated, nil
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
