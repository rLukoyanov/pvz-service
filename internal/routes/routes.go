package routes

import (
	"pvz-service/config"
	"pvz-service/internal/handlers"
	"pvz-service/internal/middlewares"
	"pvz-service/internal/repositories"
	"pvz-service/internal/services"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo, cfg *config.Config, db *pgxpool.Pool) {
	authMiddleware := middlewares.NewAuthMiddleware(cfg)

	authRepo := repositories.NewUserRepository(db)
	pvzRepo := repositories.NewPVZRepository(db)
	productRepo := repositories.NewProductRepository(db)
	receptionRepo := repositories.NewReceptionRepository(db)

	repos := repositories.Repos{
		Cfg:           cfg,
		AuthRepo:      authRepo,
		PvzRepo:       pvzRepo,
		ProductRepo:   productRepo,
		ReceptionRepo: receptionRepo,
	}

	userService := services.NewUserService(&repos)
	productService := services.NewProductService(&repos)
	pvzService := services.NewPVZService(&repos)
	ReceptionService := services.NewReceptionService(&repos)

	services := services.Services{
		UserService:      userService,
		ProductService:   productService,
		PvzService:       pvzService,
		ReceptionService: ReceptionService,
		Cfg:              cfg,
	}

	dlHandler := handlers.NewDummyLoginHandler(&services)
	authHandler := handlers.NewAuthHandler(&services)
	pvzHandler := handlers.NewPVZHandler(&services)
	receptionHandler := handlers.NewReceptionHandler(&services)
	productHandler := handlers.NewProductHandler(&services)

	e.POST("/dummyLogin", dlHandler.DummyLogin)

	e.POST("/register", authHandler.Register)
	e.POST("/login", authHandler.Login)

	g := e.Group("/pvz")
	g.Use(authMiddleware.JWTMiddleware())

	g.POST("/", pvzHandler.Create, authMiddleware.RequireRole("moderator"))
	g.GET("/:id", pvzHandler.GetByID)
	g.POST("/:id/delete_last_product", pvzHandler.DeleteLastProduct, authMiddleware.RequireRole("client"))

	e.POST("/receptions", receptionHandler.Create, authMiddleware.JWTMiddleware(), authMiddleware.RequireRole("client"))

	e.POST("/product", productHandler.AddProduct)
}
