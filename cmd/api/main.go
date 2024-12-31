package main

import (
	"log"

	"github.com/vadiraj/gopher/internal/db"
	"github.com/vadiraj/gopher/internal/env"
	"github.com/vadiraj/gopher/internal/store"
)

func main(){
	cfg:=config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr: env.GetString("DB_ADDR", "postgres://admin:admin@localhost/socialnetwork?sslmode=disable"),
			maxOpenConnns:  env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS",30),
			maxIdleTime: env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}
	db,err:=db.New(cfg.db.addr, cfg.db.maxOpenConnns,cfg.db.maxIdleConns,cfg.db.maxIdleTime)
	defer db.Close()
	log.Printf("Database connection pool established")
	if err!=nil{
		log.Panic(err)
	}
	store:=store.NewStorage(db)
	app:=&application{
		config: cfg,
		store: store,
	}
	mux:=app.mount()
	log.Fatal(app.run(mux))
}
 