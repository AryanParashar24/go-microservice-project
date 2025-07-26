package main

import (
	"log"
	"time"

	"github.com/AryanParashar24/go-microservices-project/account"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `enconfig: "DATABASE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var r account.Repository
	retry.ForeverSleep(2*time.Second, func(_int) (err error) {
		account.NewPostgressRepository(cfg.DatabaseURL) // here we are going to create a new Postgress account Repo which will be the connection to our database and URL
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer r.Close()
	log.Println("Listenning o port 8080...")
	s := account.NewService(r)             // here with s we can start a new microservice which will be connected to our URL and Database in the above code
	log.Fatal(account.ListenGRPC(s, 8080)) // here we are going to listen to the gRPC server and will start a  new service on the port 8080 from the function that have been defined above
}
