# Logging

<!-- sdd-knowledge-generated -->

## Overview

- **Files**: 3
- **Symbols**: 9

## Files

- `logging/formatter_strict_text.go` — TextFormatterConfig, init, NewTextFormatter
- `logging/logger_test.go`
- `logging/logger.go` — Config, init, SetLoggerConfig, GetLogger, GetLoggerForComponent, SetLogger

## Architecture

### Layers

**Config**: `TextFormatterConfig`, `Config`, `SetLoggerConfig`

**Other**: `init`, `NewTextFormatter`, `init`, `GetLogger`, `GetLoggerForComponent`, `SetLogger`

## Class Diagram

```mermaid
classDiagram
  class TextFormatterConfig {
    <<config>>
  }
```

## External Dependencies

- `github.com`
- `gotest.tools`

## Minimum Viable Specification

> Auto-generated specification for the **Logging** feature.

**Key Types**: TextFormatterConfig

## See Also
- [dependency graph](../architecture/dependency-graph.md) <!-- rel:strong -->
- [config](../libs/config.md) <!-- rel:strong -->
- [overview](../architecture/overview.md) <!-- rel:related -->
- [call graph](../architecture/call-graph.md) <!-- rel:related -->
- [models](../architecture/data/models.md) <!-- rel:related -->
