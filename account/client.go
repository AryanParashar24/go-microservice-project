// here we need all the functions we specified in the service.go file
package account

import (
	"context"

	"github.com/AryanParashar24/Go-Microservice-Project/account/pb" // import the generated protobuf package
	"google.golang.org/grpc"
)

// here we will be defining the same client as we defined in the Server struct of graph.go file int he main dir

// Here as we can see that the grpc and graphql servers both are present so gRPC server is helping in building up the overall server while the ghraphql server is hellping to make us use the clients of each of these 3 microservice
// As we remembered from the diagram, graphql makes request to the clients of each of these microservicess and CLI will make requests to serves of those microservices
type Client struct {
	conn    *grpc.ClientConn
	service pb.AccountServiceClient
}

func NewClient(url string) { // now since we have a lcient struct we just need a function to take url and initialize a connection to return us a Client
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := pb.NewAccountServiceClient(conn) // now we have a connection and we can use it to create a new client
	return &Client{conn, c}, nil
}

func (c *Client) Close() { // this function will be abel to close our conenction with the c CLient
	c.conn.Close()
}

// Now we'll make the fucntions to make call to all the three functions: PostAccount, GetAccount, GetAccounts been defined int he service.go
func (c *Client) PostAccount(ctx context.Context, name string) (*Account, error) {
	r, err := c.service.PostAccount( // callign to the struct PostAccount for which we are making this function and then linking it to that strcut
		ctx,
		&pb.PostAccountRequest{Name: name},
	)

	if err != nil {
		return nil, err
	}

	/* will make this function then will call the struct been defined in the service.go with the name of PostAccount with the context and name string
	Now since we took the params entereed as string and contextx thus we'll return the Account also the name and the aim of this function is to Post and Create Account
	Since its work is just to Post Account to the database or to the entries list
	*/
	return &Account{
		ID:   r.Account.Id,
		Name: r.Account.Name,
	}, nil
}

func (c *Client) GetAccount(ctx context.Context, id string) (*Account, error) {
	r, err := c.service.PostAccount(
		ctx,
		&pb.PostAccountRequest{Name: name},
	)

	if err != nil {
		return nil, err
	}

	return &Account{
		ID:   r.Account.Id,
		Name: r.Account.Name,
	}, nil
}

func (c *Client) GetAccounts(ctx context.Context, skip uint64, take uint64) (*[]Account, error) {
	r, err := c.service.GetAccounts(
		ctx,
		&pb.GetAccountsRequest{
			Skip: skip,
			Take: take,
		},
	)

	if err != nil {
		return nil, err
	}

	// here we will make an empty slice of accounts and then append the accounts to it
	accounts := []Account{}
	for _, a := range r.Accounts {
		accounts = append(accounts, Account{
			ID:   a.Id,
			Name: a.Name,
		})
	}
	return &accounts, nil
}
