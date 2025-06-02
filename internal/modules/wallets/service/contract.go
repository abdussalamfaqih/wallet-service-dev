package service

import (
	"context"

	"github.com/abdussalamfaqih/wallet-service-dev/internal/modules/wallets/presentations"
)

type Wallet interface {
	CreateAccount(ctx context.Context, req presentations.CreateAccount) error
	GetAccount(ctx context.Context, accountID string) (presentations.Account, error)
	SubmitTransaction(ctx context.Context, req presentations.CreateTransaction) error
}
