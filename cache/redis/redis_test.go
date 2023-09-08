package redis

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type redisInitFn func(t *testing.T) (*Redis, error)

func TestInit(t *testing.T) {
	redisInitFns := []redisInitFn{redisInit, redisClusterInit}

	for _, redisInit := range redisInitFns {
		t.Run("", func(t *testing.T) {
			i, err := redisInit(t)
			assert.Nil(t, err)
			assert.NotNil(t, i)
		})
	}
}

func TestInitWithTLS(t *testing.T) {
	serverCert, err := generateCertificate()
	assert.Nil(t, err)

	mr, err := miniredis.RunTLS(&tls.Config{
		Certificates: []tls.Certificate{*serverCert},
	})
	assert.NotNil(t, mr)
	assert.Nil(t, err)

	redis, err := Init(context.TODO(), fmt.Sprintf("rediss://%s", mr.Addr()), WithTLS(&tls.Config{InsecureSkipVerify: true}))
	assert.Nil(t, err)
	assert.NotNil(t, redis)
}

func TestRedis_Set(t *testing.T) {
	redisInitFns := []redisInitFn{redisInit, redisClusterInit}
	for _, redisInit := range redisInitFns {
		t.Run("", func(t *testing.T) {
			r, err := redisInit(t)
			assert.Nil(t, err)

			testData := struct {
				Field string
				F     float64
			}{"1", 200.1}

			err = r.Set(context.TODO(), "test", testData, time.Second)
			assert.Nil(t, err)

			var newValue struct {
				Field string
				F     float64
			}
			err = r.Get(context.TODO(), "test", &newValue)
			assert.Nil(t, err)
			assert.Equal(t, testData, newValue)

			ttl := r.client.TTL(context.TODO(), "test")
			assert.Equal(t, time.Second, ttl.Val())
		})
	}
}

func TestRedis_SetBytes_GetBytes(t *testing.T) {
	redisInitFns := []redisInitFn{redisInit, redisClusterInit}
	for _, redisInit := range redisInitFns {
		t.Run("", func(t *testing.T) {
			r, err := redisInit(t)
			assert.Nil(t, err)

			tests := []struct {
				key  string
				data []byte
			}{
				{
					key:  "empty_bytes",
					data: []byte{},
				},
				{
					key:  "nil_bytes",
					data: nil,
				},
				{
					key:  "random_bytes",
					data: []byte("348ufksj@#$@agh$$!"),
				},
			}

			for _, test := range tests {
				t.Run(test.key, func(t *testing.T) {
					err := r.SetBytes(context.Background(), test.key, test.data, time.Minute)
					assert.NoError(t, err)

					got, err := r.GetBytes(context.Background(), test.key)
					assert.NoError(t, err)
					if len(test.data) == 0 { // the client returns []byte{} when we expecte []byte(nil) so we compare by len
						require.Equal(t, 0, len(got))
					} else {
						assert.Equal(t, test.data, got)
					}
				})
			}
		})
	}
}

func TestRedis_MSet_MGet(t *testing.T) {
	redisInitFns := []redisInitFn{redisInit, redisClusterInit}
	for _, redisInit := range redisInitFns {
		t.Run("", func(t *testing.T) {
			r, err := redisInit(t)
			assert.Nil(t, err)

			type testDataType struct {
				Field string
				F     float64
			}
			testData := map[string]interface{}{
				"test1": testDataType{"1", 200.1},
				"test2": testDataType{"1", 200.1},
			}

			err = r.MSet(context.TODO(), testData, time.Second)
			assert.Nil(t, err)

			keys := []string{"test1", "test2"}
			values, err := r.MGet(context.TODO(), keys...)
			assert.Nil(t, err)
			assert.Equal(t, len(keys), len(values))

			for i, key := range keys {
				var data testDataType
				err := json.Unmarshal(values[i], &data)
				assert.NoError(t, err)
				assert.Equal(t, testData[key], data)
			}

			ttl := r.client.TTL(context.TODO(), "test1")
			assert.Equal(t, time.Second, ttl.Val())
			ttl = r.client.TTL(context.TODO(), "test2")
			assert.Equal(t, time.Second, ttl.Val())
		})
	}
}

