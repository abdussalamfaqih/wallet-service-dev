# REST Wallet Service

## Features
- API Create Account
- API Get Account Balance
- API Create Transactions

## Preparations
1. Have Golang with minimum version of 1.24
2. [Optional] Have Makefile installed to run below command, if unable, copy the command directly in terminal to run each command
3. [Optional] Have Goose migration tools available as its used for migration purpose, reference: [link](https://github.com/pressly/goose)

## Run The Service
To run in containerized environments:
```bash
make run-service 
```

Execute Migration scripts from root directory, 
Run the go command to execute migration,
```bash
make run-migration
```

or if you have goose cmd, this is example script using static values:
```bash
goose -dir=db/migrations postgres "user=postgres password=strong_password dbname=wallet_db sslmode=disable" up
```

To run in local, use command: 
```bash
make run-local
```
and ensure the postgres is available, can use the compose infra and run the migration after ready:
```bash
make run-db
```
> [!NOTE] 
> To ensure the connection to db works, check the `host` in `config/config.json` file using the valid value
> - for docker-compose env, host: `postgres`
> - for running locally (service or migration), can use host: `localhost`

### Available API
#### Create Account
1. account_id := uuid text format
2. amount := float numbers with 6 precision digit, cannot be negative number

```bash
curl --location 'http://localhost:8080/v1/accounts' \
--header 'Content-Type: application/json' \
--data '{
    "account_id": 456,
    "initial_balance": "200.23344"
}'
```

#### Get Account
1. account_id := uuid text format


```bash
curl --location 'http://localhost:8080/v1/accounts/456'
```

#### Create Transaction
1. from := uuid text format for sender account_id
2. to := uuid text format for receiver account_id 
3. amount := float numbers with 6 precision digit, cannot be negative number


```bash
curl --location 'http://localhost:8080/v1/transactions' \
--header 'Content-Type: application/json' \
--data '{
    "source_account_id": 123,
    "destination_account_id": 456,
    "amount": "100.12345"
}'
```
