package testy

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

const testDbDsnEnvKey = "TEST_DB_DSN"

// IntegrationTestSuite is an integration testing suite with methods
// for retrieving the real database connection.
// Just absorb the built-in IntegrationTestSuite by defining your own suite,
// you can also use it along with `testify`'s suite.
// Example:
//
//	type SomeTestSuite struct {
//		suite.Suite
//		IntegrationTestSuite
//	}
type IntegrationTestSuite struct {
	db *gorm.DB
}

// GetDb retrieves the current *gorm.DB connection, and it's lazy loaded.
func (b *IntegrationTestSuite) GetDb() *gorm.DB {
	if b.db == nil {
		db, err := NewIntegrationTestDb()
		if err != nil {
			log.Fatalln("can not connect integration test db", err)
		}
		b.db = db
	}
	return b.db
}

// NewIntegrationTestDb creates a *gorm.DB connection to a real database which is only for integration test.
// The DSN for test database connection should be set by defining the TEST_DB_DSN env.
func NewIntegrationTestDb() (*gorm.DB, error) {
	dsn, ok := os.LookupEnv(testDbDsnEnvKey)
	if !ok {
		log.Fatalln(testDbDsnEnvKey, "env not found")
	}

	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Info),
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}
