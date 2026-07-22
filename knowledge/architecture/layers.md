---
title: Architectural Layers
category: architecture
tags: [layers,architecture]
confidence: high
source: source-analysis
updated: 2026-07-22
---

# Architectural Layers

<!-- sdd-knowledge-generated -->

> Layer of each module (path/role-derived) and any calls that flow **up** the stack — a layering violation. Allowed flow: ui → controller → service → repository → model.

## Layer distribution

| Layer | Symbols |
|-------|---------|
| controller | 50 |
| service | 1 |
| model | 10 |

## Violations (0)

_No upward-flowing calls detected._

## See Also
- [blockchain](../features/blockchain.md) <!-- rel:strong -->
- [gin](../features/gin.md) <!-- rel:strong -->
- [health](../features/health.md) <!-- rel:strong -->
- [metrics](../features/metrics.md) <!-- rel:related -->
- [client](../features/client.md) <!-- rel:related -->
