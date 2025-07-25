package account
import (
	"context"
	"fmt"
	"net"

	"github.com/AryanParashar24/Go-Microservice-Project/account/pb" // import the generated protobuf package
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

)
type grpcServer struct {
	service Service
}

// gRPC server for listenning 
func Listenning(s Service, port int) error{
	lis, err:= net.listen("tcp", fmt.Sprintf(":%d", port))
	if err!= nil{
		return err
	}

	serv:= grpc.NewServer() // here we are creating a new gRPC server
	pb.RegisterAccountServiceServer(serv, &grpcServer(s))
	reflection.Register(serv) // this registers the reflection service on the gRPC server so that we can use it to introspect the service
	return serv.Serve(lis) // this starts the gRPC server and listens for incoming connections
}

func (s *grpcServer) PostAccount(ctx context.Context, r *pb.PostAccountRequest) (*pb.PostAccountResponse, error){
	a, err:= s.service.PostAccount(ctx, r.Name)
	if err!= nil{
		return nil, err
	}
	return &pb.PostAccountResponse{Account: &pb.Account{	// her in the PostAccountResponse and in GetAccountResponse, we have only one Account but in GetAccounts have multiple
		Id: a.ID,
		Name: a.Name,
	}}, nil
}

func (s *grpcServer)GetAccount(ctx context.Context, r.*pb.GetAccountRequest)(*pb.GetAccountResponse, err){
	a, err:= s.service.PostAccount(ctx, r.Id)
	if err!= nil{
		return nil, err
	}
	return &pb.GetAccountResponse{Account: &pb.Account{
		Id: a.ID,
		Name: a.Name,
	}}, nil
}

func (s *grpcServer)GetAccounts(ctx context,.Context, r.*pb.GetAccountsRequest)(*pb.GetAccountsResponse, error){
	res, err:= s.service.GetAccounts(ctx, r.Id)
	if err!= nil{
		return nil, err
	}
	// now here mutiple Accounts are present int he GetAccounts with take=1 and skip=1
	//  so to handle it we will make an empty collection of accounts
	accounts := []*pb.Account{}	// here we opened an empty collection and then we will
	for _, p:= range res{ // Now we are going to response over the accounts and will append the accounts to the empty collection
		accounts = append(accounts, // now will append the response over the range 
			&pb.Account{
				Id: p.ID,
				Nmae: p.Name,
			},
		)
	}
	return &pb.GetAccountsResponse{Accoutns: accounts}, nil	
}