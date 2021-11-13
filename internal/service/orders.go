package service

import (
	"context"
	"github.com/cookienyancloud/back/internal/repository"
	"github.com/google/uuid"
)

type OrdersService struct {
	repo repository.OrdersRepo
}

func NewOrdersService(repo repository.OrdersRepo) *OrdersService {
	return &OrdersService{
		repo: repo,
	}
}

func (s *OrdersService) Create(ctx context.Context,orderId uuid.UUID, userId string, eventId int, zonesId []int) error {

	return nil
}
