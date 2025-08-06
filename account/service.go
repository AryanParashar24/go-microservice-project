package account

import (
	"context"

	"github.com/segmentio/ksuid"
)

type Service interface { // these interfaces call these functions differently by being a part of interfaces in the services
	PostAccount(ctx context.Context, name string) (*Account, error)
	GetAccount(ctx context.Context, id string) (*Account, error)
	GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type accountService struct {
	repository Repository
}

func NewService(r Repository) Service { // Now we have struct and now we have to use a function to call them out
	return &accountService{r} // here it is a method being called here in the function
}

// now we have 3 interfaces and will call them here in diff functions using the function defined and being called in the above newService func
func (s *accountService) PostAccount(ctx context.Context, name string) (*Account, error) { // adding the struct accountService to the function PostAccount & we'll going to call the function PutAccount from the repository.go as thats the function which can create account
	a := &Account{ // here above in the parameters for the function PostAccount we took context and just the name because while creating an account if we move back to our Graphql for the entry of AccountInput we can see just a name string in the values, methods defined
		Name: name,
		ID:   ksuid.New().String(), // here we are generating a new id for the account using the uuid package
	}
	if err := s.repository.PutAccount(ctx, *a); err != nil { // here we are calling the PutAccount function from the repository interface and passing the context and account as parameters
		return nil, err // if there is an error then we return nil and the error
	}
	return a, nil
}
func (s *accountService) GetAccount(ctx context.Context, id string) (*Account, error) {

	return s.repository.GetAccountByID(ctx, id)
}
func (s *accountService) GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	if take > 100 || (skip == 0 && take == 0) { // here we are checking if the take is greater than 100 or if both skip and take are 0 then we return an error
		take = 100
	}
	return s.repository.ListAccounts(ctx, skip, take)
}
