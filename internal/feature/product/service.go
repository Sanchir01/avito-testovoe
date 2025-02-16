package product

import (
	"context"

	"github.com/google/uuid"
)

type Service struct {
	repository ServiceProduct
}

type ServiceProduct interface {
	GetProductByID(ctx context.Context, id uuid.UUID) (*DataBaseProduct, error)
	GetAllProducts(ctx context.Context) ([]*DataBaseProduct, error)
}

func NewService(r ServiceProduct) *Service {
	return &Service{
		repository: r,
	}
}

func (s *Service) GetProductByID(ctx context.Context, id uuid.UUID) (*DataBaseProduct, error) {
	product, err := s.repository.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *Service) GetAllProducts(ctx context.Context) ([]*DataBaseProduct, error) {
	products, err := s.repository.GetAllProducts(ctx)
	if err != nil {
		return nil, err
	}
	return products, nil
}
