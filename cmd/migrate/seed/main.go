package main

import (
	"fmt"
	"github.com/igorzinar/goSocial/internal/db"
	"github.com/igorzinar/goSocial/internal/env"
	"github.com/igorzinar/goSocial/internal/store"
	"log"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable")
	fmt.Println("value of addr is: ", addr)
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	seedStore := store.NewStorage(conn)

	db.Seed(seedStore)
}
