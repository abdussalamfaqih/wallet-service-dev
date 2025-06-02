package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/abdussalamfaqih/wallet-service-dev/internal/consts"
	"github.com/abdussalamfaqih/wallet-service-dev/internal/modules/wallets/presentations"
	"github.com/abdussalamfaqih/wallet-service-dev/internal/modules/wallets/repository"
)

type wallet struct {
	repo repository.WalletRepository
}

func NewWalletService(repo repository.WalletRepository) Wallet {
	return &wallet{
		repo: repo,
	}
}

func (s *wallet) GetAccount(ctx context.Context, accountID string) (repository.Account, error) {
	if err := validateAccountID(accountID); err != nil {
		return repository.Account{}, err
	}

	result, err := s.repo.GetAccount(ctx, accountID)
	if err != nil {
		return result, err
	}

	if result.ID == 0 {
		return result, errors.New("data not found")
	}

	return result, nil
}

func (s *wallet) CreateAccount(ctx context.Context, req presentations.CreateAccount) error {
	if err := validateCreateAccount(req); err != nil {
		return err
	}

	exist, err := s.repo.GetAccount(ctx, req.AccountID)
	if err != nil {
		return err
	}

	if exist.ID > 0 {
		return errors.New("data already exists")
	}

	return s.repo.CreateAccount(ctx, prepareDepositPayload(repository.Account{
		AccountID: req.AccountID,
		Balance:   req.Amount,
	},
	))
}

func (s *wallet) SubmitTransaction(ctx context.Context, req presentations.CreateTransaction) error {
	if err := validateAccountID(req.From); err != nil {
		return err
	}

	if err := validateAccountID(req.To); err != nil {
		return err
	}

	dataFrom, err := s.GetAccount(ctx, req.From)
	if err != nil {
		return err
	}

	dataTo, err := s.GetAccount(ctx, req.To)
	if err != nil {
		return err
	}

	if dataFrom.ID == 0 || dataTo.ID == 0 {
		return errors.New("data not found")
	}

	if err := validateAccounts(dataTo, dataFrom, req.Amount); err != nil {
		return err
	}

	payloadReq := prepareTrxPayload(dataFrom, dataTo, req.Amount)

	return s.repo.SubmitTransaction(ctx, payloadReq)
}

func prepareTrxPayload(from, to repository.Account, amount float64) repository.TransactionPayload {

	var payload repository.TransactionPayload

	from.Balance -= amount
	to.Balance += amount

	payload.From = from
	payload.To = to

	payload.Transaction = repository.Transaction{
		ID:              uuid.NewString(),
		ReferenceNumber: uuid.NewString(),
		Type:            consts.TransactionTypeTransfer,
		FromAccountID: sql.NullString{
			String: from.AccountID,
			Valid:  true,
		},
		ToAccountID: sql.NullString{
			String: to.AccountID,
			Valid:  true,
		},
		Amount:      amount,
		Status:      "completed",
		Description: fmt.Sprintf("transfer %v from %s to %s", amount, from.AccountID, to.AccountID),
		CreatedAt:   time.Now(),
	}

	payload.LedgerEntryFrom = repository.LedgerEntry{
		ID:            uuid.NewString(),
		TransactionID: payload.Transaction.ID,
		AccountID:     from.AccountID,
		EntryType:     consts.EntryTypeDebit,
		Amount:        amount,
		BalanceAfter:  from.Balance,
		BalanceBefore: from.Balance + amount,
		Description:   "transfer transaction",
		CreatedAt:     time.Now(),
	}

	payload.LedgerEntryTo = repository.LedgerEntry{
		ID:            uuid.NewString(),
		TransactionID: payload.Transaction.ID,
		AccountID:     to.AccountID,
		EntryType:     consts.EntryTypeCredit,
		Amount:        amount,
		BalanceAfter:  to.Balance,
		BalanceBefore: to.Balance - amount,
		Description:   "transfer transaction",
		CreatedAt:     time.Now(),
	}

	return payload
}

func prepareDepositPayload(acc repository.Account) repository.DepositPayload {
	var payload repository.DepositPayload

	payload.Account = acc
	payload.Transaction = repository.Transaction{
		ID:              uuid.NewString(),
		ReferenceNumber: uuid.NewString(),
		Type:            "deposit",
		ToAccountID: sql.NullString{
			String: acc.AccountID,
			Valid:  true,
		},
		Amount:      acc.Balance,
		Status:      "completed",
		Description: fmt.Sprintf("deposit %v to %s", acc.Balance, acc.AccountID),
		CreatedAt:   time.Now(),
	}

	payload.LedgerEntry = repository.LedgerEntry{
		ID:            uuid.NewString(),
		TransactionID: payload.Transaction.ID,
		AccountID:     acc.AccountID,
		EntryType:     consts.EntryTypeCredit,
		Amount:        acc.Balance,
		BalanceAfter:  acc.Balance,
		BalanceBefore: 0,
		Description:   "deposit transaction",
		CreatedAt:     time.Now(),
	}

	return payload
}
