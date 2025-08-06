package account

import (
	"context"
	"fmt"
	"net"

	"github.com/AryanParashar24/go-microservices-project/account/pb" // import the generated protobuf package
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedAccountServiceServer // this is the server that will implement the AccountService interface
	service                              Service
}

// gRPC server for listenning
func ListenGRPC(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	serv := grpc.NewServer() // here we are creating a new gRPC server
	pb.RegisterAccountServiceServer(serv, &grpcServer{
		UnimplementedAccountServiceServer: pb.UnimplementedAccountServiceServer{}, // this is the server that will implement the AccountService interface
		service:                           s,                                      // here we are passing the service that we have created in the service
	})
	reflection.Register(serv) // this registers the reflection service on the gRPC server so that we can use it to introspect the service
	return serv.Serv(lis)     // this starts the gRPC server and listens for incoming connections
}

func (s *grpcServer) PostAccount(ctx context.Context, r *pb.PostAccountRequest) (*pb.PostAccountResponse, error) {
	a, err := s.service.PostAccount(ctx, r.Name)
	if err != nil {
		return nil, err
	}
	return &pb.PostAccountResponse{Account: &pb.Account{ // her in the PostAccountResponse and in GetAccountResponse, we have only one Account but in GetAccounts have multiple
		Id:   a.ID,
		Name: a.Name,
	}}, nil
}

func (s *grpcServer) GetAccount(ctx context.Context, r *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	a, err := s.service.GetAccount(ctx, r.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetAccountResponse{
		Account: &pb.Account{
			Id:   a.ID,
			Name: a.Name,
		},
	}, nil
}

func (s *grpcServer) GetAccounts(ctx context.Context, r *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
	res, err := s.service.GetAccounts(ctx, r.Skip, r.Take)
	if err != nil {
		return nil, err
	}
	// now here mutiple Accounts are present int he GetAccounts with take=1 and skip=1
	//  so to handle it we will make an empty collection of accounts
	accounts := []*pb.Account{} // here we opened an empty collection and then we will
	for _, p := range res {     // Now we are going to response over the accounts and will append the accounts to the empty collection
		accounts = append(
			accounts, // now will append the response over the range
			&pb.Account{
				Id:   p.ID,
				Name: p.Name,
			},
		)
	}
	return &pb.GetAccountsResponse{Accounts: accounts}, nil
}
