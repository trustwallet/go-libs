# Cache

<!-- sdd-knowledge-generated -->

## Overview

- **Files**: 3
- **Symbols**: 27

## Files

- `cache/redis/client_interface.go` — redisClient, redisClientForTest
- `cache/redis/redis_test.go`
- `cache/redis/redis.go` — Redis, Option, ClusterOption, WithTLS, WithClusterTLS, InitWithTLS, Init, InitClusterWithTLS, InitCluster, Get, GetBytes, MGet, Scan, Set, SetBytes, MSet, Delete, Watch, IsAvailable, Reconnect, SetNX, SetXX, reconnectCluster, Close, HealthCheck

## Class Diagram

```mermaid
classDiagram
  class redisClient {
    <<interface>>
  }
  class redisClientForTest {
    <<interface>>
  }
  class Redis {
  }
```

## External Dependencies

- `github.com`

## Minimum Viable Specification

> Auto-generated specification for the **Cache** feature.

**Key Types**: redisClient, redisClientForTest, Redis

## See Also
- [call graph](../architecture/call-graph.md) <!-- rel:strong -->
- [dead code candidates](../architecture/dead-code-candidates.md) <!-- rel:strong -->
- [overview](../architecture/overview.md) <!-- rel:related -->
- [crypto and auth](../security/crypto-and-auth.md) <!-- rel:related -->
- [config](../libs/config.md) <!-- rel:related -->
