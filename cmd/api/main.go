package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"practice2/internal/handlers"       
	"practice2/internal/middleware"   
	"practice2/internal/storage"       
)

func main() {

	store := storage.NewTaskStorage()

	//handlers
	handler := handlers.NewTaskHandler(store)

	//routes
	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", handler.HandleTasks)

	//middleware
	wrappedMux := middleware.Logging(middleware.APIKey(mux))

	//serv config
	server := &http.Server{
		Addr:         ":8080",
		Handler:      wrappedMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	//shutdown
	go func() {
		log.Println("Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
