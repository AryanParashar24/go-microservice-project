package main

import (
	"log"
	"time"

	"github.com/AryanParashar24/go-microservices-project/catalog"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envcongif: "DATABASE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg) // Load configuration from environment variables
	if err != nil {
		log.Fatal(err)
	}
	var r catalog.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err := catalog.NewElasticRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer r.Close()

	log.Println("Listenning on port 8080...")
	s := catalog.NewService(r)
	log.Fatal(catalog.ListenGRPC(s, 8080))
}
