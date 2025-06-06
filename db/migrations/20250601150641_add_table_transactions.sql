-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY,
    reference_number VARCHAR(255) UNIQUE NOT NULL,
    type VARCHAR(50) NOT NULL, -- e.g., transfer, deposit, withdraw
    description TEXT,
    from_account_id INT NULL REFERENCES accounts(account_id), 
    to_account_id INT NULL REFERENCES accounts(account_id), 
    amount DECIMAL(20, 6) NOT NULL, -- always positive
    metadata JSONB,
    status VARCHAR(20) DEFAULT 'completed',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;
-- +goose StatementEnd
