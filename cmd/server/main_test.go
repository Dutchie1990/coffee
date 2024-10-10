package main_test

import (
	"coffee/coffee-server/db"
	"coffee/coffee-server/services"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"

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
	It("should return coffee by id with status 200", Label("integration"), func() {
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
})
