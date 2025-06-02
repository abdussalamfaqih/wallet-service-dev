package http

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/abdussalamfaqih/wallet-service-dev/internal/appconfig"
	"github.com/gorilla/mux"
)

func Start(ctx context.Context) {
	router := mux.NewRouter()

	cfg := appconfig.LoadConfig()

	RegisterHandlers(router, cfg)

	// Start the server
	log.Printf("Starting Server on port %s\n", cfg.App.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", cfg.App.Port), router))
}
