---
title: Message Queue (AMQP/RabbitMQ) Pattern
category: patterns
tags: [mq, rabbitmq, amqp, publish, consume, reconnect]
confidence: high
source: mq/mq.go, mq/consumer.go, mq/exchange.go
updated: 2026-07-22
---

# Message Queue (AMQP/RabbitMQ) Pattern

## Overview

The `mq` package is a thin wrapper over the `streadway/amqp` library, providing auto-reconnect, queue/exchange declaration, and publisher/consumer helpers.

## Connection and Setup

```go
client, err := mq.Connect(amqpURL,
    mq.WithConnCheckTimeout(30 * time.Second),
)
defer client.Close()

// Declare infrastructure
queue := client.InitQueue("my.queue")
exchange := client.InitExchange("my.exchange")
```

## Publishing (`publish` god-node)

`publish` is the internal function (fan-in: 3 callers, fan-out: 1 AMQP publish). It is called by all publisher methods in the package. The `Exchange` and `Queue` types both delegate to it.

```go
err := queue.Publish(ctx, message)
err := exchange.Publish(ctx, routingKey, message)
```

## Consuming

```go
consumer := client.InitConsumer("my.queue")
err := consumer.Start(ctx, func(msg mq.Message) error {
    // process message
    return nil
})
```

## Auto-Reconnect

The client attempts reconnection up to `reconnectionAttemptsNum = 5` times with a `reconnectionTimeout = 30s` delay between attempts. If all attempts fail, the error is propagated. A `connCheckTimeout` (default 10s) governs how long to wait before declaring the connection unhealthy.

## See Also
- [mq.md](../features/mq.md)
- [architecture/overview.md](../architecture/overview.md)
- [dead code candidates](../architecture/dead-code-candidates.md) <!-- rel:strong -->
- [client new req builder](../architecture/client-new-req-builder.md) <!-- rel:strong -->
- [crypto and auth](../security/crypto-and-auth.md) <!-- rel:related -->
