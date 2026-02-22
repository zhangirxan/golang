package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"golang/internal/handler"
	"golang/internal/middleware"
	"golang/internal/repository"
	_postgres "golang/internal/repository/_postgres"
	"golang/internal/usecase"
	"golang/pkg/modules"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConfig := initPostgreConfig()
	db := _postgres.NewPGXDialect(ctx, dbConfig)

	repos := repository.NewRepositories(db)
	userUsecase := usecase.NewUserUsecase(repos.UserRepository)
	userHandler := handler.NewUserHandler(userUsecase)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userHandler.GetUsers(w, r)
		case http.MethodPost:
			userHandler.CreateUser(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userHandler.GetUserByID(w, r)
		case http.MethodPut:
			userHandler.UpdateUser(w, r)
		case http.MethodDelete:
			userHandler.DeleteUser(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	chain := middleware.Logger(middleware.Auth(mux))

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", chain); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func initPostgreConfig() *modules.PostgreConfig {
	return &modules.PostgreConfig{
		Host:        "localhost",
		Port:        "5433",
		Username:    "postgres",
		Password:    "zhangir",
		DBName:      "mydb",
		SSLMode:     "disable",
		ExecTimeout: 5 * time.Second,
	}
}
