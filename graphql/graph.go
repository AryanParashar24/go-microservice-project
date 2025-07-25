package main

type Server struct { // with the help of server struct we'll call 3 diff clients as the account, catalog and order client which ha e further gRPC connections with the servers of all these clients
	// 	accountClient	*account.Client
	// 	catalogClient	*catalog.Client
	// 	orderClient		*order.Client
}

func NewGraphQLServer(accountUrl, catalogUrl, orderUrl string) (*Server, error) { // here this func will be accessing the account, catalog and order url int he forms of string while the return type will be of Server struct and error
	// 	accountClient, err := account.NewClient(accountUrl)
	// 	if err!= nil {
	// 		return nil, err
	// 	}

	// 	catalogClient, err := catalog.NewClient(catalogUrl)
	// 	if err != nil {
	// 		accountClient.Close() // if there is an error in the catalog client then we will close the account client and return nil and err
	// 		return nil, err
	// 	}

	// 	orderClient, err := order.NewClient(orderUrl)
	// 	if err != nil {
	// 		catalogClient.Close() // if there is an error in the order client then we will close the catalog client and return nil and err
	// 		accountClient.Close()
	// 		return nil, err
	// 	}

	return &Server{
		// 		accountClient: accountClient,
		// 		catalogClient: catalogClient,
		// 		orderClient:   orderClient,
	}, nil
}

// func (s *Server) Mutattion() MutationResolver { // this will return the mutation resolver which will have the server as a field
// 	return &mutationR  esolver{server: s}
// }

// func (s *Server) Query() QueryResolver { // this will return the query resolver which will have the server as a field
// 	return &queryResolver{server: s}
// }

// func (s *Server) Account() AccountResolver { // this will return the account resolver which will have the server as a field
// 	return &accountResolver{server: s}
// }

//  func (s *Server) ToExecutableSchema() graphql.ExecutableSchema { // this will return the executable schema which will have the server as a field
// 	return &executableSchema(Config
// 	{
// 		server: s,
// 	})
// }
