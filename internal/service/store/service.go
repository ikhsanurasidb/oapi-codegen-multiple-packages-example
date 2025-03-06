package store

import (
	"context"
	"errors"

	gen_store "github.com/oapi-codegen-multiple-packages-example/internal/gen/store"
	"github.com/oapi-codegen-multiple-packages-example/internal/repository/store"
)

type Service interface {
    CreateOrder(ctx context.Context, order gen_store.Order) (*gen_store.Order, error)
    GetOrderByID(ctx context.Context, orderID int64) (*gen_store.Order, error)
    DeleteOrder(ctx context.Context, orderID int64) error
    GetInventory(ctx context.Context) (map[string]int32, error)
}

type storeService struct {
    repo store.Repository
}

func NewService(repo store.Repository) Service {
    return &storeService{
        repo: repo,
    }
}

func (s *storeService) CreateOrder(ctx context.Context, order gen_store.Order) (*gen_store.Order, error) {
    if order.PetId == nil {
        return nil, errors.New("pet ID is required")
    }
    
    if order.Status == nil {
        status := gen_store.Placed
        order.Status = &status
    }
    
    return s.repo.CreateOrder(ctx, order)
}

func (s *storeService) GetOrderByID(ctx context.Context, orderID int64) (*gen_store.Order, error) {
    order, err := s.repo.GetOrderByID(ctx, orderID)
    if err != nil {
        return nil, err
    }
    
    if order == nil {
        return nil, errors.New("order not found")
    }
    
    return order, nil
}

func (s *storeService) DeleteOrder(ctx context.Context, orderID int64) error {
    return s.repo.DeleteOrder(ctx, orderID)
}

func (s *storeService) GetInventory(ctx context.Context) (map[string]int32, error) {
    return s.repo.GetInventory(ctx)
}