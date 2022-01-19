package services

import (
	"Dp218GO/models"
	"Dp218GO/repositories/mock"
	"errors"
	"github.com/golang/mock/gomock"
	assert "github.com/stretchr/testify/require"
	"testing"
)

type OrderMock struct {
	OrderService *OrderService
	RepoOrder    *mock.MockOrderRepo
}

type orderTestCase struct {
	name string
	test func(t *testing.T, mock *OrderMock)
}

func runOrderTestCases(t *testing.T, testCases []orderTestCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					tt.Error(err)
				}
			}()

			ctrl := gomock.NewController(tt)
			defer ctrl.Finish()

			mock := NewOrderMock(ctrl)

			tc.test(tt, mock)
		})
	}
}

func NewOrderMock(ctrl *gomock.Controller) *OrderMock {
	repoOrder := mock.NewMockOrderRepo(ctrl)

	orderService := NewOrderService(repoOrder)

	return &OrderMock{
		OrderService: orderService,
		RepoOrder:    repoOrder,
	}
}

func TestOrderService_CreateOrder(t *testing.T) {
	runOrderTestCases(t, []orderTestCase{
		{
			name: "Correct",
			test: func(t *testing.T, mock *OrderMock) {
				mock.RepoOrder.EXPECT().CreateOrder(models.User{}, 1, 2, 3, 1000.50).Return(models.Order{},
					nil).Times(1)

				_, err := mock.OrderService.CreateOrder(models.User{}, 1, 2, 3, 1000.50)
				assert.Equal(t, nil, err)
			},
		}, {
			name: "Incorrect",
			test: func(t *testing.T, mock *OrderMock) {
				expectedError := errors.New("expectedError")
				mock.RepoOrder.EXPECT().CreateOrder(models.User{}, 1, 2, 3, 1000.50).Return(models.Order{},
					expectedError).Times(1)

				_, err := mock.RepoOrder.CreateOrder(models.User{}, 1, 2, 3, 1000.50)
				assert.Error(t, err)
				assert.Equal(t, expectedError, err)
			},
		},
	})
}

func TestOrderService_GetAllOrders(t *testing.T) {
	runOrderTestCases(t, []orderTestCase{
		{
			name: "Correct",
			test: func(t *testing.T, mock *OrderMock) {
				mock.RepoOrder.EXPECT().GetAllOrders().Return(&models.OrderList{},
					nil).Times(1)

				_, err := mock.OrderService.GetAllOrders()
				assert.Equal(t, nil, err)
			},
		}, {
			name: "Incorrect",
			test: func(t *testing.T, mock *OrderMock) {
				expectedError := errors.New("expectedError")
				mock.RepoOrder.EXPECT().GetAllOrders().Return(&models.OrderList{},
					expectedError).Times(1)

				_, err := mock.RepoOrder.GetAllOrders()
				assert.Error(t, err)
				assert.Equal(t, expectedError, err)
			},
		},
	})
}

func TestOrderService_GetOrderByID(t *testing.T) {
	runOrderTestCases(t, []orderTestCase{
		{
			name: "Correct",
			test: func(t *testing.T, mock *OrderMock) {
				mock.RepoOrder.EXPECT().GetOrderByID(1).Return(models.Order{},
					nil).Times(1)

				_, err := mock.OrderService.GetOrderByID(1)
				assert.Equal(t, nil, err)
			},
		}, {
			name: "Incorrect",
			test: func(t *testing.T, mock *OrderMock) {
				expectedError := errors.New("expectedError")
				mock.RepoOrder.EXPECT().GetOrderByID(1).Return(models.Order{},
					expectedError).Times(1)

				_, err := mock.RepoOrder.GetOrderByID(1)
				assert.Error(t, err)
				assert.Equal(t, expectedError, err)
			},
		},
	})
}

func TestOrderService_GetOrdersByUserID(t *testing.T) {
	runOrderTestCases(t, []orderTestCase{
		{
			name: "Correct",
			test: func(t *testing.T, mock *OrderMock) {
				mock.RepoOrder.EXPECT().GetOrdersByUserID(1).Return(models.OrderList{},
					nil).Times(1)

				_, err := mock.OrderService.GetOrdersByUserID(1)
				assert.Equal(t, nil, err)
			},
		}, {
			name: "Incorrect",
			test: func(t *testing.T, mock *OrderMock) {
				expectedError := errors.New("expectedError")
				mock.RepoOrder.EXPECT().GetOrdersByUserID(1).Return(models.OrderList{},
					expectedError).Times(1)

				_, err := mock.RepoOrder.GetOrdersByUserID(1)
				assert.Error(t, err)
				assert.Equal(t, expectedError, err)
			},
		},
	})
}