func TestRedis_Get(t *testing.T) {
	redisInitFns := []redisInitFn{redisInit, redisClusterInit}
	for _, redisInit := range redisInitFns {
		t.Run("", func(t *testing.T) {
			r, err := redisInit(t)
			assert.Nil(t, err)

			type TestStruct struct {
				Field string
				F     float64
			}
			testData := TestStruct{"1", 200.1}

			err = r.Set(context.TODO(), "test", testData, time.Second)
			assert.Nil(t, err)

			var newValue TestStruct
			err = r.Get(context.TODO(), "test", &newValue)
			assert.Nil(t, err)
			assert.Equal(t, testData, newValue)

			ttl := r.client.TTL(context.TODO(), "test")
			assert.Equal(t, time.Second, ttl.Val())

			var empty interface{}
			err = r.Get(context.TODO(), "1", empty)
			assert.Equal(t, err, ErrNotFound)
			assert.Nil(t, empty)
		})
	}
}

func TestRedis_Delete(t *testing.T) {
	redisInitFns := []redisInitFn{redisInit, redisClusterInit}
	for _, redisInit := range redisInitFns {
		t.Run("", func(t *testing.T) {
			r, err := redisInit(t)
			assert.Nil(t, err)
			err = r.Set(context.TODO(), "test", []byte{0, 1}, time.Second)
			assert.Nil(t, err)

			err = r.Delete(context.TODO(), "test")
			assert.Nil(t, err)

			var v []byte
			err = r.Get(context.TODO(), "test", &v)
			assert.NotNil(t, err)
			assert.Equal(t, string([]byte{}), string(v))
		})
	}
}

func TestRedis_Watch(t *testing.T) {
	redisInitFns := []redisInitFn{redisInit, redisClusterInit}
	for _, redisInit := range redisInitFns {
		t.Run("", func(t *testing.T) {
			r, err := redisInit(t)
			assert.Nil(t, err)

			ctx := context.TODO()
			key := "test"
			err = r.Watch(ctx, func(tx *redis.Tx) error {
				n, err := tx.Get(ctx, key).Int()
				if err != nil && err != redis.Nil {
					return err
				}

				// Actual operation (local in optimistic lock).
				n++

				// Operation is commited only if the watched keys remain unchanged.
				_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
					pipe.Set(ctx, key, fmt.Sprintf("%d", n), 0)
					return nil
				})
				return err
			}, key)
			assert.Nil(t, err)

			newValue, err := r.GetBytes(context.TODO(), key)
			assert.Nil(t, err)
			assert.Equal(t, "1", string(newValue))
		})
	}
}

func TestRedis_IsAvailable(t *testing.T) {
	redisInitFns := []redisInitFn{redisInit, redisClusterInit}
	for _, redisInit := range redisInitFns {
		t.Run("", func(t *testing.T) {
			r, err := redisInit(t)
			assert.Nil(t, err)

			assert.True(t, r.IsAvailable(context.TODO()))
		})
	}

}

func TestRedis_Reconnect(t *testing.T) {
	redisInitFns := []redisInitFn{redisInit, redisClusterInit}
	for _, redisInit := range redisInitFns {
		t.Run("", func(t *testing.T) {
			r, err := redisInit(t)
			assert.Nil(t, err)

			mr, err := miniredis.Run()
			assert.NotNil(t, mr)
			assert.Nil(t, err)

			err = r.Reconnect(context.TODO(), fmt.Sprintf("redis://%s", mr.Addr()))
			assert.NoError(t, err)
		})
	}

}

func redisInit(t *testing.T) (*Redis, error) {
	mr, err := miniredis.Run()
	assert.NotNil(t, mr)
	assert.Nil(t, err)

	c, err := Init(context.TODO(), fmt.Sprintf("redis://%s", mr.Addr()))
	assert.Nil(t, err)
	assert.NotNil(t, c)

	return c, nil
}

func redisClusterInit(t *testing.T) (*Redis, error) {
	mr, err := miniredis.Run()
	assert.NotNil(t, mr)
	assert.Nil(t, err)

	c, err := InitCluster(context.TODO(), fmt.Sprintf("redis://%s", mr.Addr()))
	assert.Nil(t, err)
	assert.NotNil(t, c)

	return c, nil
}

// generateCertificate generates self-signed untrusted certificate.
// WithTLS optiom should be used with insecureSkipVerify: true in tests.
// In real world scenario, the certificate should be trusted (insecureSkipVerify: false).
// See the proper and idiomatic way to generate a self signed certificate.
// https://go.dev/src/crypto/tls/generate_cert.go
func generateCertificate() (*tls.Certificate, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(1 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, err
	}

	certOut := new(bytes.Buffer)
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return nil, err
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, err
	}

	keyOut := new(bytes.Buffer)
	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		return nil, err
	}

	serverCert, err := tls.X509KeyPair(certOut.Bytes(), keyOut.Bytes())
	if err != nil {
		return nil, err
	}

	return &serverCert, nil
}
