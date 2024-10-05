package controllers

import (
	"coffee/coffee-server/helpers"
	"coffee/coffee-server/services"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

// GET /coffees

func GetAllCoffees(w http.ResponseWriter, r *http.Request, coffee services.CoffeeService) {
	all, err := coffee.GetAllCoffees()
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJson(w, err, http.StatusInternalServerError)
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelop{"coffees": all})
}

// GET /coffees/{id}

func GetCoffeesById(w http.ResponseWriter, r *http.Request, coffeeService services.CoffeeService) {
	id := chi.URLParam(r, "id")

	// Get the coffee by ID - this returns a *Coffee (pointer)
	coffeePointer, err := coffeeService.GetCoffeesById(id)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJson(w, err, http.StatusInternalServerError)
		return
	}

	// Since coffeePointer is *Coffee, we can pass it directly to the response
	helpers.WriteJson(w, http.StatusOK, helpers.Envelop{"coffee": coffeePointer})
}

// POST /coffees

func CreateCoffee(w http.ResponseWriter, r *http.Request, coffee services.CoffeeService) {
	var coffeeData services.Coffee
	err := json.NewDecoder(r.Body).Decode(&coffeeData)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJson(w, err, http.StatusInternalServerError)
		return
	}
	coffeeCreated, err := coffee.CreateCoffee(coffeeData)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJson(w, err, http.StatusInternalServerError)
		return
	}
	helpers.WriteJson(w, http.StatusOK, helpers.Envelop{"coffees": coffeeCreated})
}

func UpdateCoffeeById(w http.ResponseWriter, r *http.Request, coffee services.CoffeeService) {
	var coffeeData services.Coffee
	err := json.NewDecoder(r.Body).Decode(&coffeeData)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJson(w, err, http.StatusInternalServerError)
		return
	}

	id := chi.URLParam(r, "id")

	coffeeUpdated, err := coffee.UpdateCoffee(id, coffeeData)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJson(w, err, http.StatusInternalServerError)
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelop{"coffees": coffeeUpdated})
}

func DeleteCoffee(w http.ResponseWriter, r *http.Request, coffee services.CoffeeService) {
	id := chi.URLParam(r, "id")

	err := coffee.DeleteCoffee(id)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJson(w, err, http.StatusInternalServerError)
		return
	}
}
