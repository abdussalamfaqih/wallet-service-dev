package presentations

type (
	CreateAccount struct {
		AccountID      int    `json:"account_id"`
		InitialBalance string `json:"initial_balance"`
	}

	CreateTransaction struct {
		SourceAccountID      int    `json:"source_account_id"`
		DestinationAccountID int    `json:"destination_account_id"`
		Amount               string `json:"amount"`
	}

	Account struct {
		AccountID int    `json:"account_id"`
		Balance   string `json:"balance"`
	}
)
