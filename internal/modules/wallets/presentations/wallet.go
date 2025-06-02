package presentations

type (
	CreateAccount struct {
		AccountID string  `json:"account_id"`
		Amount    float64 `json:"amount"`
	}

	CreateTransaction struct {
		From   string  `json:"from"`
		To     string  `json:"to"`
		Amount float64 `json:"amount"`
	}

	Account struct {
		AccountID string  `json:"account_id"`
		Balance   float64 `json:"balance"`
	}
)
