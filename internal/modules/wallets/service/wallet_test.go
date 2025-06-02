package service

import (
	"context"
	"errors"
	"testing"

	"github.com/abdussalamfaqih/wallet-service-dev/internal/modules/wallets/presentations"
	"github.com/abdussalamfaqih/wallet-service-dev/internal/modules/wallets/repository"
	"github.com/abdussalamfaqih/wallet-service-dev/internal/modules/wallets/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetAccount(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mRepo := mocks.NewMockWalletRepository(mockCtl)

	errTest := errors.New("test")
	ctx := context.TODO()

	testTables := []struct {
		name   string
		mock   func()
		err    error
		req    string
		result presentations.Account
	}{
		{
			name:   "FAILED validation",
			err:    errors.New("invalid account_id format"),
			req:    "accountID",
			result: presentations.Account{},
			mock:   func() {},
		},
		{
			name:   "FAILED db error",
			err:    errTest,
			req:    "c5934062-e368-4bc0-be95-c2265bb7430f",
			result: presentations.Account{},
			mock: func() {
				mRepo.EXPECT().GetAccount(ctx, "c5934062-e368-4bc0-be95-c2265bb7430f").Return(repository.Account{}, errTest)
			},
		},
		{
			name:   "FAILED data not found",
			err:    errors.New("data not found"),
			req:    "c5934062-e368-4bc0-be95-c2265bb7430f",
			result: presentations.Account{},
			mock: func() {
				mRepo.EXPECT().GetAccount(ctx, "c5934062-e368-4bc0-be95-c2265bb7430f").Return(repository.Account{}, nil)
			},
		},
		{
			name: "SUCCESS",
			err:  nil,
			req:  "c5934062-e368-4bc0-be95-c2265bb7430f",
			result: presentations.Account{
				AccountID: "c5934062-e368-4bc0-be95-c2265bb7430f",
			},
			mock: func() {
				mRepo.EXPECT().GetAccount(ctx, "c5934062-e368-4bc0-be95-c2265bb7430f").Return(repository.Account{
					ID:        1,
					AccountID: "c5934062-e368-4bc0-be95-c2265bb7430f",
				}, nil)
			},
		},
	}

	svc := NewWalletService(mRepo)
	for _, tt := range testTables {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			resp, err := svc.GetAccount(ctx, tt.req)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.result, resp)
		})
	}
}

