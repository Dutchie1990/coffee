package main

import (
	"coffee/coffee-server/db"
	"coffee/coffee-server/router"
	"coffee/coffee-server/services"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/lpernett/godotenv"
)

type Config struct {
	Port string
}

type Application struct {
	Config Config
	Models services.Models
}

func (app *Application) Serve() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env file")
	}
	port := os.Getenv("PORT")
	fmt.Printf(`API is listering on %s`, port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router.Routes(),
	}

	return srv.ListenAndServe()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env file")
	}

	cfg := Config{
		Port: os.Getenv("PORT"),
	}

	dsn := os.Getenv("DSN")
	dbConn, err := db.ConnectPostgres(dsn)
	if err != nil {
		log.Fatal("Error connection DB")
	}

	defer dbConn.DB.Close()

	// Assert dbConn.DB to *sql.DB
	sqlDB, ok := dbConn.DB.(*sql.DB)
	if !ok {
		fmt.Println("Error: dbConn.DB is not a *sql.DB")
		return
	}

	app := &Application{
		Config: cfg,
		Models: services.New(sqlDB),
	}

	err = app.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
