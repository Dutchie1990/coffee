package services_test

import (
	"coffee/coffee-server/mocks"
	"coffee/coffee-server/services"
	"database/sql"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("Coffee Service", func() {
	var (
		mockDB        *mocks.DBInterface
		coffeeService *services.Coffee
	)

	BeforeEach(func() {
		mockDB = new(mocks.DBInterface)
		coffeeService = &services.Coffee{}
	})

	AfterEach(func() {
		mockDB.AssertExpectations(GinkgoT()) // Verify that all expectations were met
	})

	Describe("GetAllCoffees", func() {

		It("should return an error if the query fails", func() {
			mockDB.On("QueryContext", mock.Anything, "SELECT id, name, roast, image, region, price, grind_unit, created_at, updated_at FROM coffees", []interface{}{}).Return(nil, errors.New("query failed"))

			coffees, err := coffeeService.GetAllCoffees()
			Expect(err).To(HaveOccurred())
			Expect(coffees).To(BeNil())
		})
	})

	Describe("CreateCoffee", func() {
		It("should create a new coffee and return it", func() {
			newCoffee := services.Coffee{Name: "Mocha", Roast: "Medium", Image: "image3.png", Region: "Ethiopia", Price: 15.0, GrindUnit: 1}

			mockDB.On("ExecContext", mock.Anything, "INSERT INTO coffees(name, roast, image, region, price, grind_unit, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) returning *",
				newCoffee.Name, newCoffee.Roast, newCoffee.Image, newCoffee.Region, newCoffee.Price, newCoffee.GrindUnit, mock.Anything, mock.Anything).Return(sql.Result(nil), nil)

			createdCoffee, err := coffeeService.CreateCoffee(newCoffee)
			Expect(err).To(BeNil())
			Expect(createdCoffee).To(Equal(&newCoffee))
		})

		It("should return an error if the insert fails", func() {
			newCoffee := services.Coffee{Name: "Mocha", Roast: "Medium", Image: "image3.png", Region: "Ethiopia", Price: 15.0, GrindUnit: 1}

			mockDB.On("ExecContext", mock.Anything, "INSERT INTO coffees(name, roast, image, region, price, grind_unit, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) returning *",
				newCoffee.Name, newCoffee.Roast, newCoffee.Image, newCoffee.Region, newCoffee.Price, newCoffee.GrindUnit, mock.Anything, mock.Anything).Return(nil, errors.New("insert failed"))

			createdCoffee, err := coffeeService.CreateCoffee(newCoffee)
			Expect(err).To(HaveOccurred())
			Expect(createdCoffee).To(BeNil())
		})
	})

	Describe("GetCoffeesById", func() {
		It("should return a coffee by ID", func() {
			coffee := services.Coffee{ID: "1", Name: "Espresso", Roast: "Dark", Image: "image1.png", Region: "Brazil", Price: 10.0, GrindUnit: 1}

			mockDB.On("QueryRowContext", mock.Anything, "SELECT id, name, roast, image, region, price, grind_unit, created_at, updated_at FROM coffees WHERE id=$1", "1").Return(&sql.Row{})
			mockDB.On("Row.Scan").Return(nil)

			result, err := coffeeService.GetCoffeesById("1")
			Expect(err).To(BeNil())
			Expect(result).To(Equal(&coffee))
		})

		It("should return an error if coffee not found", func() {
			mockDB.On("QueryRowContext", mock.Anything, "SELECT id, name, roast, image, region, price, grind_unit, created_at, updated_at FROM coffees WHERE id=$1", "1").Return(&sql.Row{})
			mockDB.On("Row.Scan").Return(errors.New("no rows in result set"))

			result, err := coffeeService.GetCoffeesById("1")
			Expect(err).To(HaveOccurred())
			Expect(result).To(BeNil())
		})
	})

	Describe("UpdateCoffee", func() {
		It("should update an existing coffee and return it", func() {
			coffee := services.Coffee{ID: "1", Name: "Espresso", Roast: "Dark", Image: "image1.png", Region: "Brazil", Price: 10.0, GrindUnit: 1}

			mockDB.On("ExecContext", mock.Anything, "UPDATE coffees SET name = $1, roast = $2, image = $3, region = $4, price = $5, grind_unit = $6, updated_at = $7 WHERE id = $8 returning *",
				coffee.Name, coffee.Roast, coffee.Image, coffee.Region, coffee.Price, coffee.GrindUnit, mock.Anything, "1").Return(sql.Result(nil), nil)

			updatedCoffee, err := coffeeService.UpdateCoffee("1", coffee)
			Expect(err).To(BeNil())
			Expect(updatedCoffee).To(Equal(&coffee))
		})

		It("should return an error if the update fails", func() {
			coffee := services.Coffee{ID: "1", Name: "Espresso", Roast: "Dark", Image: "image1.png", Region: "Brazil", Price: 10.0, GrindUnit: 1}

			mockDB.On("ExecContext", mock.Anything, "UPDATE coffees SET name = $1, roast = $2, image = $3, region = $4, price = $5, grind_unit = $6, updated_at = $7 WHERE id = $8 returning *",
				coffee.Name, coffee.Roast, coffee.Image, coffee.Region, coffee.Price, coffee.GrindUnit, mock.Anything, "1").Return(nil, errors.New("update failed"))

			updatedCoffee, err := coffeeService.UpdateCoffee("1", coffee)
			Expect(err).To(HaveOccurred())
			Expect(updatedCoffee).To(BeNil())
		})
	})

	Describe("DeleteCoffee", func() {
		It("should delete a coffee by ID", func() {
			mockDB.On("ExecContext", mock.Anything, "DELETE FROM coffees WHERE id = $1", "1").Return(sql.Result(nil), nil)

			err := coffeeService.DeleteCoffee("1")
			Expect(err).To(BeNil())
		})

		It("should return an error if the deletion fails", func() {
			mockDB.On("ExecContext", mock.Anything, "DELETE FROM coffees WHERE id = $1", "1").Return(nil, errors.New("delete failed"))

			err := coffeeService.DeleteCoffee("1")
			Expect(err).To(HaveOccurred())
		})
	})
})
