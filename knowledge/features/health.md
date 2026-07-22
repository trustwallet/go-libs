---
title: Health Check Feature
category: features
tags: [health,http,healthcheck]
confidence: high
source: source-analysis
updated: 2026-07-22
---

# Health

<!-- sdd-knowledge-generated -->

## Overview

- **Files**: 2
- **Symbols**: 10
- **Controllers**: WithHealthCheckRoute, WithReadinessCheckRoute

## Files

- `health/http_test.go`
- `health/http.go` — CheckFunc, Option, server, WithHealthCheckRoute, WithReadinessCheckRoute, WithPort, WithHealthChecks, WithReadinessChecks, handle, StartHealthCheckServer

## Architecture

### Layers

**Controller**: `WithHealthCheckRoute`, `WithReadinessCheckRoute`

**Other**: `CheckFunc`, `Option`, `server`, `WithPort`, `WithHealthChecks`, `WithReadinessChecks`, `handle`, `StartHealthCheckServer`

## Class Diagram

```mermaid
classDiagram
  class server {
  }
```

## External Dependencies

- `github.com`

## Minimum Viable Specification

> Auto-generated specification for the **Health** feature.

**Key Types**: server

## See Also
- [config](../libs/config.md) <!-- rel:strong -->
- [call graph](../architecture/call-graph.md) <!-- rel:strong -->
- [middleware health bridge](../architecture/middleware-health-bridge.md) <!-- rel:strong -->
- [overview](../architecture/overview.md) <!-- rel:strong -->
- [project structure](../architecture/project-structure.md) <!-- rel:related -->
