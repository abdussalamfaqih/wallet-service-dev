package service

import (
	"errors"
	"strconv"

	"github.com/abdussalamfaqih/wallet-service-dev/internal/modules/wallets/presentations"
	"github.com/abdussalamfaqih/wallet-service-dev/internal/modules/wallets/repository"
)

func validateAccounts(to, from repository.Account, amount float64) error {
	if from.Balance < amount {
		return errors.New("sender balance less than amount")
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

	amount, err := strconv.ParseFloat(p.InitialBalance, 64)
	if err != nil {
		return err
	}
	if amount < 1 {
		return errors.New("amount cannot be less than 1.00")
	}

	if amount > 1000000 {
		return errors.New("amount cannot be greater than than 1000000.00")
	}

	return nil
}

func validateAccountID(id int) error {
	if id < 1 {
		return errors.New("invalid account_id format")
	}

	return nil
}
