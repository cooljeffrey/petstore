package service

import (
	"context"
	"github.com/cooljeffrey/petstore/model"
	"github.com/go-kit/kit/log"
)

type StoreService interface {
	GetInventoriesByStatus(ctx context.Context) (map[string]int64, error)
	PlaceOrder(ctx context.Context, order *model.Order) (*model.Order, error)
	FindOrderByID(ctx context.Context, id int64) (*model.Order, error)
	DeleteOrderByID(ctx context.Context, id int64) error
}

type storeService struct {
	logger  log.Logger
	storage model.Storage
}

func NewStoreService(logger log.Logger, storage model.Storage) StoreService {
	return &storeService{
		logger:  logger,
		storage: storage,
	}
}

func (s storeService) GetInventoriesByStatus(ctx context.Context) (map[string]int64, error) {
	return s.storage.RetrieveStoreInventoriesByStatus()
}

func (s storeService) PlaceOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	return s.storage.CreateOrder(order)
}

func (s storeService) FindOrderByID(ctx context.Context, id int64) (*model.Order, error) {
	return s.storage.RetrieveOrderByID(id)
}

func (s storeService) DeleteOrderByID(ctx context.Context, id int64) error {
	return s.storage.DeleteOrderByID(id)
}
