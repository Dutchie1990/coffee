package db_test

import (
	"coffee/coffee-server/db"
	"coffee/coffee-server/mocks"
	"database/sql"
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DB Connection", func() {
	var (
		mockDB = &mocks.DBInterface{}
		err    error
	)

	BeforeEach(func() {
		mockDB = new(mocks.DBInterface)

		// Set up expectations for the mock methods
		mockDB.On("SetMaxOpenConns", 10).Return()
		mockDB.On("SetMaxIdleConns", 5).Return()
		mockDB.On("SetConnMaxLifetime", time.Minute).Return()
		mockDB.On("Ping").Return(nil) // Default to a successful Ping

		// Mock sql.Open to return a valid *sql.DB instance
		db.SqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
			return &sql.DB{}, nil // You can return a more complex struct if needed
		}
	})

	AfterEach(func() {
		// Reset sqlOpen to the original sql.Open
		db.SqlOpen = sql.Open
	})

	Describe("testDB", func() {
		// Context("when Ping succeeds", func() {
		// 	It("should return nil", func() {
		// 		mockDB.On("Ping").Return(nil)
		// 		_, err = db.ConnectPostgres("postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable")
		// 		Expect(err).To(BeNil())
		// 	})
		// })

		Context("when Ping fails", func() {
			It("should return an error", func() {
				mockDB.On("Ping").Return(errors.New("ping failed"))
				_, err = db.ConnectPostgres("postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("ping failed"))
			})
		})
	})
})
