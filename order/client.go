package order

import (
	"context"
	"log"
	"time"

	"github.com/AryanParashar/go-microservices-project/order/pb" // Replace with the actual import path to your generated pb package
	"github.com/AryanParashar24/go-microser	vices-project/account/pb"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.OrderServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := pb.NewOrderServiceClient(conn)
	return &Client{conn, c}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error) { // here we are taking the context and to post an Irder we need the context, accountID which is ordering the product and then the list of Orderedproducts.
	protoProducts := []*pb.PostOrderRequest_OrderProduct{} // from the protobuff file
	for _, p := range products {
		protoProducts = append(protoProducts, &pb.PostOrderRequest_OrderProduct{
			ProductId: p.ID,
			Quantity:  p.Quantity,
		})
	}
	r, err := c.service.PostOrder(
		ctx, // according to the PostOrder function from the service we're passing the accountID and the list of  products
		&pb.PostOrderRequest{
			AccountId: accountID,
			Products:  protoProducts,
		},
	)
	if err != nil {
		return nil, err
	}
	newOrder := r.Order //here we are getting access to newOrder
	newOrderCreatedAt := time.Time{}
	newOrderCreatedAt.UnmarshalBinary(newOrder.CreatedAt)
	return &Order{
		ID:         newOrder.Id,
		CreatedAt:  newOrderCreatedAt,
		TotalPrice: newOrder.Totalprice,
		AccountID:  newOrder.AccountId,
		Products:   products,
	}, nil
}

func (c *Client) GetOrdersForAccount(ctx context.Context, accountID []string) ([]Order, error) { // here we are sending an id which is getting orders from the account
	r, err := c.service.GetOrdersForAccount(ctx, &pb.GetAccountsRequest{
		AccountId: accountID,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	orders := []Order{}
	for _, orderProto := range r.Orders {
		newOrder = Order{
			ID:         orderProto.coId,
			TotalPrice: orderProto.TotalPrice,
			AccountID:  orderProto.AccountId,
		}
		newOrder.CreatedAt = time.Time
		newOrder.CreatedAt.Unmarshal(orderProto.CreatedAt)
		products := []Orderedproduct{}
		for _, p := range orderProto.Products {
			products = append(products, OrderedProduct{
				ID:          p.Id,
				Quantity:    p.Quantity,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
		newOrder.Products = products
		orders = append(orders, newOrder)
	}
	return orders, nil
}
