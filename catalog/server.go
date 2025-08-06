package catalog

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/AryanParashar24/go-microservices-project/catalog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedCatalogServiceServer // this is the server that will implement the AccountService interface
	service                              Service
}

// gRPC server for listenning
func ListenGRPC(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	serv := grpc.NewServer() // here we are creating a new gRPC server
	pb.RegisterCatalogServiceServer(serv, &grpcServer{
		UnimplementedCatalogServiceServer: pb.UnimplementedCatalogServiceServer{}, // this is the server that will implement the AccountService interface
		service:                           s,                                      // here we are passing the service that we have created in the service
	})
	reflection.Register(serv) // this registers the reflection service on the gRPC server so that we can use it to introspect the service
	return serv.Serve(lis)    // this starts the gRPC server and listens for incoming connections
}

func (s *grpcServer) PostProduct(ctx context.Context, r *pb.PostProductRequest) (*pb.PostProductResponse, error) {
	p, err := s.service.PostProduct(ctx, r.Name, r.Description, r.Price)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.PostProductResponse{Product: &pb.Product{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}}, nil
}
func (s *grpcServer) GetProduct(ctx context.Context, r *pb.PostProductRequest) (*pb.PostProductResponse, error) {
	p, err := s.service.GetProduct(ctx, r.Id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.GetProductResponse{
		Product: &pb.Product{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		},
	}, nil
}
func (s *grpcServer) GetProducts(ctx context.Context, r *pb.PostProductRequest) (*pb.PostProductResponse, error) {
	var res []Product
	var err error
	if r.Query != "" {
		res, err = s.service.SearchProducts(ctx, r.Query, r.Skip, r.Take)
	} else if len(r.Ids) != 0 {
		res, err = s.service.GetProductByIDs(ctx, r.Ids)
	} else {
		res, err := s.service.GetProducts(ctx, r.Name, r.Skip, r.Take)
	}
	if err != nil {
		log.Println(err)
		return nil, err
	}

	products := []*pb.Product{} // colelction of products which will going to get retruned from this function after ranging over them as in the code below
	for _, p := range res {     // ranging over the response that has been hit from the GetProduct func or GetProductByIDs function which are all returning the collection of products
		//each of these products will be available to us via p
		products = append(products, &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		},
		)
	}
	return &pb.GetProductsResponse{Products: products}, nil
}
