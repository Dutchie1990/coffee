package main_test

import (
	"bytes"
	"coffee/coffee-server/db"
	"coffee/coffee-server/services"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/lpernett/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type Config struct {
	Port string
}

type Application struct {
	Config Config
	Models services.Models
}

var (
	dbConn *db.DB  // To hold the test database connection
	sqlDB  *sql.DB // The raw SQL database object
)

var _ = BeforeSuite(func() {
	// Load test-specific environment variables
	err := godotenv.Load("../../.env")
	Expect(err).NotTo(HaveOccurred())

	// Connect to test DB
	dsn := os.Getenv("DSN")
	dbConn, err = db.ConnectPostgres(dsn)
	Expect(err).NotTo(HaveOccurred())

	var ok bool
	sqlDB, ok = dbConn.DB.(*sql.DB)
	Expect(ok).To(BeTrue(), "dbConn.DB is not a *sql.DB")
})

var _ = AfterSuite(func() {
	sqlDB.Close()
})

var _ = Describe("Main", Label("E2E"), func() {
	BeforeEach(func() {
		// Clean up database or reset state before each test
		_, err := sqlDB.Exec("DELETE FROM coffees")
		Expect(err).To(BeNil())
	})
	It("should return all coffees with status 200", func() {
		// Insert a coffee into the database
		_, err := sqlDB.Exec("INSERT INTO coffees (id, name, roast, image, region, price, grind_unit) VALUES ('550e8400-e29b-41d4-a716-446655440000','Espresso', 'Dark', 'image1.png', 'Brazil', 10.0, 1)")
		Expect(err).To(BeNil())

		// Make a GET request to the API to fetch the coffees
		res, err := http.Get("http://localhost:8080/api/v1/coffees")
		Expect(err).NotTo(HaveOccurred())
		defer res.Body.Close()

		// Ensure that the response status code is 200 OK
		Expect(res.StatusCode).To(Equal(http.StatusOK))

		// Parse the JSON response body
		var response map[string][]services.Coffee
		err = json.NewDecoder(res.Body).Decode(&response)
		Expect(err).NotTo(HaveOccurred())

		// Ensure that the response contains at least one coffee
		Expect(response["coffees"]).To(HaveLen(1))

		// Validate the content of the returned coffee
		coffee := response["coffees"][0]
		Expect(coffee.Name).To(Equal("Espresso"))
		Expect(coffee.Roast).To(Equal("Dark"))
		Expect(coffee.Region).To(Equal("Brazil"))
		Expect(float64(coffee.Price)).To(Equal(10.0))
		Expect(int(coffee.GrindUnit)).To(Equal(1))
	})
	It("should return coffee by id with status 200", func() {
		// Insert a coffee into the database
		_, err := sqlDB.Exec("INSERT INTO coffees (id, name, roast, image, region, price, grind_unit) VALUES ('550e8400-e29b-41d4-a716-446655440000','Espresso', 'Dark', 'image1.png', 'Brazil', 10.0, 1)")
		Expect(err).To(BeNil())

		// Make a GET request to the API to fetch the coffees
		res, err := http.Get("http://localhost:8080/api/v1/coffees/coffee/550e8400-e29b-41d4-a716-446655440000")
		Expect(err).NotTo(HaveOccurred())
		defer res.Body.Close()

		// Ensure that the response status code is 200 OK
		Expect(res.StatusCode).To(Equal(http.StatusOK))

		var response map[string]services.Coffee
		err = json.NewDecoder(res.Body).Decode(&response)
		Expect(err).NotTo(HaveOccurred())

		// Ensure that the response contains one key "coffee"
		Expect(response).To(HaveKey("coffee"))

		// Access the coffee from the response map
		coffee := response["coffee"]

		// Validate the content of the returned coffee
		Expect(coffee.Name).To(Equal("Espresso"))
		Expect(coffee.Roast).To(Equal("Dark"))
		Expect(coffee.Region).To(Equal("Brazil"))
		Expect(float64(coffee.Price)).To(Equal(10.0))
		Expect(int(coffee.GrindUnit)).To(Equal(1))
	})
	It("should return coffee by id with status 200", func() {
		coffee := services.Coffee{
			Name:      "Test-coffee",
			Roast:     "Test-roast",
			Region:    "Test-region",
			Image:     "Test-image",
			Price:     12.50,
			GrindUnit: 1,
		}

		coffeeJson, err := json.Marshal(coffee)
		Expect(err).NotTo(HaveOccurred())

		// Make a GET request to the API to fetch the coffees
		res, err := http.Post("http://localhost:8080/api/v1/coffees/coffee", "application/json", bytes.NewBuffer(coffeeJson))
		Expect(err).NotTo(HaveOccurred())
		defer res.Body.Close()

		// Ensure that the response status code is 200 OK
		Expect(res.StatusCode).To(Equal(http.StatusOK))

		rows, err := sqlDB.Query(`SELECT * FROM coffees`)
		Expect(err).NotTo(HaveOccurred())
		defer rows.Close()

		// Iterate over the results and validate
		for rows.Next() {
			var ID, name, roast, region, image string
			var price float64
			var grindUnit int
			var created_at, updated_at time.Time

			// Scan the result into the respective variables
			err := rows.Scan(&ID, &name, &roast, &region, &image, &price, &grindUnit, &created_at, &updated_at)
			Expect(err).NotTo(HaveOccurred())

			// Validate the content of the returned coffee
			Expect(name).To(Equal("Test-coffee"))
			Expect(roast).To(Equal("Test-roast"))
			Expect(region).To(Equal("Test-region"))
			Expect(image).To(Equal("Test-image"))
			Expect(float64(price)).To(Equal(12.5))
			Expect(int(grindUnit)).To(Equal(1))
		}
	})
})
