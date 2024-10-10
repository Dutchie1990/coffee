package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"coffee/coffee-server/controllers"
	"coffee/coffee-server/mocks"
	"coffee/coffee-server/services"
)

var (
	mockedCoffee *mocks.CoffeeService
	recorder     *httptest.ResponseRecorder
	request      *http.Request
)

var _ = Describe("Coffee controller", Label("unit"), func() {
	BeforeEach(func() {
		recorder = httptest.NewRecorder()

		mockedCoffee = new(mocks.CoffeeService)
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
			mockedCoffee.On("GetAllCoffees").Return(mockCoffees, nil)

			controllers.GetAllCoffees(recorder, request, mockedCoffee)

			Expect(recorder.Code).To(Equal(http.StatusOK))

			var response map[string][]services.Coffee
			err := json.Unmarshal(recorder.Body.Bytes(), &response)
			Expect(err).NotTo(HaveOccurred())

			Expect(response["coffees"]).To(HaveLen(3))
			Expect(response["coffees"][0].Name).To(Equal("Latte"))
			Expect(response["coffees"][1].Name).To(Equal("Espresso"))
		})

		It("should log the error and not write a response", func() {
			mockedCoffee.On("GetAllCoffees").Return(nil, errors.New("New database error"))
			controllers.GetAllCoffees(recorder, request, mockedCoffee)

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

			mockCoffee := &services.Coffee{ // Use & to create a pointer to Coffee
				Name:      "Latte",
				Roast:     "Light",
				Image:     "image2.png",
				Region:    "Colombia",
				Price:     12.0,
				GrindUnit: 1,
			}
			mockedCoffee.On("GetCoffeesById", "").Return(mockCoffee, nil)

			controllers.GetCoffeesById(recorder, request, mockedCoffee)

			Expect(recorder.Code).To(Equal(http.StatusOK))

			var response map[string]services.Coffee
			err := json.Unmarshal(recorder.Body.Bytes(), &response)
			Expect(err).NotTo(HaveOccurred())

			Expect(response).To(HaveLen(1))
			Expect(response["coffee"].Name).To(Equal("Latte"))
			Expect(response["coffee"].Price).To(Equal(float32(12.0)))
		})
		It("Return error if not coffee not found", func() {
			mockedCoffee.On("GetCoffeesById", "").Return(nil, errors.New("The coffee is not found"))

			controllers.GetCoffeesById(recorder, request, mockedCoffee)

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

	Describe("CreateCoffee", func() {
		var (
			mockCoffeeJson []byte
			mockCoffee     services.Coffee
		)

		It("Create coffee should be succesfull", func() {
			mockCoffee = services.Coffee{
				ID:        "",
				Name:      "Latte",
				Roast:     "Light",
				Image:     "image2.png",
				Region:    "Colombia",
				Price:     12.0,
				GrindUnit: 1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			// Marshal the struct to JSON
			mockCoffeeJson, _ = json.Marshal(mockCoffee)
			request, _ = http.NewRequest(http.MethodPost, "/api/v1/coffees/coffee", bytes.NewBuffer(mockCoffeeJson))
			request.Header.Set("Content-Type", "application/json")
			mockedCoffee.On("CreateCoffee", mock.MatchedBy(func(c services.Coffee) bool {
				return c.Name == mockCoffee.Name &&
					c.Roast == mockCoffee.Roast &&
					c.Image == mockCoffee.Image &&
					c.Region == mockCoffee.Region &&
					c.Price == mockCoffee.Price &&
					c.GrindUnit == mockCoffee.GrindUnit
			})).Return(&mockCoffee, nil)

			controllers.CreateCoffee(recorder, request, mockedCoffee)
			Expect(recorder.Code).To(Equal(http.StatusOK))

			mockedCoffee.AssertCalled(GinkgoT(), "CreateCoffee", mock.MatchedBy(func(c services.Coffee) bool {
				return c.Name == mockCoffee.Name &&
					c.Roast == mockCoffee.Roast &&
					c.Image == mockCoffee.Image &&
					c.Region == mockCoffee.Region &&
					c.Price == mockCoffee.Price &&
					c.GrindUnit == mockCoffee.GrindUnit
			}))
		})
		It("Create coffee should be fails return error", func() {
			mockCoffee = services.Coffee{ // Use & to create a pointer to Coffee
				ID:        "",
				Name:      "Latte",
				Roast:     "Light",
				Image:     "image2.png",
				Region:    "Colombia",
				Price:     12.0,
				GrindUnit: 1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			// Marshal the struct to JSON
			mockCoffeeJson, _ = json.Marshal(mockCoffee)
			request, _ = http.NewRequest(http.MethodPost, "/api/v1/coffees/coffee", bytes.NewBuffer(mockCoffeeJson))
			request.Header.Set("Content-Type", "application/json")
			mockedCoffee.On("CreateCoffee", mock.MatchedBy(func(c services.Coffee) bool {
				return c.Name == mockCoffee.Name &&
					c.Roast == mockCoffee.Roast &&
					c.Image == mockCoffee.Image &&
					c.Region == mockCoffee.Region &&
					c.Price == mockCoffee.Price &&
					c.GrindUnit == mockCoffee.GrindUnit
			})).Return(nil, errors.New("Database error"))

			controllers.CreateCoffee(recorder, request, mockedCoffee)
			Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

			responseBody := recorder.Body.String()
			Expect(responseBody).To(ContainSubstring("Database error"))
		})
		It("should return an error when the JSON decoder fails", func() {
			invalidJson := `{"Name": "Latte", "Roast": "Light", "Price": "invalid_number"}`

			request, _ := http.NewRequest(http.MethodPost, "/api/v1/coffees/coffee", bytes.NewBuffer([]byte(invalidJson)))
			request.Header.Set("Content-Type", "application/json")

			// Call the handler function with the invalid request
			controllers.CreateCoffee(recorder, request, mockedCoffee)

			// Assert that the status code is 500 Internal Server Error
			Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

			responseBody := recorder.Body.String()
			Expect(responseBody).To(ContainSubstring("cannot unmarshal string into Go struct field Coffee.price of type float32"))
		})
	})
	Describe("UpdateCoffee", func() {
		var (
			mockCoffeeJson []byte
			mockCoffee     services.Coffee
		)
		It("Should successfully update a coffee", func() {
			mockCoffee = services.Coffee{
				ID:        "",
				Name:      "Latte",
				Roast:     "Light",
				Image:     "image2.png",
				Region:    "Colombia",
				Price:     12.0,
				GrindUnit: 1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			// Marshal the struct to JSON
			mockCoffeeJson, _ = json.Marshal(mockCoffee)
			request, _ = http.NewRequest(http.MethodPut, "/api/v1/coffees/coffee/12345", bytes.NewBuffer(mockCoffeeJson))
			request.Header.Set("Content-Type", "application/json")
			mockedCoffee.On("UpdateCoffee", "", mock.MatchedBy(func(c services.Coffee) bool {
				return c.Name == mockCoffee.Name &&
					c.Roast == mockCoffee.Roast &&
					c.Image == mockCoffee.Image &&
					c.Region == mockCoffee.Region &&
					c.Price == mockCoffee.Price &&
					c.GrindUnit == mockCoffee.GrindUnit
			})).Return(&mockCoffee, nil)

			controllers.UpdateCoffeeById(recorder, request, mockedCoffee)

			Expect(recorder.Code).To(Equal(http.StatusOK))

			mockedCoffee.AssertCalled(GinkgoT(), "UpdateCoffee", "", mock.MatchedBy(func(c services.Coffee) bool {
				return c.Name == mockCoffee.Name &&
					c.Roast == mockCoffee.Roast &&
					c.Image == mockCoffee.Image &&
					c.Region == mockCoffee.Region &&
					c.Price == mockCoffee.Price &&
					c.GrindUnit == mockCoffee.GrindUnit
			}))
		})

		It("Should fail when coffee is invalid", func() {
			invalidJson := `{"Name": "Latte", "Roast": "Light", "Price": "invalid_number"}`

			request, _ = http.NewRequest(http.MethodPut, "/api/v1/coffees/coffee/12345", bytes.NewBuffer([]byte(invalidJson)))
			request.Header.Set("Content-Type", "application/json")

			controllers.UpdateCoffeeById(recorder, request, mockedCoffee)

			Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			responseBody := recorder.Body.String()
			Expect(responseBody).To(ContainSubstring("cannot unmarshal string into Go struct field Coffee.price of type float32"))

		})
		It("Should fail at database error", func() {
			mockCoffee = services.Coffee{
				ID:        "",
				Name:      "Latte",
				Roast:     "Light",
				Image:     "image2.png",
				Region:    "Colombia",
				Price:     12.0,
				GrindUnit: 1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			// Marshal the struct to JSON
			mockCoffeeJson, _ = json.Marshal(mockCoffee)
			request, _ = http.NewRequest(http.MethodPut, "/api/v1/coffees/coffee/12345", bytes.NewBuffer(mockCoffeeJson))
			request.Header.Set("Content-Type", "application/json")
			mockedCoffee.On("UpdateCoffee", "", mock.MatchedBy(func(c services.Coffee) bool {
				return c.Name == mockCoffee.Name &&
					c.Roast == mockCoffee.Roast &&
					c.Image == mockCoffee.Image &&
					c.Region == mockCoffee.Region &&
					c.Price == mockCoffee.Price &&
					c.GrindUnit == mockCoffee.GrindUnit
			})).Return(nil, errors.New("Database error"))

			controllers.UpdateCoffeeById(recorder, request, mockedCoffee)

			Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

			responseBody := recorder.Body.String()
			Expect(responseBody).To(ContainSubstring("Database error"))

		})
	})
	Describe("DeleteCoffe", func() {

		BeforeEach(func() {
			request, _ = http.NewRequest(http.MethodDelete, "/api/v1/coffees/coffee/12345", nil)
		})
		It("Should succesfull delete coffee", func() {
			mockedCoffee.On("DeleteCoffee", "").Return(nil)
			controllers.DeleteCoffee(recorder, request, mockedCoffee)
			Expect(recorder.Code).To(Equal(http.StatusOK))
		})
		It("Should fail on database error", func() {
			mockedCoffee.On("DeleteCoffee", "").Return(errors.New("Database error"))
			controllers.DeleteCoffee(recorder, request, mockedCoffee)
			Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

			responseBody := recorder.Body.String()
			Expect(responseBody).To(ContainSubstring("Database error"))
		})
	})
})
