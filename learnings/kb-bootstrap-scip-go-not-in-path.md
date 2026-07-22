---
title: scip-go Unavailable When Go Not in bash Login PATH
date: 2026-07-22
pr: TBD
area: [tooling, kb-bootstrap]
files: [go.mod]
symptom: SCIP indexing step produces "go: command not found" when running inside the pod even though go is installed
tags: [scip, go, path, kb-bootstrap, tooling]
summary: The bash login shell in the KB-bootstrap pod does not have Go in PATH, so scip-go install silently fails; sdd-knowledge falls back to AST path normally.
---

## The pattern

During KB bootstrap, step 2.5 attempts to install `scip-go` via `go install`. The command runs inside `bash -lc '...'` (login shell), but in this pod Go is not in the login-shell PATH. The install silently fails, no `index.scip` is produced, and `sdd-knowledge --full` falls through to AST-based analysis.

## The rule

This is expected and acceptable for the bootstrap pod — AST analysis still produces a valid KB. If SCIP-enriched call edges are needed, the repo's CI (which has Go in PATH) can run `scip-go` as part of `knowledge-sync.yml`.

## Detection

`bash -lc 'go version'` returns non-zero or "command not found" in the pod.

## Related

- [architecture/call-graph.md](../knowledge/architecture/call-graph.md)
