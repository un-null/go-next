package main

import (
	"backend/internal/delivery/http"
	"backend/internal/repository"
	"backend/internal/usecase"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	_ = godotenv.Load()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	e := echo.New()

	e.Validator = &CustomValidator{validator: validator.New()}

	// Middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.DELETE},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	// DI
	authMiddleware := http.NewAuthMiddleware(jwtSecret)

	userRepo := repository.NewUserRepository()
	userUC := usecase.NewUserUseCase(userRepo)

	productRepo := repository.NewProductRepository()
	productUC := usecase.NewProductUseCase(productRepo)

	cartRepo := repository.NewCartRepository()
	cartUC := usecase.NewCartUseCase(cartRepo)

	userHandler := http.NewUserHandler(userUC, jwtSecret)
	productHandler := http.NewProductHandler(productUC)
	cartHandler := http.NewCartHandler(cartUC)

	// Route groupin

	api := e.Group("/api")

	// Public endpoints
	public := api.Group("")
	public.POST("/signup", userHandler.SignUp)
	public.POST("/login", userHandler.Login)
	productHandler.RegisterRoutes(api)

	// Protected endpoints
	protected := api.Group("")
	protected.Use(authMiddleware.Middleware)
	protected.GET("/users", userHandler.GetUsers)
	cartHandler.RegisterRoutes(protected)

	log.Println("Server running on :8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
