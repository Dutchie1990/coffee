package controllers

import (
	"coffee/coffee-server/helpers"
	"coffee/coffee-server/services"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

var coffee services.Coffee

// GET /coffees

func GetAllCoffees(w http.ResponseWriter, r *http.Request) {
	all, err := coffee.GetAllCoffees()
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelop{"coffees": all})
}

// GET /coffees/{id}

func GetCoffeesById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	coffee, err := coffee.GetCoffeesById(id)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelop{"coffees": coffee})
}

// POST /coffees

func CreateCoffee(w http.ResponseWriter, r *http.Request) {
	var coffeeData services.Coffee
	err := json.NewDecoder(r.Body).Decode(&coffeeData)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	coffeeCreated, err := coffeeData.CreateCoffee((coffeeData))
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	helpers.WriteJson(w, http.StatusOK, helpers.Envelop{"coffees": coffeeCreated})
}

func UpdateCoffeeById(w http.ResponseWriter, r *http.Request) {
	var coffeeData services.Coffee
	err := json.NewDecoder(r.Body).Decode(&coffeeData)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	id := chi.URLParam(r, "id")

	coffeeUpdated, err := coffeeData.UpdateCoffee(id, coffeeData)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelop{"coffees": coffeeUpdated})
}

func DeleteCoffee(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := coffee.DeleteCoffee(id)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
}
