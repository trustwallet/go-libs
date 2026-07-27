---
title: HTTP Client — NewReqBuilder (God Node)
category: architecture
tags: [http, client, builder, god-node, request-pipeline]
confidence: high
source: client/request.go, client/client_execute.go, client/jsonrpc.go
updated: 2026-07-22
---

# HTTP Client — `NewReqBuilder` (God Node)

`NewReqBuilder` is the central entry point for all outbound HTTP requests in the library (fan-in: 7 callers). It is a **builder** that constructs a `Req` value which is then executed by `Request.Execute`.

## Contract

```go
func NewReqBuilder() *ReqBuilder
```

Returns a `*ReqBuilder` with defaults: empty headers, `pathMetricEnabled: true`.

**Builder methods (all return `*ReqBuilder` for chaining):**

| Method | Purpose |
|---|---|
| `.Method(string)` | HTTP verb (GET, POST, …) |
| `.PathStatic(string)` | Static path — no template parameters |
| `.Pathf(template, values...)` | Parameterized path (use instead of PathStatic when path has vars) |
| `.Headers(map)` | Merge extra request headers |
| `.Body(any)` | JSON-encode as request body |
| `.WriteTo(any)` | JSON-decode response into this pointer |
| `.WriteRawResponseTo(*http.Response)` | Capture raw HTTP response (headers, status) without JSON decoding |
| `.Query(url.Values)` | URL query parameters |
| `.MetricName(string)` | Override the Prometheus label for this request |
| `.Build()` | Finalize and return a `*Req` (copy of internal state — safe to reuse builder) |

## Why It Is Central

Every HTTP call in the library (including JSON-RPC calls in `jsonrpc.go`, batch JSON-RPC in `jsonrpc_batch.go`, cached client calls in `clientcache.go`, and the wrapper shims in `client_wrapper.go`) ultimately constructs a `*Req` via `NewReqBuilder().…Build()` and passes it to `Request.Execute`. This is the single construction path — there is no alternate constructor.

## Callers

- `client/client_execute.go` — `Request.Execute` receives the built `Req`
- `client/jsonrpc.go` — `RpcCall`, `RpcCallRaw` build their POST requests here
- `client/jsonrpc_batch.go` — batch JSON-RPC calls
- `client/clientcache.go` — cached HTTP GET calls
- `client/client_wrapper.go` — deprecated wrapper shims (`.Get`, `.Post`, etc.)

## Important Note on `pathMetricEnabled`

Legacy wrapper functions (`Get`, `Post`, `PostRaw`) call the unexported `.pathMetricEnabled(false)` on the builder to suppress path-based metric labels (they predate the builder API). New callers should always use `NewReqBuilder()` directly and rely on the default `pathMetricEnabled: true`.

## See Also
- [client.md](../features/client.md)
- [architecture/overview.md](overview.md)
- [crypto and auth](../security/crypto-and-auth.md) <!-- rel:related -->
- [mq pattern](../patterns/mq-pattern.md) <!-- rel:related -->
- [integration testing](../tests/integration-testing.md) <!-- rel:weak -->
