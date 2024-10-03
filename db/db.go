package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
)

type DB struct {
	DB DBInterface
}

const (
	maxOpenDbConn = 10
	maxIdleDbConn = 5
	maxDbLifetime = 5 * time.Minute
)

var SqlOpen = sql.Open

func ConnectPostgres(dsn string) (*DB, error) {
	d, err := SqlOpen("pgx", dsn)
	if err != nil {
		return nil, err
	}

	d.SetMaxOpenConns(maxOpenDbConn)
	d.SetMaxIdleConns(maxIdleDbConn)
	d.SetConnMaxLifetime(maxDbLifetime)

	err = testDB(d)
	if err != nil {
		return nil, err
	}

	return &DB{DB: d}, nil
}

func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		fmt.Println("Error", err)
		return err
	}
	fmt.Println("*** Pinged database successfully! ***")
	return nil
}
