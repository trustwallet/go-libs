package gin

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"gotest.tools/assert"
)

func createTestContext(t *testing.T, w *httptest.ResponseRecorder, rawURL string, headers map[string]string) *gin.Context {
	c, _ := gin.CreateTestContext(w)

	u, err := url.Parse(rawURL)
	if err != nil {
		t.Fatal(err)
	}

	c.Request = &http.Request{
		URL:    u,
		Header: make(http.Header),
	}
	for k, v := range headers {
		c.Request.Header.Add(k, v)
	}

	return c
}

func TestHmacVerifier(t *testing.T) {
	verifier := NewHmacVerifier(
		WithHmacVerifierSigKeys("some-key", "some-other-key", "k"),
	)

	someRouteHandler := func(c *gin.Context) { c.Status(http.StatusOK) }
	signedRouteHandler := verifier.SignedHandler(someRouteHandler, func(c *gin.Context) (string, error) {
		pt := c.Query("asset") + ":" + c.Query("coin")
		return pt, nil
	})

	t.Run("unauthorized when signature not provided", func(t *testing.T) {
		w := httptest.NewRecorder()
		signedRouteHandler(createTestContext(t, w, "", nil))
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("signature verification success for all keys", func(t *testing.T) {
		asset := "123"
		coin := "60"
		plainText := asset + ":" + coin
		rawURL := fmt.Sprintf("http://does.not.matter?asset=%s&coin=%s", asset, coin)

		for _, k := range verifier.keys {
			h := hmac.New(sha256.New, k)
			h.Write([]byte(plainText))
			sig := base64.StdEncoding.EncodeToString(h.Sum(nil))

			w := httptest.NewRecorder()
			signedRouteHandler(createTestContext(t, w, rawURL, map[string]string{
				"X-REQ-SIG": sig,
			}))
			assert.Equal(t, http.StatusOK, w.Code)
		}
	})

	t.Run("signature verification failure", func(t *testing.T) {
		asset := "123"
		coin := "60"
		rawURL := fmt.Sprintf("http://does.not.matter?asset=%s&coin=%s", asset, coin)

		w := httptest.NewRecorder()
		signedRouteHandler(createTestContext(t, w, rawURL, map[string]string{
			"X-REQ-SIG": "JBdWTO5yR2GB0TOT8YcM7AjWaJaMtVrAFOYUlZRNlYg=",
		}))
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("bad request when signature cannot be extracted", func(t *testing.T) {
		h := NewHmacVerifier(
			WithHmacVerifierSigKeys("some-key", "some-other-key", "k"),
			WithHmacVerifierSigFunction(func(c *gin.Context) (string, error) {
				return "", errors.New("some error")
			}),
		).SignedHandler(someRouteHandler, func(c *gin.Context) (string, error) {
			return "whatever", nil
		})

		rawURL := "http://does.not.matter"
		w := httptest.NewRecorder()
		h(createTestContext(t, w, rawURL, map[string]string{}))
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("bad request when plaintext cannot be extracted", func(t *testing.T) {
		h := NewHmacVerifier(
			WithHmacVerifierSigKeys("some-key", "some-other-key", "k"),
		).SignedHandler(someRouteHandler, func(c *gin.Context) (string, error) {
			return "", errors.New("plaintext cannot be extracted")
		})

		rawURL := "http://does.not.matter"
		w := httptest.NewRecorder()
		h(createTestContext(t, w, rawURL, map[string]string{}))
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("override signature encoder", func(t *testing.T) {
		h := NewHmacVerifier(
			WithHmacVerifierSigKeys("some-key", "some-other-key", "k"),
			WithHmacVerifierSigEncoder(func(b []byte) string {
				return "some-static-sig"
			}),
		).SignedHandler(someRouteHandler, func(c *gin.Context) (string, error) {
			return "whatever", nil
		})

		rawURL := "http://does.not.matter"
		w := httptest.NewRecorder()
		h(createTestContext(t, w, rawURL, map[string]string{
			HmacDefaultSignatureHeader: "some-static-sig",
		}))
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func BenchmarkHmacVerifier_verifySignature(b *testing.B) {
	const msgByteSize = 512
	const numValidKeys = 100
	validKeys := make([]string, numValidKeys)
	for i := range validKeys {
		validKeys[i] = fmt.Sprintf("key-%d", i)
	}
	verifier := NewHmacVerifier(WithHmacVerifierSigKeys(validKeys...))

	for i := 0; i < b.N; i++ {
		msg := randBytes(msgByteSize)
		sig := string(randBytes(64))
		err := verifier.verifySignature(msg, sig)
		if err == nil {
			b.Errorf("expected error")
		}
	}
}

func randBytes(n int) []byte {
	var charset = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return b
}
