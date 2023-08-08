package main

import (
	"database/sql"
	"dmmak/simple-rest-crud/internal/handler"
	"dmmak/simple-rest-crud/internal/repo"
	"flag"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	var connStr string
	flag.StringVar(&connStr, "pgConn", "", "PostgresSQL connection string")
	db, err := sql.Open("pgx", connStr)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)

	pgRepo := repo.NewPGRepo(db)
	httpHandler := handler.New(pgRepo)

	http.HandleFunc("/post", httpHandler.Route(http.MethodPost))
	http.HandleFunc("/post/", httpHandler.Route(http.MethodGet, http.MethodDelete))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
