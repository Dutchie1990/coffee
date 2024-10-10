package services_test

import (
	"coffee/coffee-server/services"
	"database/sql"
	"log"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
)

var (
	db            *sql.DB
	coffeeService services.CoffeeService
)

var _ = BeforeSuite(func() {
	dsn := "host=localhost port=5432 user=root password=secret dbname=coffee sslmode=disable timezone=UTC connect_timeout=5"
	var err error
	// Set up your database connection here
	connStr := dsn
	db, err = sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize the Models struct with the database connection
	models := services.New(db)
	coffeeService = models.Coffee

})

var _ = AfterSuite(func() {
	if db != nil {
		db.Close()
	}
})

var _ = Describe("Coffee Service", Label("integration"), func() {
	BeforeEach(func() {
		// Clean up database or reset state before each test
		_, err := db.Exec("DELETE FROM coffees")
		Expect(err).To(BeNil())
	})

	Describe("GetAllCoffees", func() {
		It("should return an empty slice if there are no coffees", func() {
			coffees, err := coffeeService.GetAllCoffees()
			Expect(err).To(BeNil())
			Expect(coffees).To(BeEmpty())
		})

		It("should return all coffees", func() {
			_, err := db.Exec("INSERT INTO coffees (id, name, roast, image, region, price, grind_unit) VALUES ('550e8400-e29b-41d4-a716-446655440000','Espresso', 'Dark', 'image1.png', 'Brazil', 10.0, 1)")
			Expect(err).To(BeNil())

			coffees, err := coffeeService.GetAllCoffees()
			Expect(err).To(BeNil())
			Expect(coffees).To(HaveLen(1))
		})
	})

	Describe("CreateCoffee", func() {
		It("should create a new coffee and return it", func() {
			newCoffee := services.Coffee{Name: "Mocha", Roast: "Medium", Image: "image3.png", Region: "Ethiopia", Price: 15.0, GrindUnit: 1}

			createdCoffee, err := coffeeService.CreateCoffee(newCoffee)
			Expect(err).To(BeNil())
			Expect(createdCoffee).To(Equal(&newCoffee))
		})
	})

	Describe("GetCoffeesById", func() {
		It("should return a coffee by ID", func() {
			// Insert a coffee into the database
			_, err := db.Exec("INSERT INTO coffees (id, name, roast, image, region, price, grind_unit) VALUES ('550e8400-e29b-41d4-a716-446655440000', 'Espresso', 'Dark', 'image1.png', 'Brazil', 10.0, 1)")
			Expect(err).To(BeNil())

			result, err := coffeeService.GetCoffeesById("550e8400-e29b-41d4-a716-446655440000")
			Expect(err).To(BeNil())
			Expect(result.ID).To(Equal("550e8400-e29b-41d4-a716-446655440000"))
		})

		It("should return an error if coffee not found", func() {
			result, err := coffeeService.GetCoffeesById("nonexistent")
			Expect(err).To(HaveOccurred())
			Expect(result).To(BeNil())
		})
	})

	Describe("UpdateCoffee", func() {
		It("should update an existing coffee and return it", func() {
			// Insert coffee to be updated
			_, err := db.Exec("INSERT INTO coffees (id, name, roast, image, region, price, grind_unit) VALUES ('550e8400-e29b-41d4-a716-446655440000', 'Espresso', 'Dark', 'image1.png', 'Brazil', 10.0, 1)")
			Expect(err).To(BeNil())

			coffee := services.Coffee{Name: "Latte", Roast: "Light", Image: "image2.png", Region: "Colombia", Price: 12.0, GrindUnit: 1}
			updatedCoffee, err := coffeeService.UpdateCoffee("550e8400-e29b-41d4-a716-446655440000", coffee)
			Expect(err).To(BeNil())
			Expect(updatedCoffee.Name).To(Equal("Latte"))
		})
	})

	Describe("DeleteCoffee", func() {
		It("should delete a coffee by ID", func() {
			// Insert a coffee to delete
			_, err := db.Exec("INSERT INTO coffees (id, name, roast, image, region, price, grind_unit) VALUES ('550e8400-e29b-41d4-a716-446655440000', 'Espresso', 'Dark', 'image1.png', 'Brazil', 10.0, 1)")
			Expect(err).To(BeNil())

			err = coffeeService.DeleteCoffee("550e8400-e29b-41d4-a716-446655440000")
			Expect(err).To(BeNil())

			coffees, err := coffeeService.GetAllCoffees()
			Expect(err).To(BeNil())
			Expect(coffees).To(BeEmpty())
		})
	})
})
