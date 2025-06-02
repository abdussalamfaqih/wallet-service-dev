package http

import (
	"encoding/json"
	"net/http"

	"github.com/abdussalamfaqih/wallet-service-dev/internal/modules/wallets/delivery/http/middlewares"
	"github.com/abdussalamfaqih/wallet-service-dev/internal/modules/wallets/presentations"
	"github.com/abdussalamfaqih/wallet-service-dev/internal/modules/wallets/service"
	"github.com/gorilla/mux"
)

type WalletHandler struct {
	ucase service.Wallet
}

func NewWalletHandler(r *mux.Router, ucase service.Wallet) {
	handler := &WalletHandler{
		ucase: ucase,
	}

	r.Use(middlewares.CommonMiddleware, middlewares.LoggingMiddleware)

	r.HandleFunc("/v1/accounts", handler.CreateAccountHandler).Methods(http.MethodPost)
	r.HandleFunc("/v1/accounts/{account_id}", handler.GetAccountHandler).Methods(http.MethodGet)
	r.HandleFunc("/v1/transactions", handler.CreateTransactionHandler).Methods(http.MethodPost)
}

func (handler *WalletHandler) GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	accountID := mux.Vars(r)["account_id"]

	w.Header().Add("Content-Type", "application/json")
	result, err := handler.ucase.GetAccount(ctx, accountID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponsePayload{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	if result.AccountID == "" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ResponsePayload{
			Code:    http.StatusNotFound,
			Message: "DATA_NOT_FOUND",
		})

		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponsePayload{
		Code:    http.StatusOK,
		Message: "SUCCESS",
		Data:    result,
	})
}

func (handler *WalletHandler) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	var reqData presentations.CreateAccount

	ctx := r.Context()

	json.NewDecoder(r.Body).Decode(&reqData)

	err := handler.ucase.CreateAccount(ctx, reqData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponsePayload{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ResponsePayload{
		Code:    http.StatusOK,
		Message: "SUCCESS",
	})
}

func (handler *WalletHandler) CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var reqData presentations.CreateTransaction

	ctx := r.Context()

	json.NewDecoder(r.Body).Decode(&reqData)

	err := handler.ucase.SubmitTransaction(ctx, reqData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponsePayload{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ResponsePayload{
		Code:    http.StatusOK,
		Message: "SUCCESS",
	})
}
