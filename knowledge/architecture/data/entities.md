# Domain Entities & Data Transfer Objects

<!-- sdd-knowledge-generated -->

## DTOs, Requests & Responses

| Name | Type | Role | File | Line |
|------|------|------|------|------|
| TransactionsResponse | class | Response | `blockchain/binance/api/model.go` | 4 |
| NodeInfoResponse | class | Response | `blockchain/binance/model.go` | 6 |
| TransactionsInBlockResponse | class | Response | `blockchain/binance/model.go` | 12 |
| Request | class | Request | `client/client.go` | 17 |
| constructHttpRequest | method | Request | `client/client_execute.go` | 63 |
| RpcRequest | class | Request | `client/jsonrpc.go` | 17 |
| RpcResponse | class | Response | `client/jsonrpc.go` | 24 |
| Query | method | Query | `client/request.go` | 73 |
| Event | class | Event | `eventer/client.go` | 19 |
| cacheResponse | class | Response | `middleware/cache.go` | 31 |
| Message | type | DTO | `mq/mq.go` | 18 |

## See Also
- [blockchain](../../features/blockchain.md) <!-- rel:strong -->
- [client](../../features/client.md) <!-- rel:strong -->
- [eventer](../../features/eventer.md) <!-- rel:related -->
- [middleware](../../features/middleware.md) <!-- rel:related -->
- [mq](../../features/mq.md) <!-- rel:weak -->
