// For our mutations we shoudl be abel to create an Account, product and an order
// once we have theese queeries enxt step will be of linking these queries to the microservices
// CreateAccount
// CreateProduct
// CreateOrder

package main

// import "context"

// type mutationResolver struct {
// 	server *Server // this will have the server as a field which will be used to call the account, catalog and order clients
// }

// Mutation server is for cerating multiple mutations like Create Account , create Order and create Product as we have shown inthe Schema  diagram

// func (r *mutationResolver) cerateAccount(ctx context.Context, in AccountInput) (*Account, error) {
// 	account, err := r.server.accountClient.CreateAccount(ctx, &input)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return account, nil
// }

// func (r *mutationResolver) createProduct(ctx context.Context, in ProductInput) (*Product, error) {
// 	product, err := r.server.catalogClient.CreateProduct(ctx, &input)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return product, nil
// }

// func (r *mutationResolver) createOrder(ctx context.Context, in OrderInput) (*Order, error) {
// 	order, err := r.server.orderClient.CreateOrder(ctx, &input)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return order, nil
// }
