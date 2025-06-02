package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/abdussalamfaqih/wallet-service-dev/internal/consts"
)

type WalletRepository interface {
	CreateAccount(ctx context.Context, payload DepositPayload) error
	GetAccount(ctx context.Context, accountID string) (Account, error)
	SubmitTransaction(ctx context.Context, payload TransactionPayload) error
}

type (
	Account struct {
		ID        int       `json:"id"`
		AccountID string    `json:"account_id"`
		Currency  string    `json:"currency"`
		Status    string    `json:"status"`
		Balance   float64   `json:"balance"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	Transaction struct {
		ID              string                 `json:"id"`
		ReferenceNumber string                 `json:"reference_number"`
		Type            consts.TransactionType `json:"type"`
		Description     string                 `json:"description,omitempty"`
		FromAccountID   sql.NullString         `json:"from_account_id,omitempty"`
		ToAccountID     sql.NullString         `json:"to_account_id,omitempty"`
		Amount          float64                `json:"amount"`
		Metadata        json.RawMessage        `json:"metadata,omitempty"`
		Status          string                 `json:"status"`
		CreatedAt       time.Time              `json:"created_at"`
	}

	LedgerEntry struct {
		ID            string           `json:"id"`
		TransactionID string           `json:"transaction_id"`
		AccountID     string           `json:"account_id"`
		Amount        float64          `json:"amount"`
		EntryType     consts.EntryType `json:"entry_type"`
		BalanceBefore float64          `json:"balance_before,omitempty"`
		BalanceAfter  float64          `json:"balance_after,omitempty"`
		Description   string           `json:"description,omitempty"`
		CreatedAt     time.Time        `json:"created_at"`
	}

	DepositPayload struct {
		Account     Account
		Transaction Transaction
		LedgerEntry LedgerEntry
	}

	TransactionPayload struct {
		From            Account
		To              Account
		Amount          float64
		Transaction     Transaction
		LedgerEntryFrom LedgerEntry
		LedgerEntryTo   LedgerEntry
	}
)
