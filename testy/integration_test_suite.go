package testy

import (
	"context"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/trustwallet/go-libs/cache/redis"
)

const (
	testDbDsnEnvKey    = "TEST_DB_DSN"
	testRedisUrlEnvKey = "TEST_REDIS_URL"
)

// IntegrationTestSuite is an integration testing suite with methods
// for retrieving the real database and redis connection.
// Just absorb the built-in IntegrationTestSuite by defining your own suite,
// you can also use it along with `testify`'s suite.
// Example:
//
//	type SomeTestSuite struct {
//		suite.Suite
//		IntegrationTestSuite
//	}
type IntegrationTestSuite struct {
	db    *gorm.DB
	redis *redis.Redis
}

// GetDb retrieves the current *gorm.DB connection, and it's lazy loaded.
func (s *IntegrationTestSuite) GetDb() *gorm.DB {
	if s.db == nil {
		db, err := NewIntegrationTestDb()
		if err != nil {
			log.Fatalln("can not connect integration test db", err)
		}
		s.db = db
	}
	return s.db
}

// GetRedis retrieves the current *redis.Redis connection, and it's lazy loaded.
func (s *IntegrationTestSuite) GetRedis() *redis.Redis {
	if s.redis == nil {
		r, err := NewIntegrationTestRedis()
		if err != nil {
			log.Fatalln("can not connect integration redis db", err)
		}
		s.redis = r
	}
	return s.redis
}

// NewIntegrationTestDb creates a *gorm.DB connection to a real database which is only for integration test.
// The DSN for test database connection should be set by defining the TEST_DB_DSN env.
func NewIntegrationTestDb() (*gorm.DB, error) {
	return gorm.Open(
		postgres.Open(MustGetTestDbDSN()),
		&gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Info),
			SkipDefaultTransaction: true,
		},
	)
}

// NewIntegrationTestRedis creates a *redis.Redis connection to a real redis pool which is only for integration test.
// The url for test redis connection should be set by defining the TEST_REDIS_URL env.
func NewIntegrationTestRedis() (*redis.Redis, error) {
	url, ok := os.LookupEnv(testRedisUrlEnvKey)
	if !ok {
		log.Fatalln(testRedisUrlEnvKey, "env not found")
	}
	return redis.Init(context.Background(), url)
}

func MustGetTestDbDSN() string {
	dsn, ok := os.LookupEnv(testDbDsnEnvKey)
	if !ok {
		log.Fatal(testDbDsnEnvKey, "env not found")
	}
	return dsn
}
