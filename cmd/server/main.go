package main

import (
	"context"
	"pvz-service/config"
	"pvz-service/internal/database"
	"pvz-service/internal/handlers"
	"pvz-service/internal/logger"
	"pvz-service/internal/middlewares"
	"pvz-service/internal/repositories"

	"github.com/labstack/echo/v4"
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

	// Register Swagger
	handlers.RegisterSwagger(e)

	authMiddleware := middlewares.NewAuthMiddleware(cfg)

	dlHandler := handlers.NewDummyLoginHandler(cfg)
	e.POST("/dummyLogin", dlHandler.DummyLogin)

	authRepo := repositories.NewUserRepository(pool)
	authHandler := handlers.NewAuthHandler(authRepo, cfg)
	e.POST("/register", authHandler.Register)
	e.POST("/login", authHandler.Login)

	g := e.Group("/pvz")
	g.Use(authMiddleware.JWTMiddleware())
	pvzRepo := repositories.NewPVZRepository(pool)
	pvzHandler := handlers.NewPVZHandler(pvzRepo)

	g.POST("/", pvzHandler.Create, authMiddleware.RequireRole("moderator"))
	g.GET("/:id", pvzHandler.GetByID)

	// Reception endpoints
	ReceptionRepo := repositories.NewReceptionRepository(pool)
	ReceptionHandler := handlers.NewReceptionHandler(ReceptionRepo)
	e.POST("/receptions", ReceptionHandler.Create, authMiddleware.JWTMiddleware(), authMiddleware.RequireRole("client"))

	ProductRepo := repositories.NewProductRepository(pool)
	productHandler := handlers.NewProductHandler(ProductRepo)
	e.POST("/product", productHandler.AddProduct)

	e.Logger.Fatal(e.Start(":8080"))

	//graseful shutdown
}
