---
title: Dependency Graph
category: architecture
tags: [dependency,graph,packages]
confidence: high
source: source-analysis
updated: 2026-07-22
---

# Dependency Graph

<!-- sdd-knowledge-generated -->

## Feature Dependencies

```mermaid
flowchart LR
  blockchain --> client
  blockchain --> eventer
  client --> blockchain
  client --> eventer
  eventer --> blockchain
  eventer --> client
  metrics --> blockchain
  metrics --> client
  metrics --> eventer
  middleware --> metrics
  mock --> blockchain
  mock --> client
  mock --> eventer
  mq --> metrics
  mq --> middleware
  testy --> cache
  worker --> metrics
  worker --> middleware
```

## External Dependencies

| Package | Import Count |
|---------|--------------|
| `github.com` | 88 |
| `gorm.io` | 9 |
| `golang.org` | 5 |
| `gotest.tools` | 3 |

## See Also
- [logging](../features/logging.md) <!-- rel:strong -->
- [gin](../features/gin.md) <!-- rel:strong -->
- [worker](../features/worker.md) <!-- rel:strong -->
- [httplib](../features/httplib.md) <!-- rel:strong -->
- [slice](../features/slice.md) <!-- rel:strong -->
