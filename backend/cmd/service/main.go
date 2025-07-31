package main

import (
	"backend/internal/delivery/http"
	"backend/internal/repository"
	"backend/internal/usecase"
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// auth middleware
func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// retrieve userID from query paramater for testing
		userIDStr := c.QueryParam("user_id")
		if userIDStr != "" {
			if userID, err := strconv.Atoi(userIDStr); err == nil {
				c.Set("userID", userID)
			}
		}
		return next(c)
	}
}

func main() {
	e := echo.New()

	// Allow requests from frontend dev server (5173)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{echo.GET, echo.POST, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Use(authMiddleware)

	// User setup
	userRepo := repository.NewUserRepository()
	userUC := usecase.NewUserUseCase(userRepo)
	userHandler := http.NewUserHandler(userUC)
	userHandler.RegisterRoutes(e)

	// Product setup
	productRepo := repository.NewProductRepository()
	productUC := usecase.NewProductUseCase(productRepo)
	productHandler := http.NewProductHandler(productUC)
	productHandler.RegisterRoutes(e)

	// Cart setup
	cartRepo := repository.NewCartRepository()
	cartUC := usecase.NewCartUseCase(cartRepo)
	cartHandler := http.NewCartHandler(cartUC)
	cartHandler.RegisterRoutes(e)

	log.Println("Server running on :8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
