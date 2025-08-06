package order

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/AryanParashar24/go-microservices-project/account"
	"github.com/AryanParashar24/go-microservices-project/catalog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedOrderServiceServer // this is the server that will implement the AccountService interface
	service       Service
	accountClient Account.Client
	catalogClient Catalog.Client
}

func ListenGRPC(s service, accountURL, catalogURL string, port int) error {
	accountClient, err := account.NewClient(accountURL)
	if err != nil {
		return err
	}
	catalogClient, err := catalog.NewClient(catalogURL)
	if err != nil {
		return err
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		accountClient.Close()
		return err
	}
	serv := grpc.NewServer()
	pb.RegisterOrderServiceServer(serv, &grpcServer{
		UnimplementedOrderServiceServer: pb.UnimplementedOrderServiceServer{}, // this is the server that will implement the AccountService interface
		service: s,
		accountClient: accountClient,
		catalogClient: catalogClient,
	})
	reflection.Register(serv)
	return serv.Serv(lis)
}

func (s *grpcServer) PostOrder(ctx context.Context, r *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	_, err := s.accountClient.GetAccount(ctx, r.AccountId)
	if err != nil {
		log.Println("Error getting account:", err)
		return nil, errors.New("access not found")
	}

	// get orderd products
	productIDs := []string{}
	/// here we will be managing with the CLients of all the 3 diffeerent microservices in our applciation architecture
	for _, p := range r.Products {
		productIDs = append(productIDs, p.ProductId) // here we are appending the productIDs to the
	}
	orderedProdcts, err := s.catalog.Client.GetProducts(ctx, 0, 0, productIDs, "") // here we are pulling the client of our catalog microservice to get the products
	if err != nil {
		log.Println("Error getting products: ", err)
		return nil, errors.New("products not found")
	}
	// construct products
	products := []OrderedProducts{} // here the products slice from our orderedProduct is been pulled with the following methods been defined
	for _, p := range orderedProducts {
		product := OrderedProduct{
			ID:          p.ID,
			Quantity:    0,
			Price:       p.Price,
			Name:        p.Name,
			Description: p.Description,
		}
		for _, rp := range r.Products {
			if rp.ProductId == p.ID {
				product.Quantity = rp.Quantity
				break
			}
		}
		if product.Quantity != 0 {
			products = append(products, product)
		}
	}
	// call service implementation
	order, err := s.service.PostOrder(ctx, r.AccountId, products)
	if err != nil {
		log.Println("Error posting order:", err)
		return nil, errors.New("could not post order")
	}

	orderProto := &pb.Order{
		Id:         order.ID,
		AccountId:  order.AccountID,
		TotalPrice: order.TotalPrice,
		Products:   []*pb.Order_Orderproduct{},
	}

	orderProto.CreatedAt, _ = order.CreatedAt.MarshalBinary()
	for _, p := range order.Products {
		orderProto.Products = append(orderProto.Products, &pb.Order_OrderProduct{ // herr we will be appending products to the order list defined in protoc file
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			price:       p.Price,
			Quantity:    p.Quantity,
		})
	}
	return &pb.PostOrderResponse{
		Order: orderProto,
	}, nil
}

func (s *grpcServer) GetOrdersFromAccount(ctx context.Context, r *pb.GetOrderForAccountRequest) (*pb.GetOrderForAccountResponse, error) {
	accountOrders, err := s.service.GetOrdersForAccount(ctx, r.AccountId)	// if we want to get orders for a particular account then we can pass on the account id and then that will call GetOrdersForAccount 
	if err != nil {
		log.Println(err)
		return nil, err
	}

	productIDMap := map[string]bool{}
	for _, o := range accountOrders {
		for _, p := range o.Products {	// we also wants to extract all those products which are present in the orders of the account becasue we r ranging over accountOrders
			productIDMap[p.ID] = true	// each order will be having the list of orderedProducts along their orderiDs arranged in the productMap as we described in the line above
		}
	}
	productIDs := []string{}
	for id := range productIDMap {
		productIDs = append(productIDs, id)
	}
	products, err := s.catalogClient.GetProducts(ctx, 0, 0, productIDs, "")
	if err!= nil{
		log.println("Error getting account products: ", err)
		return nil, err
	}
	orders:= []*pb.Order{}	// here we r creating a slice of orders from the protobuff file we created for our Orders, where we defined for OrderedProduct and various order methods like id, name, description, quantity, price, etc.
	for _, o:= &pb.Order{
		op := &pb.Order{
		    AccountId: o.AccountID,
			Id: o.ID,
			TotalPrice: o.TotalPrice,
			products: []*pb.Order_OrderProduct{},		
		}
		op.CreatedAt, _ = o.CreatedAt.MarshalBinary()

		for _, product:= range o.Products {
			for _, p:= range products{
				if p.ID == product.ID{
					product.Name = p.Name
					product.Description = p.Description
					product.price = p.Price
					break
				}
			}
			op.Products = append(op.Products, &pb.Order_OrderProduct{	// here op which we defined with the pb for the Orders from the protobuff file will now handle Products string and will append its items to pb file for Order_OrderProducts
				Id: productsID,
				name: products.Name,
				Description: products.Description,
				Price: products.Price,
				Quantity: product.Quantity,
			})
		}
		orderes = append(orders, op)
	}
	return &pb.GetOrdersForAccountResponse{Orders: orders}, nil
}
