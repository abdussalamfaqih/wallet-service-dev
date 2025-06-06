-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS ledger_entries (
    id UUID PRIMARY KEY,
    transaction_id UUID NOT NULL REFERENCES transactions(id),
    account_id INT NOT NULL REFERENCES accounts(account_id),
    amount DECIMAL(20, 6) NOT NULL,
    entry_type VARCHAR(50) NOT NULL, -- debit, credit
    balance_before DECIMAL(20, 6),
    balance_after DECIMAL(20, 6),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ledger_entries;
-- +goose StatementEnd 

