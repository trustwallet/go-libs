---
category: architecture
subcategory: data
confidence: low
documentType: explanation
scope: repo
contentHash: 8058ea13eece
tags: [architecture]
source: architecture/data/models.md
verified: 2026-07-22
splitPartIndex: 5
splitPartTotal: 5
canonical: false
synthetic: split-part
---

## Data Models & Schemas (Part 5)

| Field | Type | Optional |
|-------|------|----------|
| `ContractAddress` | `string` | no |
| `Name` | `string` | no |
| `OriginalSymbol` | `string` | no |
| `Owner` | `string` | no |
| `Symbol` | `string` | no |
| `TotalSupply` | `string` | no |

## TokenBalance

_struct · `blockchain/binance/model.go`:65_

| Field | Type | Optional |
|-------|------|----------|
| `Free` | `string` | no |
| `Frozen` | `string` | no |
| `Locked` | `string` | no |
| `Symbol` | `string` | no |

## TransactionData

_struct · `blockchain/binance/model.go`:38_

| Field | Type | Optional |
|-------|------|----------|
| `OrderData` | `struct { Symbol string `json:"symbol"` OrderType string `json:"orderType"` Side string `json:"side"` Price string `json:"price"` Quantity string `json:"quantity"` TimeInForce string `json:"timeInForce"` OrderID string `json:"orderId"` }` | no |

## TransactionsInBlockResponse

_struct · `blockchain/binance/model.go`:12_

| Field | Type | Optional |
|-------|------|----------|
| `BlockHeight` | `int` | no |
| `Tx` | `[]Tx` | no |

## TransactionsResponse

_struct · `blockchain/binance/api/model.go`:4_

| Field | Type | Optional |
|-------|------|----------|
| `Total` | `int` | no |
| `Tx` | `[]Tx` | no |

## Tx

_struct · `blockchain/binance/api/model.go`:11_

| Field | Type | Optional |
|-------|------|----------|
| `Hash` | `string` | no |
| `BlockHeight` | `int` | no |
| `BlockTime` | `int64` | no |
| `Type` | `Type` | no |
| `Fee` | `int` | no |
| `Code` | `int` | no |
| `Source` | `int` | no |
| `Sequence` | `int` | no |
| `Memo` | `string` | no |
| `Log` | `string` | no |
| `Data` | `string` | no |
| `Asset` | `string` | no |
| `Amount` | `float64` | no |
| `FromAddr` | `string` | no |
| `ToAddr` | `string` | no |

## Tx

_struct · `blockchain/binance/model.go`:19_

| Field | Type | Optional |
|-------|------|----------|
| `TxHash` | `string` | no |
| `BlockHeight` | `int` | no |
| `TxType` | `TxType` | no |
| `TimeStamp` | `time.Time` | no |
| `FromAddr` | `interface{}` | no |
| `ToAddr` | `interface{}` | no |
| `Value` | `string` | no |
| `TxAsset` | `string` | no |
| `TxFee` | `string` | no |
| `OrderID` | `string` | no |
| `Code` | `int` | no |
| `Data` | `string` | no |
| `Memo` | `string` | no |
| `Source` | `int` | no |
| `SubTransactions` | `[]SubTransactions` | no |
| `Sequence` | `int` | no |

## worker

_struct · `worker/worker.go`:54_

| Field | Type | Optional |
|-------|------|----------|
| `name` | `string` | no |
| `workerFn` | `func() error` | no |
| `stopFn` | `func() error` | no |
| `options` | `*WorkerOptions` | no |

## WorkerOptions

_struct · `worker/options.go`:9_

| Field | Type | Optional |
|-------|------|----------|
| `Interval` | `time.Duration` | no |
| `RunImmediately` | `bool` | no |
| `RunConsequently` | `bool` | no |
| `PerformanceMetric` | `metrics.PerformanceMetric` | no |

## See Also
- [client](../../features/client.md) <!-- rel:strong -->
- [blockchain](../../features/blockchain.md) <!-- rel:strong -->
- [database](../../features/database.md) <!-- rel:strong -->
- [mq](../../features/mq.md) <!-- rel:strong -->
- [worker pattern](../../patterns/worker-pattern.md) <!-- rel:related -->
