package router

import (
	"coffee/coffee-server/controllers"
	"coffee/coffee-server/services"
	"net/http"
)

func CoffeeHandler(coffeeService services.CoffeeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		controllers.GetAllCoffees(w, r, coffeeService)
	}
}
func CoffeeByIdHandler(coffeeService services.CoffeeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		controllers.GetCoffeesById(w, r, coffeeService)
	}
}
func CreateCoffeeHandler(coffeeService services.CoffeeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		controllers.CreateCoffee(w, r, coffeeService)
	}
}
func UpdateCoffeeHandler(coffeeService services.CoffeeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		controllers.UpdateCoffeeById(w, r, coffeeService)
	}
}
func DeleteCoffeeHandler(coffeeService services.CoffeeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		controllers.DeleteCoffee(w, r, coffeeService)
	}
}
