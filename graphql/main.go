package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	AccountURL string `envconfig: "ACCOUNT_SERVICE_URL"`
	CatalogURL string `envconfig: "CATALOG_SERVICE_URL"`
	OrderURL   string `envconfig: "ORDER_SERVICE_URL"`
}

func main() {
	var cfg AppConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	s, err := NewGraphQLServer(cfg.AccountURL, cfg.CatalogURL, cfg.OrderURL) // here as we can go back and see these are the fields that are being accepted by the NewGraphQL Server i.e. Account,Catalog,Order URLs
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/graphql", handler.New(s.ToExecutableSchema()))
	// In graphql we gets a UI named playground where we can run our queries to be able to interact with the graphQL server and can handle our server along with the database
	http.Handle("/playground", playground.Handler("Aryan", "/graphql")) // here we are using the handle package to handle the playground and the graphql server

	log.Fatal(http.ListenAndServe(":8080", nil)) // here we are starting the server on port 8080 and passing nil as the handler
}
