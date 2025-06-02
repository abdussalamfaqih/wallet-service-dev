package presentations

type CreateAccount struct {
	AccountID string  `json:"account_id"`
	Amount    float64 `json:"amount"`
}

type CreateTransaction struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}
