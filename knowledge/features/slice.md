---
title: Slice Utilities
category: features
tags: [slice,filter,partition,search,generic]
confidence: high
source: source-analysis
updated: 2026-07-22
---

# Slice

<!-- sdd-knowledge-generated -->

## Overview

- **Files**: 6
- **Symbols**: 5

## Files

- `slice/filter_test.go`
- `slice/filter.go` — Filter
- `slice/partition_test.go`
- `slice/partition.go` — Partition
- `slice/search_test.go`
- `slice/search.go` — Contains, ValueAt, Min

## Architecture

### Layers

**Middleware**: `Filter`

**Other**: `Partition`, `Contains`, `ValueAt`, `Min`

## External Dependencies

- `github.com`
- `golang.org`

## Minimum Viable Specification

> Auto-generated specification for the **Slice** feature.

**Key Types**: none

## See Also
- [overview](../architecture/overview.md) <!-- rel:strong -->
- [dependency graph](../architecture/dependency-graph.md) <!-- rel:strong -->
- [config](../libs/config.md) <!-- rel:strong -->
- [project structure](../architecture/project-structure.md) <!-- rel:related -->
- [call graph](../architecture/call-graph.md) <!-- rel:related -->
