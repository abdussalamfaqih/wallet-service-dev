package http

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	accID, err := strconv.Atoi(accountID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponsePayload{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	if accID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponsePayload{
			Code:    http.StatusBadRequest,
			Message: "invalid accountID",
		})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	result, err := handler.ucase.GetAccount(ctx, accID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponsePayload{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	if result.AccountID == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ResponsePayload{
			Code:    http.StatusNotFound,
			Message: "DATA_NOT_FOUND",
		})

		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (handler *WalletHandler) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	var reqData presentations.CreateAccount

	ctx := r.Context()

	json.NewDecoder(r.Body).Decode(&reqData)
	if reqData.AccountID == 0 || reqData.InitialBalance == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponsePayload{
			Code:    http.StatusBadRequest,
			Message: "invalid request body",
		})
		return
	}

	err := handler.ucase.CreateAccount(ctx, reqData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponsePayload{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

func (handler *WalletHandler) CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var reqData presentations.CreateTransaction

	ctx := r.Context()

	json.NewDecoder(r.Body).Decode(&reqData)
	if reqData.SourceAccountID == 0 || reqData.DestinationAccountID == 0 || reqData.Amount == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponsePayload{
			Code:    http.StatusBadRequest,
			Message: "invalid request body",
		})
		return
	}

	err := handler.ucase.SubmitTransaction(ctx, reqData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponsePayload{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
}
