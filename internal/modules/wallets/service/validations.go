package service

import (
	"errors"

	"github.com/abdussalamfaqih/wallet-service-dev/internal/modules/wallets/presentations"
	"github.com/abdussalamfaqih/wallet-service-dev/internal/modules/wallets/repository"
	"github.com/google/uuid"
)

func validateAccounts(to, from repository.Account, amount float64) error {
	if from.Balance < amount {
		return errors.New("sender balance less than amount")
	}

	if to.AccountID == from.AccountID {
		return errors.New("request payload invalid")
	}

	if amount < 1 {
		return errors.New("amount cannot be less than 1.00")
	}

	if amount > 1000000 {
		return errors.New("amount cannot be greater than than 1000000.00")
	}

	return nil
}

func validateCreateAccount(p presentations.CreateAccount) error {
	if err := validateAccountID(p.AccountID); err != nil {
		return err
	}

	if p.Amount < 1 {
		return errors.New("amount cannot be less than 1.00")
	}

	if p.Amount > 1000000 {
		return errors.New("amount cannot be greater than than 1000000.00")
	}

	return nil
}

func validateAccountID(s string) error {
	_, err := uuid.Parse(s)
	if err != nil {
		return errors.New("invalid account_id format")
	}

	return nil
}
