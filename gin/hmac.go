package gin

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var ErrInvalidSignature = errors.New("invalid signature")

type StrFromCtx func(c *gin.Context) (string, error)

// HmacDefaultSignatureHeader defines the default header name where clients should place the signature.
const HmacDefaultSignatureHeader = "X-REQ-SIG"

type HmacSha256Verifier struct {
	keys       [][]byte
	sigFN      StrFromCtx
	sigEncoder func([]byte) string
}

func NewHmacSha256Verifier(keys []string) *HmacSha256Verifier {
	keysB := make([][]byte, len(keys))
	for i := range keys {
		keysB[i] = []byte(keys[i])
	}

	return &HmacSha256Verifier{
		keys: keysB,
		sigFN: func(c *gin.Context) (string, error) {
			return c.GetHeader(HmacDefaultSignatureHeader), nil
		},
		sigEncoder: func(b []byte) string {
			return base64.StdEncoding.EncodeToString(b)
		},
	}
}

// WithSigFunction can be used to override signature location.
// As a query string param for example:
//
//	func(c *gin.Context) (string, error) {
//		return c.Query("sig"), nil
//	}
func (v *HmacSha256Verifier) WithSigFunction(sigFN StrFromCtx) *HmacSha256Verifier {
	v.sigFN = sigFN
	return v
}

// WithSigEncoder can be used to override default signature encoder (base64).
func (v *HmacSha256Verifier) WithSigEncoder(e func(b []byte) string) *HmacSha256Verifier {
	v.sigEncoder = e
	return v
}

// SignedHandler can be used to construct signed handlers.
// plaintextFN defines the message that should be signed by clients.
// For example:
//
//	func(c *gin.Context) (string, error) {
//		return c.Query("asset") + ":" + c.Query("coin"), nil
//	}
func (v *HmacSha256Verifier) SignedHandler(h gin.HandlerFunc, plaintextFN StrFromCtx) gin.HandlerFunc {
	return func(c *gin.Context) {
		plaintext, err := plaintextFN(c)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("cannot extract plaintext"))
			return
		}

		sig, err := v.sigFN(c)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("cannot extract signature"))
			return
		}

		if err := v.verifySignature([]byte(plaintext), sig); err != nil {
			_ = c.AbortWithError(http.StatusUnauthorized, errors.New("cannot verify signature"))
			return
		}

		h(c)
	}
}

func (v *HmacSha256Verifier) verifySignature(msg []byte, sig string) error {
	for _, signatureKey := range v.keys {
		h := hmac.New(sha256.New, signatureKey)
		h.Write(msg)
		sum := h.Sum(nil)
		encoded := v.sigEncoder(sum)
		if sig == encoded {
			return nil
		}
	}

	return ErrInvalidSignature
}
