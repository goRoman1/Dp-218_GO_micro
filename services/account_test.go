package services

import (
	"Dp218GO/models"
	"Dp218GO/services/mock"
	"errors"
	"github.com/golang/mock/gomock"
	assert "github.com/stretchr/testify/require"
	"testing"
	"time"
)

// UseCasesMock is a struct which exists of repositories which are mocked and our service.
type accountUseCasesMock struct {
	AccountServiceUC       *AccountService
	RepoPaymentType        *mock.MockPaymentTypeRepo
	RepoAccountTransaction *mock.MockAccountTransactionRepo
	RepoAccount            *mock.MockAccountRepo
	Clock                  *mock.MockClock
}

type accountTestCase struct {
	name string
	test func(t *testing.T, mock *accountUseCasesMock)
}

// We can keep this function without changes in our next test-cases. Except of 'mock' declaration.
func runTestCases(t *testing.T, testCases []accountTestCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					tt.Error(err)
				}
			}()

			ctrl := gomock.NewController(tt)
			defer ctrl.Finish()

			// Here we should change if our struct name will be different to 'UseCasesMock'.
			mock := newAccountUseCasesMock(ctrl)

			tc.test(tt, mock)
		})
	}
}

func newAccountUseCasesMock(ctrl *gomock.Controller) *accountUseCasesMock {
	repoAccount := mock.NewMockAccountRepo(ctrl)
	repoAccountTransaction := mock.NewMockAccountTransactionRepo(ctrl)
	repoPaymentType := mock.NewMockPaymentTypeRepo(ctrl)
	clock := mock.NewMockClock(ctrl)

	// We created 'clock' for mocking 'time.Now()'
	// Transfer 'clock' here just because it doesn't work in any other way.
	accountServiceUC := NewAccountService(repoAccount, repoAccountTransaction, repoPaymentType, clock)

	return &accountUseCasesMock{
		AccountServiceUC:       accountServiceUC,
		RepoPaymentType:        repoPaymentType,
		RepoAccountTransaction: repoAccountTransaction,
		RepoAccount:            repoAccount,
		Clock:                  clock,
	}
}

func Test_Account_AddMoneyToAccount(t *testing.T) {
	runTestCases(t, []accountTestCase{
		{ // In this case we are going by happy path.
			name: "Correct",
			test: func(t *testing.T, mock *accountUseCasesMock) {

				// Create a variable with the exact time for mocking time.Now().
				var currentTime = time.Date(2021, 12, 19, 12, 21, 00, 00, time.UTC)

				// With help of mocks we can call the functions of repositories without deployment.
				// 'EXPECT' means that the function will be called.
				// The next we call the function we need ex:'GetPaymentTypeByID'
				// 'Return' let us set the values which will be returned. We can also return an error.
				// With 'Times' we set how many times the function will be called.
				mock.RepoPaymentType.EXPECT().GetPaymentTypeById(2).
					Return(models.PaymentType{}, nil).Times(1)

				// Here we are mocking the time of our 'Clock' which is a wrapper of the system service 'Time'
				// With the value of 'currentTime'.
				mock.Clock.EXPECT().Now().Return(currentTime).Times(1)

				// Into 'DateTime' we put the 'currentTime'.
				// So now for the test we have the same time into the struct and into the mocked 'Clock'.
				accTransaction := &models.AccountTransaction{
					DateTime:    currentTime,
					PaymentType: models.PaymentType{},
					AccountFrom: models.Account{},
					AccountTo:   models.Account{},
					Order:       models.Order{},
					AmountCents: 50}

				//We call this func like in the order of calls into the real 'AddMoneyToAccount'.
				mock.RepoAccountTransaction.EXPECT().AddAccountTransaction(accTransaction).
					Return(nil).Times(1)

				// In this case we expect that function will be called without any errors.
				err := mock.AccountServiceUC.AddMoneyToAccount(accTransaction.AccountTo, 50)

				// Compare that expected value of error is nil.
				assert.Equal(t, nil, err)
			},
		}, { // In this case we are going by getting the error.
			name: "Incorrect.Got error from GetPaymentTypeByID",
			test: func(t *testing.T, mock *accountUseCasesMock) {

				// Describe which error we'll get.
				expectedError := errors.New("expectedError")

				// Call 'GetPaymentTypeById' and return here our 'expectedError'.
				mock.RepoPaymentType.EXPECT().GetPaymentTypeById(2).
					Return(models.PaymentType{}, expectedError).Times(1)

				accTransaction := &models.AccountTransaction{
					DateTime:    time.Now(),
					PaymentType: models.PaymentType{},
					AccountFrom: models.Account{},
					AccountTo:   models.Account{},
					Order:       models.Order{},
					AmountCents: 50}

				// Calling 'AddMoneyToAccount' will return us the error, because we had the error into the func before.
				err := mock.AccountServiceUC.AddMoneyToAccount(accTransaction.AccountTo, 50)

				assert.Error(t, err)
				assert.Equal(t, expectedError, err)
			},
		},
	})
}

func Test_Account_TakeMoneyFromAccount(t *testing.T) {
	var currentTime = time.Date(2021, 12, 19, 12, 21, 00, 00, time.UTC)
	accTransaction := &models.AccountTransaction{
		DateTime:    currentTime,
		PaymentType: models.PaymentType{},
		AccountFrom: models.Account{ID: 1},
		AccountTo:   models.Account{},
		Order:       models.Order{},
		AmountCents: 100,
	}
	accTransList := &models.AccountTransactionList{
		AccountTransactions: []models.AccountTransaction{
			{AccountTo: models.Account{ID: 1}, AmountCents: 100},
			{AccountTo: models.Account{ID: 1}, AmountCents: 50},
		},
	}

	runTestCases(t, []accountTestCase{
		{
			name: "Correct",
			test: func(t *testing.T, mock *accountUseCasesMock) {

				mock.RepoPaymentType.EXPECT().GetPaymentTypeById(3).
					Return(models.PaymentType{}, nil).Times(1)

				mock.Clock.EXPECT().Now().Return(currentTime).Times(1)

				mock.RepoAccountTransaction.EXPECT().
					GetAccountTransactionsInTimePeriod(time.UnixMilli(0), currentTime, accTransaction.AccountFrom).
					Return(accTransList, nil).Times(1)

				mock.RepoAccountTransaction.EXPECT().AddAccountTransaction(accTransaction).
					Return(nil).Times(1)

				err := mock.AccountServiceUC.TakeMoneyFromAccount(accTransaction.AccountFrom, 100)

				assert.Equal(t, nil, err)
			},
		},
		{
			name: "Incorrect, not enough money",
			test: func(t *testing.T, mock *accountUseCasesMock) {

				mock.RepoPaymentType.EXPECT().GetPaymentTypeById(3).
					Return(models.PaymentType{}, nil).Times(1)

				mock.Clock.EXPECT().Now().Return(currentTime).Times(1)

				mock.RepoAccountTransaction.EXPECT().
					GetAccountTransactionsInTimePeriod(time.UnixMilli(0), currentTime, accTransaction.AccountFrom).
					Return(accTransList, nil).Times(1)

				err := mock.AccountServiceUC.TakeMoneyFromAccount(accTransaction.AccountFrom, 200)

				assert.Error(t, err)
				assert.Equal(t, ErrNotEnoughMoneyToTake, err)
			},
		},
	})
}
