package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/abdussalamfaqih/wallet-service-dev/pkg/db"
)

type walletRepo struct {
	db *db.Repository
}

func NewWalletRepository(db *db.Repository) WalletRepository {
	return &walletRepo{
		db: db,
	}
}

func (r *walletRepo) CreateAccount(ctx context.Context, payload DepositPayload) error {
	return r.db.WithTransaction(ctx, func(ctx context.Context, repo *db.Repository) error {
		_, err := repo.Exec(ctx,
			"INSERT INTO accounts (account_id, currency, balance) VALUES ($1, $2, $3)",
			payload.Account.AccountID, "USD", payload.Account.Balance,
		)
		if err != nil {
			return fmt.Errorf("failed to insert accounts data: %w", err)
		}

		_, err = repo.Exec(ctx, `INSERT INTO transactions(id,
							reference_number,
							type,
							from_account_id,
							to_account_id,
							amount,
							description,
							status,
							created_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			payload.Transaction.ID,
			payload.Transaction.ReferenceNumber,
			payload.Transaction.Type,
			payload.Transaction.FromAccountID,
			payload.Transaction.ToAccountID,
			payload.Transaction.Amount,
			payload.Transaction.Description,
			payload.Transaction.Status,
			payload.Transaction.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert transactions data: %w", err)
		}

		_, err = repo.Exec(
			ctx, ` INSERT INTO ledger_entries 
				(
					id,
					transaction_id,
					account_id,
					entry_type,
					amount,
					balance_before,
					balance_after,
					description,
					created_at
				)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			payload.LedgerEntry.ID,
			payload.LedgerEntry.TransactionID,
			payload.LedgerEntry.AccountID,
			payload.LedgerEntry.EntryType,
			payload.LedgerEntry.Amount,
			payload.LedgerEntry.BalanceBefore,
			payload.LedgerEntry.BalanceAfter,
			payload.LedgerEntry.Description,
			payload.LedgerEntry.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert ledger_entries data: %w", err)
		}

		return nil
	})
}

func (r *walletRepo) GetAccount(ctx context.Context, accountID string) (Account, error) {
	var result Account
	err := r.db.QueryRow(ctx, "SELECT id, account_id, status, balance, created_at, updated_at FROM accounts WHERE account_id = $1", accountID).Scan(&result.ID, &result.AccountID, &result.Status, &result.Balance, &result.CreatedAt, &result.UpdatedAt)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return result, nil
	}

	if err != nil {
		return result, fmt.Errorf("failed to query data: %w", err)
	}

	return result, nil
}

func (r *walletRepo) SubmitTransaction(ctx context.Context, payload TransactionPayload) error {
	return r.db.WithTransaction(ctx, func(ctx context.Context, repo *db.Repository) error {

		var idFrom int
		err := repo.QueryRow(ctx, "SELECT id from accounts WHERE account_id = $1 FOR UPDATE", payload.From.AccountID).Scan(&idFrom)
		if err != nil {
			return fmt.Errorf("failed to query  from data: %w", err)
		}

		var idTo int
		err = repo.QueryRow(ctx, "SELECT id from accounts WHERE account_id = $1 FOR UPDATE", payload.To.AccountID).Scan(&idTo)
		if err != nil {
			return fmt.Errorf("failed to query to data: %w", err)
		}

		if idFrom == 0 || idTo == 0 {
			return errors.New("data not found")
		}

		_, err = repo.Exec(ctx, `INSERT INTO transactions (id,
							reference_number,
							type,
							from_account_id,
							to_account_id,
							amount,
							description,
							status,
							created_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			payload.Transaction.ID,
			payload.Transaction.ReferenceNumber,
			payload.Transaction.Type,
			payload.Transaction.FromAccountID,
			payload.Transaction.ToAccountID,
			payload.Transaction.Amount,
			payload.Transaction.Description,
			payload.Transaction.Status,
			payload.Transaction.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert transactions data: %w", err)
		}

		_, err = repo.Exec(
			ctx, ` INSERT INTO ledger_entries 
				(
					id,
					transaction_id,
					account_id,
					entry_type,
					amount,
					balance_before,
					balance_after,
					description,
					created_at
				)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9), 
					($10, $11, $12, $13, $14, $15, $16, $17, $18)`,
			payload.LedgerEntryFrom.ID,
			payload.LedgerEntryFrom.TransactionID,
			payload.LedgerEntryFrom.AccountID,
			payload.LedgerEntryFrom.EntryType,
			payload.LedgerEntryFrom.Amount,
			payload.LedgerEntryFrom.BalanceBefore,
			payload.LedgerEntryFrom.BalanceAfter,
			payload.LedgerEntryFrom.Description,
			payload.LedgerEntryFrom.CreatedAt,
			payload.LedgerEntryTo.ID,
			payload.LedgerEntryTo.TransactionID,
			payload.LedgerEntryTo.AccountID,
			payload.LedgerEntryTo.EntryType,
			payload.LedgerEntryTo.Amount,
			payload.LedgerEntryTo.BalanceBefore,
			payload.LedgerEntryTo.BalanceAfter,
			payload.LedgerEntryTo.Description,
			payload.LedgerEntryTo.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert ledger_entries data: %w", err)
		}

		_, err = repo.Exec(ctx, `UPDATE accounts SET balance = $1, updated_at = $2
		WHERE account_id = $3`, payload.From.Balance, time.Now(), payload.From.AccountID)
		if err != nil {
			return fmt.Errorf("failed to update accounts data from: %w", err)
		}

		_, err = repo.Exec(ctx, `UPDATE accounts SET balance = $1, updated_at = $2
		WHERE account_id = $3`, payload.To.Balance, time.Now(), payload.To.AccountID)
		if err != nil {
			return fmt.Errorf("failed to update accounts data to: %w", err)
		}

		return nil
	})
}
