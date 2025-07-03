package main

import (
	"log"

	"backend/internal/delivery/http"
	"backend/internal/repository"
	"backend/internal/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Allow requests from frontend dev server (5173)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{echo.GET},
	}))

	repo := repository.NewUserRepository()
	uc := usecase.NewUserUseCase(repo)
	handler := http.NewHandler(uc)

	handler.RegisterRoutes(e)

	log.Println("Server running on :8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
