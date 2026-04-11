package app

import (
	"log"
	"os"
	"practice-7/internal/controller/http/v1"
	"practice-7/internal/usecase"
	"practice-7/internal/usecase/repo"
	"practice-7/pkg/logger"
	"practice-7/pkg/postgres"

	"github.com/gin-gonic/gin"
)

func Run() {
	l := logger.New("info")

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "host=localhost user=postgres password=zhangir dbname=auth_service port=5433 sslmode=disable"
	}

	pg, err := postgres.New(dbURL)
	if err != nil {
		log.Fatalf("postgres.New: %v", err)
	}

	userRepo := repo.NewUserRepo(pg)
	userUseCase := usecase.NewUserUseCase(userRepo)

	router := gin.Default()
	v1.NewRouter(router, userUseCase, l)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	log.Printf("Server running on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("router.Run: %v", err)
	}
}
