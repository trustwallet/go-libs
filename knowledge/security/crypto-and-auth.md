---
title: Security — Crypto and Auth Patterns
category: security
tags: [crypto, aes, ecdsa, signing, hmac, gin, auth]
confidence: high
source: crypto/aes.go, crypto/sign.go, gin/hmac.go
updated: 2026-07-22
---

# Security — Crypto and Auth Patterns

## AES Encryption (`crypto` package)

Provides AES-GCM symmetric encryption/decryption. Key intended for encrypting sensitive data at rest or in transit between trusted services.

```go
// Encrypt
ciphertext, err := crypto.EncryptAES(key, plaintext)

// Decrypt
plaintext, err := crypto.DecryptAES(key, ciphertext)
```

AES-GCM provides both confidentiality and integrity (authenticated encryption). Do NOT use the raw `crypto/aes` block cipher directly in consumer code — use this wrapper.

## ECDSA Signing (`crypto` package)

Provides ECDSA signature creation and verification, used for signing messages between services.

```go
sig, err := crypto.Sign(privKey, message)
ok, err := crypto.Verify(pubKey, message, sig)
```

## HMAC Request Authentication (`gin` package)

`gin.HmacMiddleware(secret string)` is a Gin middleware that authenticates incoming HTTP requests using HMAC-SHA256 signatures. The signature is computed over the request body and compared against a header value. Used for service-to-service authentication over HTTP.

## Secret Handling

- Secrets (Redis passwords, DB DSNs, AMQP URLs) are passed as constructor arguments or `Option` functions — they are never stored as package-level globals.
- TLS configuration for Redis is passed as `redis.WithTLS(cfg)` / `redis.WithClusterTLS(cfg)` — TLS is opt-in, not required.
- For production Redis TLS, use at minimum `&tls.Config{MinVersion: tls.VersionTLS12}`.
- The crypto package does not hard-code any keys — callers own key management.

## See Also
- [crypto.md](../features/crypto.md)
- [gin.md](../features/gin.md)
- [overview](../architecture/overview.md) <!-- rel:strong -->
- [client new req builder](../architecture/client-new-req-builder.md) <!-- rel:related -->
- [mq pattern](../patterns/mq-pattern.md) <!-- rel:related -->
