package main

import (
	"fmt"
	"log"

	"github.com/SaidovZohid/hotel-project/api"
	"github.com/SaidovZohid/hotel-project/config"
	"github.com/SaidovZohid/hotel-project/storage"
	"github.com/go-redis/redis/v9"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/SaidovZohid/hotel-project/api/docs"
)

func main() {
	cfg := config.Load(".")
	log.Println(cfg.Auth)
	psqlUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
	)

	psqlConn, err := sqlx.Connect("postgres", psqlUrl)

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	rda := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
	strg := storage.NewStorage(psqlConn)
	inMemory := storage.NewInMemoryStorage(rda)
	apiServer := api.New(&api.RouteOptions{
		Cfg: &cfg,
		Stgr: strg,
		InMemory: inMemory,
	})

	err = apiServer.Run(cfg.HttpPort)
	if err != nil {
		log.Fatalln(err)
	}
}
