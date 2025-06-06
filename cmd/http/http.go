package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abdussalamfaqih/wallet-service-dev/internal/appconfig"
	"github.com/gorilla/mux"
)

func Start(ctx context.Context, cfg appconfig.Config) {
	router := mux.NewRouter()

	RegisterHandlers(router, cfg)

	// Start the server
	startServer(cfg, router)
}

func startServer(cfg appconfig.Config, router http.Handler) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.App.Port),
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Starting Server on port %s\n", cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)

	// Catch SIGINT, SIGTERM
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a context with timeout to finish requests
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
