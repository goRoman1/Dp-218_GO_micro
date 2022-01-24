package service

import (
	"OrderService/proto"
	"OrderService/repository"
	"context"
)

type OrderInterface interface {
	CreateOrder(ctx context.Context, info *proto.TripInfo) (*proto.Order, error)
}

type OrderService struct {
	Repo *repository.OrderRepo
	*proto.UnimplementedOrderServiceServer
}

func (os *OrderService) mustEmbedUnimplementedOrderServiceServer() {
	//TODO implement me
	panic("implement me")
}

func NewOrderService(repo *repository.OrderRepo) *OrderService {
	return &OrderService{Repo: repo}
}

func (os *OrderService) CreateOrder(ctx context.Context, info *proto.TripInfo) (*proto.Order, error) {
	return os.Repo.CreateOrder(ctx, info)
}
