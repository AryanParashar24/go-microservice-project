package order

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error)
	GetOrdersForAccount(ctx context.Context, accoubtID string) ([]Order, error)
}

type Order struct {
	ID         string           `json:"id"`
	CreatedAt  time.Time        `json:"created_at"`
	TotalPrice float64          `json:"total_price"`
	AccountID  string           `json:"account_id"`
	Products   []OrderedProduct `json:"products"`
}

type OrderedProduct struct {
	ID         string
	Name       string
	Desription string
	Price      float64
	Quantity   uint32
}

type orderService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &orderService{r}
}

func (s orderService) PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error) {
	o := &Order{
		ID:        ksuid.New().String(), // generating a new unique id for the order
		CreatedAt: time.Now().UTC(),
		AccountID: accountID,
		Products:  products,
	}

	o.TotalPrice = 0.0
	for _, p := range products {
		o.TotalPrice += p.Price * float64(p.Quantity) // calculating the total price of the order by multiplying the price of each product with its quantity
	}
	err := s.repository.PutOrder(ctx, *o) // here we are putting the order in the repository
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (s orderService) GetOrdersForAccount(ctx context.Context, accoutID string) ([]Order, error) {
	return s.repository.GetOrdersForAccount(ctx, accoutID) // here we are getting the orders for the account from the repository
}
