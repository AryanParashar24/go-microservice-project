// Here in the qraphql we'll have 2 queries which will be GetAllAccounts and GetAllProducts thus we have listed Accounts and Products only
// Accounts // Now we'll have to have different account and product ids which will be listed in a seperate file as in the account_resolver.go and

// Products

package main

// type queryResolver struct {
// 	server *Server // this will have the server with the name Server as a field which will be used to call the account, catalog and order clients
// }

// as we have seen in the graphQL diagram where we drew 2 blocks as the queries and another one as the mutations with mutations been resolved in the graph, account, mutation, queries resolvers
//  now we are letf with defining about the queries which are for the Accounts, Products and the Orders

// Here in the queries_resolver as we can see above we'll only define the Accounts and the Products queries while the Orders as was defined earlier in the account_resolver we will be defining to handle Orders in the account_resolver.go
// func (r *queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id *string)([]*Account, err){

// }

// func (r *queryResolver) Products(ctx context.Context, pagination *PaginationInput, query *string, id *string)([]*Product, error){

// }
