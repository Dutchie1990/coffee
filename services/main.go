package services

import (
	"database/sql"
	"time"
)

const dbTimeout = time.Second * 3

type Models struct {
	Coffee       CoffeeService
	JsonResponse JsonResponse
}

func New(dbPool *sql.DB) Models {
	return Models{
		Coffee:       &CoffeeServiceImpl{DB: dbPool}, // Initialize the concrete CoffeeService
		JsonResponse: JsonResponse{},
	}
}
