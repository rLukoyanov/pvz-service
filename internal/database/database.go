package database

import (
	"context"
	"pvz-service/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

func ConnectDB(cfg *config.Config, ctx context.Context) *pgxpool.Pool {
	var pool *pgxpool.Pool
	var err error

	maxRetries := 5
	retryInterval := 5 * time.Second

	for i := 1; i <= maxRetries; i++ {
		pool, err = pgxpool.New(ctx, cfg.DATABASE_URL)

		if err == nil {
			err = pool.Ping(ctx)
			if err == nil {
				logrus.Infof("Успешно подключено к базе данных на попытке %d", i)
				return pool
			}
		}

		logrus.Errorf("ошибка подключения к базе данных (попытка %d): %s", i, err.Error())

		if i < maxRetries {
			time.Sleep(retryInterval)
		}
	}

	logrus.Fatal("Не удалось подключиться к базе данных после нескольких попыток.")
	return nil
}
