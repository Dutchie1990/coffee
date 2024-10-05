package router

import (
	"coffee/coffee-server/services"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func Routes(coffeeService services.CoffeeService) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Get("/api/v1/coffees", CoffeeHandler(coffeeService))
	router.Get("/api/v1/coffees/coffee/{id}", CoffeeByIdHandler(coffeeService))
	router.Post("/api/v1/coffees/coffee", CreateCoffeeHandler(coffeeService))
	router.Put("/api/v1/coffees/coffee/{id}", UpdateCoffeeHandler(coffeeService))
	router.Delete("/api/v1/coffees/coffee/{id}", DeleteCoffeeHandler(coffeeService))

	return router
}
