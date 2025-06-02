package http

import (
	"github.com/abdussalamfaqih/wallet-service-dev/internal/appconfig"
	"github.com/abdussalamfaqih/wallet-service-dev/internal/bootstrap"
	httpDelivery "github.com/abdussalamfaqih/wallet-service-dev/internal/modules/wallets/delivery/http"
	"github.com/abdussalamfaqih/wallet-service-dev/internal/modules/wallets/repository"
	"github.com/abdussalamfaqih/wallet-service-dev/internal/modules/wallets/service"
	"github.com/gorilla/mux"
)

func RegisterHandlers(r *mux.Router, cfg appconfig.Config) {

	dbClient := bootstrap.NewDB(cfg.Database)
	repo := repository.NewWalletRepository(dbClient)
	service := service.NewWalletService(repo)

	httpDelivery.NewWalletHandler(r, service)
}