func TestOrderService_GetOrdersByScooterID(t *testing.T) {
	runOrderTestCases(t, []orderTestCase{
		{
			name: "Correct",
			test: func(t *testing.T, mock *OrderMock) {
				mock.RepoOrder.EXPECT().GetOrdersByScooterID(1).Return(models.OrderList{},
					nil).Times(1)

				_, err := mock.OrderService.GetOrdersByScooterID(1)
				assert.Equal(t, nil, err)
			},
		}, {
			name: "Incorrect",
			test: func(t *testing.T, mock *OrderMock) {
				expectedError := errors.New("expectedError")
				mock.RepoOrder.EXPECT().GetOrdersByScooterID(1).Return(models.OrderList{},
					expectedError).Times(1)

				_, err := mock.RepoOrder.GetOrdersByScooterID(1)
				assert.Error(t, err)
				assert.Equal(t, expectedError, err)
			},
		},
	})
}

func TestOrderService_GetScooterMileageByID(t *testing.T) {
	runOrderTestCases(t, []orderTestCase{
		{
			name: "Correct",
			test: func(t *testing.T, mock *OrderMock) {
				mock.RepoOrder.EXPECT().GetScooterMileageByID(1).Return(100.500,
					nil).Times(1)

				_, err := mock.OrderService.GetScooterMileageByID(1)
				assert.Equal(t, nil, err)
			},
		}, {
			name: "Incorrect",
			test: func(t *testing.T, mock *OrderMock) {
				expectedError := errors.New("expectedError")
				mock.RepoOrder.EXPECT().GetScooterMileageByID(1).Return(100.500,
					expectedError).Times(1)

				_, err := mock.RepoOrder.GetScooterMileageByID(1)
				assert.Error(t, err)
				assert.Equal(t, expectedError, err)
			},
		},
	})
}

func TestOrderService_GetUserMileageByID(t *testing.T) {
	runOrderTestCases(t, []orderTestCase{
		{
			name: "Correct",
			test: func(t *testing.T, mock *OrderMock) {
				mock.RepoOrder.EXPECT().GetUserMileageByID(1).Return(100.500,
					nil).Times(1)

				_, err := mock.OrderService.GetUserMileageByID(1)
				assert.Equal(t, nil, err)
			},
		}, {
			name: "Incorrect",
			test: func(t *testing.T, mock *OrderMock) {
				expectedError := errors.New("expectedError")
				mock.RepoOrder.EXPECT().GetUserMileageByID(1).Return(100.500,
					expectedError).Times(1)

				_, err := mock.RepoOrder.GetUserMileageByID(1)
				assert.Error(t, err)
				assert.Equal(t, expectedError, err)
			},
		},
	})
}

func TestOrderService_UpdateOrder(t *testing.T) {
	runOrderTestCases(t, []orderTestCase{
		{
			name: "Correct",
			test: func(t *testing.T, mock *OrderMock) {
				mock.RepoOrder.EXPECT().UpdateOrder(1, models.Order{}).Return(models.Order{},
					nil).Times(1)

				_, err := mock.OrderService.UpdateOrder(1, models.Order{})
				assert.Equal(t, nil, err)
			},
		}, {
			name: "Incorrect",
			test: func(t *testing.T, mock *OrderMock) {
				expectedError := errors.New("expectedError")
				mock.RepoOrder.EXPECT().UpdateOrder(1, models.Order{}).Return(models.Order{},
					expectedError).Times(1)

				_, err := mock.RepoOrder.UpdateOrder(1, models.Order{})
				assert.Error(t, err)
				assert.Equal(t, expectedError, err)
			},
		},
	})
}

func TestOrderService_DeleteOrder(t *testing.T) {
	runOrderTestCases(t, []orderTestCase{
		{
			name: "Correct",
			test: func(t *testing.T, mock *OrderMock) {
				mock.RepoOrder.EXPECT().DeleteOrder(1).
					Return(nil).Times(1)

				err := mock.OrderService.DeleteOrder(1)
				assert.Equal(t, nil, err)
			},
		}, {
			name: "Incorrect",
			test: func(t *testing.T, mock *OrderMock) {
				expectedError := errors.New("expectedError")
				mock.RepoOrder.EXPECT().DeleteOrder(1).Return(
					expectedError).Times(1)

				err := mock.RepoOrder.DeleteOrder(1)
				assert.Error(t, err)
				assert.Equal(t, expectedError, err)
			},
		},
	})
}
