---
category: guides
subcategory: troubleshooting
confidence: low
documentType: how-to
scope: org
contentHash: 70b22a808cf0
tags: [anti-pattern, troubleshooting, faq]
source: (synthesized)
verified: 2026-07-22
synthetic: synthesized-faq
---

## Common Mistakes and Anti-Patterns

<!-- sdd-knowledge-synthesized -->

> Auto-generated from anti-patterns found across the knowledge base.

### Architecture

**From [Metrics System](../../architecture/metrics-system.md):**
The `metrics` package provides a **Prometheus-based performance metrics abstraction** that is used across the library (middleware, worker, HTTP client). All metric constructors accept a `prometheus.Registerer` to avoid global registry conflicts in multi-service deployments.

**From [Metrics System](../../architecture/metrics-system.md):**
A no-op implementation of `PerformanceMetric`. Use when metrics are disabled or in tests where you do not want Prometheus counters incremented.

### Code-conventions

**From [Anti-Patterns (failed approaches)](../../code-conventions/anti-patterns-failed-approaches.md):**
<!-- Add failed approaches here. Each anti-pattern should include:
- **type**: anti-pattern
- **discovered**: YYYY-MM-DD

### Security

**From [Security — Crypto and Auth Patterns](../../security/security-crypto-and-auth-patterns.md):**
AES-GCM provides both confidentiality and integrity (authenticated encryption). Do NOT use the raw `crypto/aes` block cipher directly in consumer code — use this wrapper.

**From [Security — Crypto and Auth Patterns](../../security/security-crypto-and-auth-patterns.md):**
- Secrets (Redis passwords, DB DSNs, AMQP URLs) are passed as constructor arguments or `Option` functions — they are never stored as package-level globals.
- TLS configuration for Redis is passed as `redis.WithTLS(cfg)` / `redis.WithClusterTLS(cfg)` — TLS is opt-in, not required.
- For production Redis TLS, use at minimum `&tls.Config{MinVersion: tls.VersionTLS12}`.

## See Also
- [crypto and auth](../../security/crypto-and-auth.md) <!-- rel:strong -->
- [metrics system](../../architecture/metrics-system.md) <!-- rel:strong -->
- [overview](../../architecture/overview.md) <!-- rel:related -->
- [constraints](../../code-conventions/constraints.md) <!-- rel:related -->
- [anti patterns failed approaches](../../code-conventions/code-style/anti-patterns-failed-approaches.md) <!-- rel:related -->
