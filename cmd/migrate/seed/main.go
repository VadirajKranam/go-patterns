package main

import (
	"log"

	"github.com/vadiraj/gopher/internal/db"
	"github.com/vadiraj/gopher/internal/env"
	"github.com/vadiraj/gopher/internal/store"
)

func main(){
	addr:=env.GetString("DB_ADDR", "postgres://admin:admin@localhost/socialnetwork?sslmode=disable")
	log.Print(addr)
	conn,err:=db.New(addr,3,3,"15m")
	if err!=nil{
		log.Fatal(err)
	}
	defer conn.Close()
	store:=store.NewStorage(conn)
	db.Seed(store,conn)
}