func TestCreateAccount(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mRepo := mocks.NewMockWalletRepository(mockCtl)

	errTest := errors.New("test")
	ctx := context.TODO()

	testTables := []struct {
		name string
		mock func()
		err  error
		req  presentations.CreateAccount
	}{
		{
			name: "FAILED validation #1",
			err:  errors.New("invalid account_id format"),
			req: presentations.CreateAccount{
				AccountID: "accountID",
			},
			mock: func() {},
		},
		{
			name: "FAILED validation #2",
			err:  errors.New("amount cannot be less than 1.00"),
			req: presentations.CreateAccount{
				AccountID: "c5934062-e368-4bc0-be95-c2265bb7430f",
			},
			mock: func() {},
		},
		{
			name: "FAILED get db error",
			err:  errTest,
			req: presentations.CreateAccount{
				AccountID: "c5934062-e368-4bc0-be95-c2265bb7430f",
				Amount:    100,
			},
			mock: func() {
				mRepo.EXPECT().GetAccount(ctx, "c5934062-e368-4bc0-be95-c2265bb7430f").Return(repository.Account{}, errTest)
			},
		},
		{
			name: "FAILED data exists",
			err:  errors.New("data already exists"),
			req: presentations.CreateAccount{
				AccountID: "c5934062-e368-4bc0-be95-c2265bb7430f",
				Amount:    100,
			},
			mock: func() {
				mRepo.EXPECT().GetAccount(ctx, "c5934062-e368-4bc0-be95-c2265bb7430f").Return(repository.Account{
					ID: 1,
				}, nil)
			},
		},
		{
			name: "FAILED insert error",
			err:  errTest,
			req: presentations.CreateAccount{
				AccountID: "c5934062-e368-4bc0-be95-c2265bb7430f",
				Amount:    100,
			},
			mock: func() {
				mRepo.EXPECT().GetAccount(ctx, "c5934062-e368-4bc0-be95-c2265bb7430f").Return(repository.Account{}, nil)

				mRepo.EXPECT().CreateAccount(ctx, gomock.AssignableToTypeOf(repository.DepositPayload{})).Return(errTest)
			},
		},
		{
			name: "SUCCESS",
			err:  nil,
			req: presentations.CreateAccount{
				AccountID: "c5934062-e368-4bc0-be95-c2265bb7430f",
				Amount:    100,
			},
			mock: func() {
				mRepo.EXPECT().GetAccount(ctx, "c5934062-e368-4bc0-be95-c2265bb7430f").Return(repository.Account{}, nil)

				mRepo.EXPECT().CreateAccount(ctx, gomock.AssignableToTypeOf(repository.DepositPayload{})).Return(nil)
			},
		},
	}

	svc := NewWalletService(mRepo)
	for _, tt := range testTables {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := svc.CreateAccount(ctx, tt.req)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestSubmitTransaction(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mRepo := mocks.NewMockWalletRepository(mockCtl)

	errTest := errors.New("test")
	ctx := context.TODO()

	testTables := []struct {
		name string
		mock func()
		err  error
		req  presentations.CreateTransaction
	}{
		{
			name: "FAILED validation #1",
			err:  errors.New("invalid account_id format"),
			req: presentations.CreateTransaction{
				To:   "accountID",
				From: "accountID",
			},
			mock: func() {},
		},
		{
			name: "FAILED get db #1",
			err:  errTest,
			req: presentations.CreateTransaction{
				To:     "3128e237-fbd1-4271-9e40-17b132db5859",
				From:   "c5934062-e368-4bc0-be95-c2265bb7430f",
				Amount: 10,
			},
			mock: func() {
				mRepo.EXPECT().GetAccount(ctx, "c5934062-e368-4bc0-be95-c2265bb7430f").Return(repository.Account{}, errTest)
			},
		},
		{
			name: "FAILED get db #1",
			err:  errTest,
			req: presentations.CreateTransaction{
				To:     "3128e237-fbd1-4271-9e40-17b132db5859",
				From:   "c5934062-e368-4bc0-be95-c2265bb7430f",
				Amount: 10,
			},
			mock: func() {
				mRepo.EXPECT().GetAccount(ctx, "c5934062-e368-4bc0-be95-c2265bb7430f").Return(repository.Account{
					ID: 1,
				}, nil)
				mRepo.EXPECT().GetAccount(ctx, "3128e237-fbd1-4271-9e40-17b132db5859").Return(repository.Account{}, errTest)
			},
		},
		{
			name: "SUCCESS",
			err:  nil,
			req: presentations.CreateTransaction{
				To:     "3128e237-fbd1-4271-9e40-17b132db5859",
				From:   "c5934062-e368-4bc0-be95-c2265bb7430f",
				Amount: 10,
			},
			mock: func() {
				mRepo.EXPECT().GetAccount(ctx, "c5934062-e368-4bc0-be95-c2265bb7430f").Return(repository.Account{
					ID:        1,
					Balance:   50,
					AccountID: "c5934062-e368-4bc0-be95-c2265bb7430f",
				}, nil)
				mRepo.EXPECT().GetAccount(ctx, "3128e237-fbd1-4271-9e40-17b132db5859").Return(repository.Account{
					ID:        2,
					AccountID: "3128e237-fbd1-4271-9e40-17b132db5859",
					Balance:   50,
				}, nil)

				mRepo.EXPECT().SubmitTransaction(ctx, gomock.AssignableToTypeOf(repository.TransactionPayload{})).Return(nil)
			},
		},
	}

	svc := NewWalletService(mRepo)
	for _, tt := range testTables {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := svc.SubmitTransaction(ctx, tt.req)
			assert.Equal(t, tt.err, err)
		})
	}
}
