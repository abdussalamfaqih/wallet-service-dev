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
		req    int
		result presentations.Account
	}{
		{
			name:   "FAILED validation",
			err:    errors.New("invalid account_id format"),
			req:    0,
			result: presentations.Account{},
			mock:   func() {},
		},
		{
			name:   "FAILED db error",
			err:    errTest,
			req:    2,
			result: presentations.Account{},
			mock: func() {
				mRepo.EXPECT().GetAccount(ctx, 2).Return(repository.Account{}, errTest)
			},
		},
		{
			name:   "FAILED data not found",
			err:    errors.New("data not found"),
			req:    2,
			result: presentations.Account{},
			mock: func() {
				mRepo.EXPECT().GetAccount(ctx, 2).Return(repository.Account{}, nil)
			},
		},
		{
			name: "SUCCESS",
			err:  nil,
			req:  2,
			result: presentations.Account{
				AccountID: 2,
				Balance:   "0",
			},
			mock: func() {
				mRepo.EXPECT().GetAccount(ctx, 2).Return(repository.Account{
					ID:        1,
					AccountID: 2,
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
				AccountID:      0,
				InitialBalance: "0",
			},
			mock: func() {},
		},
		{
			name: "FAILED validation #2",
			err:  errors.New("amount cannot be less than 1.00"),
			req: presentations.CreateAccount{
				AccountID:      2,
				InitialBalance: "0",
			},
			mock: func() {},
		},
		{
			name: "FAILED get db error",
			err:  errTest,
			req: presentations.CreateAccount{
				AccountID:      2,
				InitialBalance: "100",
			},
			mock: func() {
				mRepo.EXPECT().GetAccount(ctx, 2).Return(repository.Account{}, errTest)
			},
		},
		{
			name: "FAILED data exists",
			err:  errors.New("data already exists"),
			req: presentations.CreateAccount{
				AccountID:      2,
				InitialBalance: "100",
			},
			mock: func() {
				mRepo.EXPECT().GetAccount(ctx, 2).Return(repository.Account{
					ID: 1,
				}, nil)
			},
		},
		{
			name: "FAILED insert error",
			err:  errTest,
			req: presentations.CreateAccount{
				AccountID:      2,
				InitialBalance: "100",
			},
			mock: func() {
				mRepo.EXPECT().GetAccount(ctx, 2).Return(repository.Account{}, nil)

				mRepo.EXPECT().CreateAccount(ctx, gomock.AssignableToTypeOf(repository.DepositPayload{})).Return(errTest)
			},
		},
		{
			name: "SUCCESS",
			err:  nil,
			req: presentations.CreateAccount{
				AccountID:      2,
				InitialBalance: "100",
			},
			mock: func() {
				mRepo.EXPECT().GetAccount(ctx, 2).Return(repository.Account{}, nil)

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
				DestinationAccountID: 0,
				SourceAccountID:      0,
				Amount:               "0",
			},
			mock: func() {},
		},

		{
			name: "FAILED validation #2",
			err:  errors.New("request payload invalid"),
			req: presentations.CreateTransaction{
				DestinationAccountID: 1,
				SourceAccountID:      1,
				Amount:               "0",
			},
			mock: func() {},
		},
		{
			name: "FAILED get db #1",
			err:  errTest,
			req: presentations.CreateTransaction{
				DestinationAccountID: 2,
				SourceAccountID:      3,
				Amount:               "10",
			},
			mock: func() {
				mRepo.EXPECT().GetAccount(ctx, 3).Return(repository.Account{}, errTest)
			},
		},
		{
			name: "FAILED get db #1",
			err:  errTest,
			req: presentations.CreateTransaction{
				DestinationAccountID: 2,
				SourceAccountID:      3,
				Amount:               "10",
			},
			mock: func() {
				mRepo.EXPECT().GetAccount(ctx, 3).Return(repository.Account{
					ID: 1,
				}, nil)
				mRepo.EXPECT().GetAccount(ctx, 2).Return(repository.Account{}, errTest)
			},
		},
		{
			name: "FAILED data not found",
			err:  errors.New("data not found"),
			req: presentations.CreateTransaction{
				DestinationAccountID: 2,
				SourceAccountID:      3,
				Amount:               "10",
			},
			mock: func() {
				mRepo.EXPECT().GetAccount(ctx, 3).Return(repository.Account{
					ID: 1,
				}, nil)
				mRepo.EXPECT().GetAccount(ctx, 2).Return(repository.Account{}, nil)
			},
		},
		{
			name: "Error Insert",
			err:  errTest,
			req: presentations.CreateTransaction{
				DestinationAccountID: 2,
				SourceAccountID:      3,
				Amount:               "10",
			},
			mock: func() {
				mRepo.EXPECT().GetAccount(ctx, 3).Return(repository.Account{
					ID:        1,
					Balance:   50,
					AccountID: 3,
				}, nil)
				mRepo.EXPECT().GetAccount(ctx, 2).Return(repository.Account{
					ID:        2,
					AccountID: 2,
					Balance:   50,
				}, nil)

				mRepo.EXPECT().SubmitTransaction(ctx, gomock.AssignableToTypeOf(repository.TransactionPayload{})).Return(errTest)
			},
		},
		{
			name: "SUCCESS",
			err:  nil,
			req: presentations.CreateTransaction{
				DestinationAccountID: 2,
				SourceAccountID:      3,
				Amount:               "10",
			},
			mock: func() {
				mRepo.EXPECT().GetAccount(ctx, 3).Return(repository.Account{
					ID:        1,
					Balance:   50,
					AccountID: 3,
				}, nil)
				mRepo.EXPECT().GetAccount(ctx, 2).Return(repository.Account{
					ID:        2,
					AccountID: 2,
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
