package catalog

import (
	"context"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PutProduct(ctx context.Context, name, description string, price float64) (*Prouct, error)
	GetProduct(ctx context.Context, id string) (*Product, error)
	GetProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	GetProductByIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
}

type Products struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       stiring `json:"price"`
}

type catalogService struct {
	repository repository
}

func NewService(r Repository) Service {
	return &catalogService{r}
}

func (s *catalogService) PutProduct(ctx context.Context, name, description string, price float64) (*Prouct, error) {
	p := &productDocument{
		Name:        name,
		Description: description,
		Price:       price,
		ID:          ksuid.New().String(), // generate a new ID for the product
	}
	if err := s.repository.putProduct(ctx, *p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *catalogService) GetProduct(ctx context.Context, id string) (*Product, error) {
	return s.repository.GetProductsByID(ctx, id)
}

func (s *catalogService) GetProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repository.ListProducts(ctx, skip, take)
}

func (s *catalogService) GetProductByIDs(ctx context.Context, ids []string) ([]Product, error) {
	return s.repository.ListProductWithIDs(ctx, ids)
}

func (s *catalogService) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repository.SearchProducts(ctx, query, skip, take)
}
