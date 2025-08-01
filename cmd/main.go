package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"
	"uhuaha/computers-management/internal/db"
	"uhuaha/computers-management/internal/db/postgres"
	"uhuaha/computers-management/internal/handler"
	"uhuaha/computers-management/internal/router"
	"uhuaha/computers-management/internal/service"

	"github.com/bdlm/log"
)

const PORT = ":8081"

func main() {
	dbConn, err := db.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	repository := postgres.NewRepository(dbConn)
	notifier := service.NewNotifier()
	computerMgmtService := service.NewComputerMgmtService(repository, notifier)
	handler := handler.New(computerMgmtService)
	router := router.New(handler)

	server := &http.Server{
		Addr:    PORT,
		Handler: router,
	}

	go func() {
		log.Info("Listening and serving on port " + PORT + " ...")

		err := http.ListenAndServe(PORT, router)
		if err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shut down server: %v", err)
	}
}
