package product

import (
	"context"

	"github.com/google/uuid"
)

type Service struct {
	repository *Repository
}

func NewService(r *Repository) *Service {
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
