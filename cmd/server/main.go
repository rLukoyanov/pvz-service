package main

import (
	"context"
	"net/http"
	"pvz-service/config"
	"pvz-service/internal/database"
	"pvz-service/internal/logger"
	"pvz-service/internal/pkg/jwt"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
	// init config
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	// init logger
	logger.InitLogger(cfg.LOG_LEVEL, cfg.MODE)

	// init db
	pool := database.ConnectDB(cfg, ctx)
	_ = pool
	// init echo
	e := echo.New()

	e.POST("/dummyLogin", func(c echo.Context) error {
		type req struct {
			Role string `json:"role"`
		}

		var r req
		if err := c.Bind(&r); err != nil {
			logrus.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
		}

		if r.Role != "client" && r.Role != "moderator" {
			logrus.Error("invalid role")
			return echo.NewHTTPError(http.StatusBadRequest, "invalid role")
		}

		token, err := jwt.GenerateToken(r.Role, cfg)
		if err != nil {
			logrus.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "could not generate token")
		}

		return c.JSON(http.StatusOK, echo.Map{
			"token": token,
		})
	})

	e.Logger.Fatal(e.Start(":8080"))

	//graseful shutdown
}
