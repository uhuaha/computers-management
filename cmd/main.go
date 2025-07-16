package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"uhuaha/computers-management/internal/handler"
	"uhuaha/computers-management/internal/router"
)

const PORT = ":8080"

func main() {
	// TODO:
	// dbConnection := ...
	// repository := postgres.NewComputerRepository(dbConnection)
	// service := service.NewComputerMgmtService(repository)
	handler := handler.New() // gets the service
	router := router.New(handler)

	server := &http.Server{
		Addr:    PORT,
		Handler: router,
	}

	go func() {
		log.Println("Listening and serving on port 8080...")
		
		err := http.ListenAndServe(PORT, router)
		if err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shut down server: %v", err)
	}
}
