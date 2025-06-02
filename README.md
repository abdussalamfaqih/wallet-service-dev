# REST Wallet Service

## Features
- API Create Account
- API Get Account Balance
- API Create Transactions

## Preparations
1. Have Golang with minimum version of 1.21
2. Have Goose migration tools available as its used for migration purpose, reference: [link](https://github.com/pressly/goose)

## Run The Service
To run in containerized environments:
```bash
docker-compose -f deployment/docker-compose.yml --project-directory . up 
```

Execute Migration scripts from root directory, example script using static values:
```bash
goose -dir=db/migrations postgres "user=postgres password=strong_password dbname=wallet_db sslmode=disable" up
```

To run in local, use command: 
```bash
go run main.go run-http
```
and ensure the postgres is available, can use the compose infra and run the migration after ready:
```bash
docker-compose -f deployment/docker-compose-db.yml up
```

### Available API
#### Create Account
1. account_id := uuid text format
2. amount := float numbers with 6 precision digit, cannot be negative number

```bash
curl --location 'http://localhost:8080/v1/accounts' \
--header 'Content-Type: application/json' \
--data '{
    "account_id": "3128e237-fbd1-4271-9e40-17b132db5859", 
    "amount": 100
}'
```

#### Get Account
1. account_id := uuid text format


```bash
curl --location 'http://localhost:8080/v1/accounts/{account_id}'
```

#### Create Transaction
1. from := uuid text format for sender account_id
2. to := uuid text format for receiver account_id 
3. amount := float numbers with 6 precision digit, cannot be negative number


```bash
curl --location 'http://localhost:8080/v1/transactions' \
--header 'Content-Type: application/json' \
--data '{
    "from": "3128e237-fbd1-4271-9e40-17b132db5859",
    "to": "c5934062-e368-4bc0-be95-c2265bb7430f",
    "amount": 20
}'
```

### Project Structure
```bash
├── README.md
├── cmd
│   ├── http
│   │   ├── http.go
│   │   └── wallet.go
│   └── root.go
├── config
│   ├── config.json
│   └── config.json.example
├── db
│   └── migrations
│       ├── 20250601150631_add_table_accounts.sql
│       ├── 20250601150641_add_table_transactions.sql
│       └── 20250601181846_add_table_ledger_entries.sql
├── deployment
│   ├── Dockerfile
│   ├── docker-compose-db.yml
│   └── docker-compose.yml
├── go.mod
├── go.sum
├── internal
│   ├── appconfig
│   │   └── config.go
│   ├── bootstrap
│   │   └── db.go
│   ├── consts
│   │   ├── common.go
│   │   └── wallet.go
│   ├── modules
│   │   └── wallets
│   │       ├── delivery
│   │       │   └── http
│   │       │       ├── middlewares
│   │       │       │   └── common.go
│   │       │       ├── response.go
│   │       │       └── wallet.go
│   │       ├── presentations
│   │       │   └── wallet.go
│   │       ├── repository
│   │       │   ├── contract.go
│   │       │   └── wallet.go
│   │       └── service
│   │           ├── contract.go
│   │           ├── validations.go
│   │           └── wallet.go
│   └── utils
│       └── decimal.go
├── main.go
└── pkg
    ├── config
    │   └── viper.go
    ├── db
    │   ├── contract.go
    │   └── postgres.go
    ├── decimal
    └── logger
        └── zap.go
```

