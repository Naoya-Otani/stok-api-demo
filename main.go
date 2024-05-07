package main

import (
	"log"
	"net/http"

	"example.com/stok-api-demo/api"
	"example.com/stok-api-demo/config"
	"example.com/stok-api-demo/db"
	"example.com/stok-api-demo/middleware"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	config.LoadEnv()
	db, err := db.NewDBConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	loggedRouter := middleware.Logger(http.HandlerFunc(api.HandleProducts(db)))
	http.Handle("/products", loggedRouter)

	log.Fatal(http.ListenAndServe(":8080", nil))
}