package main

import (
	"log"
	"time"

	"github.com/AryanParashar24/go-microservices-project/order"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig: "DATABASE_URL"`
	AccountURL  string `envconfig: "ACCOUNT_SERVICE_URL"`
	CatalogURL  string `envconfig: "CATALOG_SERVICE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var r order.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = order.NewPostgressRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer r.Close()
	log.Println("listenning on port 8080 ...")
	s := order.NewService(r)
	log.fatal(order.ListenGRPC(s, cfg.AccountURL, cfg.CatalogURL, 8080))
}
