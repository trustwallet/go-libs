---
title: Integration Testing with testy
category: tests
tags: [integration-testing, postgres, redis, testify, suite, tagged]
confidence: high
source: testy/integration_test_suite.go, testy/tagged.go
updated: 2026-07-22
---

# Integration Testing with `testy`

## Overview

The `testy` package provides two utilities for integration testing:

1. **`IntegrationTestSuite`** — a composable struct that lazy-loads real Postgres and Redis connections from environment variables.
2. **Tagged test filtering** — a mechanism to run integration tests conditionally based on build tags or environment variables.

## `IntegrationTestSuite`

Embed this struct in your testify suite:

```go
type MyTestSuite struct {
    suite.Suite
    testy.IntegrationTestSuite
}

func (s *MyTestSuite) TestSomething() {
    db := s.GetDb()      // *gorm.DB — real Postgres
    redis := s.GetRedis() // *redis.Redis — real Redis
    // ...
}
```

**Environment variables required:**
- `TEST_DB_DSN` — Postgres DSN for test database
- `TEST_REDIS_URL` — Redis URL for test instance

Connections are lazy-initialized on first call to `GetDb()` / `GetRedis()` and cached for the suite lifetime. A fatal log occurs if the connection fails.

## Tagged Test Filtering

`testy.Tagged(t, "integration")` (or similar) allows tests to be skipped unless a specific tag is active. This enables the same test binary to run unit-only or integration modes without recompilation.

## See Also
- [testy.md](../features/testy.md)
- [architecture/overview.md](../architecture/overview.md)
- [database transactions](../architecture/database-transactions.md) <!-- rel:strong -->
- [client new req builder](../architecture/client-new-req-builder.md) <!-- rel:related -->
- [metrics system](../architecture/metrics-system.md) <!-- rel:related -->
