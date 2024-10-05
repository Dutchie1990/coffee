package controllers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"coffee/coffee-server/controllers"
	"coffee/coffee-server/mocks"
	"coffee/coffee-server/services"
)

var (
	mock     *mocks.CoffeeService
	recorder *httptest.ResponseRecorder
	request  *http.Request
)

var ()

var _ = Describe("Coffee controller", func() {
	BeforeEach(func() {
		recorder = httptest.NewRecorder()

		mock = new(mocks.CoffeeService)
	})

	Describe("GetAllCoffees", func() {
		BeforeEach(func() {
			request, _ = http.NewRequest(http.MethodGet, "/api/v1/coffees", nil)
		})
		It("should return all coffees with status 200", func() {
			mockCoffees := []*services.Coffee{
				{Name: "Latte", Roast: "Light", Image: "image2.png", Region: "Colombia", Price: 12.0, GrindUnit: 1},
				{Name: "Espresso", Roast: "Dark", Image: "image3.png", Region: "Italy", Price: 10.0, GrindUnit: 2},
				{Name: "Cappuccino", Roast: "Medium", Image: "image4.png", Region: "Italy", Price: 11.0, GrindUnit: 1},
			}
			mock.On("GetAllCoffees").Return(mockCoffees, nil)

			controllers.GetAllCoffees(recorder, request, mock)

			Expect(recorder.Code).To(Equal(http.StatusOK))

			var response map[string][]services.Coffee
			err := json.Unmarshal(recorder.Body.Bytes(), &response)
			Expect(err).NotTo(HaveOccurred())

			Expect(response["coffees"]).To(HaveLen(3))
			Expect(response["coffees"][0].Name).To(Equal("Latte"))
			Expect(response["coffees"][1].Name).To(Equal("Espresso"))
		})

		It("should log the error and not write a response", func() {
			mock.On("GetAllCoffees").Return(nil, errors.New("New database error"))
			controllers.GetAllCoffees(recorder, request, mock)

			Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

			// Access the response body
			responseBody := recorder.Body.String() // Get content as string
			Expect(responseBody).NotTo(BeEmpty())  // Ensure the body is not empty

			// Check the response for expected error structure
			var response map[string]any // Adjust this according to your error response structure
			err := json.Unmarshal(recorder.Body.Bytes(), &response)
			Expect(err).NotTo(HaveOccurred())

			Expect(response["error"]).To(Equal(true))
			Expect(response["message"]).To(Equal("New database error"))
		})
	})
	Describe("GetAllCoffeeById", func() {
		BeforeEach(func() {
			request, _ = http.NewRequest(http.MethodGet, "/api/v1/coffees/coffee/12345", nil)
		})
		It("Return coffee by id", func() {

			mockCoffees := &services.Coffee{ // Use & to create a pointer to Coffee
				Name:      "Latte",
				Roast:     "Light",
				Image:     "image2.png",
				Region:    "Colombia",
				Price:     12.0,
				GrindUnit: 1,
			}
			mock.On("GetCoffeesById", "").Return(mockCoffees, nil)

			controllers.GetCoffeesById(recorder, request, mock)

			Expect(recorder.Code).To(Equal(http.StatusOK))

			var response map[string]services.Coffee
			err := json.Unmarshal(recorder.Body.Bytes(), &response)
			Expect(err).NotTo(HaveOccurred())

			Expect(response).To(HaveLen(1))
			Expect(response["coffee"].Name).To(Equal("Latte"))
			Expect(response["coffee"].Price).To(Equal(float32(12.0)))
		})
		It("Return error if not coffee not found", func() {
			mock.On("GetCoffeesById", "").Return(nil, errors.New("The coffee is not found"))

			controllers.GetCoffeesById(recorder, request, mock)

			Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

			// Access the response body
			responseBody := recorder.Body.String() // Get content as string
			Expect(responseBody).NotTo(BeEmpty())  // Ensure the body is not empty

			// Check the response for expected error structure
			var response map[string]any // Adjust this according to your error response structure
			err := json.Unmarshal(recorder.Body.Bytes(), &response)
			Expect(err).NotTo(HaveOccurred())

			Expect(response["error"]).To(Equal(true))
			Expect(response["message"]).To(Equal("The coffee is not found"))
		})
	})
	// Describe("CreateCoffee", func() {})
	// Describe("UpdateCoffee", func() {})
	// Describe("DeleteCoffe", func() {})
})
