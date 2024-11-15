package main

import (
	"database/sql"
	"github.com/igorzinar/goSocial/internal/db"
	"github.com/igorzinar/goSocial/internal/env"
	"github.com/igorzinar/goSocial/internal/store"
	"log"
)

const version = "0.0.1"

func main() {

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
	}

	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Panic(err)
		}
	}(db)
	log.Println("database connection pool established")
	storage := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  storage,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
