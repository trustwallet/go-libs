---
title: Go-libs Constraints and Conventions
category: code-conventions
tags: [constraints, public-repo, no-secrets, circular-imports, prometheus, testy]
confidence: high
source: CLAUDE.md
updated: 2026-07-22
---

# Go-libs Constraints and Conventions

These constraints apply to ALL work in `github.com/trustwallet/go-libs`:

## Public Repo Security

This is a **public GitHub repository**. Never commit:
- Secrets, API keys, tokens, passwords
- Private DSNs or connection strings
- Internal infrastructure details (cluster names, internal IPs, VPN details)
- Any `.env` values

## No Circular Imports Between Packages

Each package must be importable independently. Do NOT create imports between sibling packages (e.g. `client` importing `middleware`, `database` importing `metrics` directly). Dependencies flow ONE WAY: utility packages → infrastructure packages → integration packages.

## Prometheus Registry Isolation

All metric constructors (`NewPerformanceMetric`, `NewHttpServerMetric`, etc.) accept a `prometheus.Registerer` parameter. **Never** use:
- `prometheus.DefaultRegisterer` directly
- `prometheus.MustRegister` with global state
- Package-level `var _ = prometheus.Register(...)` auto-registration

This keeps consumer services' registries independent and prevents registration panics in multi-service test scenarios.

## `testy` Package is Test-Only

The `testy` package provides integration test helpers (real DB + Redis connections). Never import it from non-test (production) code. It should only appear in `_test.go` files.

## See Also
- [architecture/overview.md](../architecture/overview.md)
- [security/crypto-and-auth.md](../security/crypto-and-auth.md)
- [common mistakes and anti patterns](../guides/troubleshooting/common-mistakes-and-anti-patterns.md) <!-- rel:strong -->
- [metrics system](../architecture/metrics-system.md) <!-- rel:related -->
- [integration testing](../tests/integration-testing.md) <!-- rel:related -->
