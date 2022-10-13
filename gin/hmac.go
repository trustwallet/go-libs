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

type HmacVerifier struct {
	keys       [][]byte
	sigFN      StrFromCtx
	sigEncoder func([]byte) string
}

func NewHmacVerifier(options ...func(verifier *HmacVerifier)) *HmacVerifier {
	verifier := &HmacVerifier{
		sigFN: func(c *gin.Context) (string, error) {
			return c.GetHeader(HmacDefaultSignatureHeader), nil
		},
		sigEncoder: func(b []byte) string {
			return base64.StdEncoding.EncodeToString(b)
		},
	}

	for _, o := range options {
		o(verifier)
	}
	return verifier
}

// WithHmacVerifierSigKeys is used to set the valid signature keys.
func WithHmacVerifierSigKeys(keys ...string) func(*HmacVerifier) {
	return func(v *HmacVerifier) {
		keysB := make([][]byte, len(keys))
		for i := range keys {
			keysB[i] = []byte(keys[i])
		}
		v.keys = keysB
	}
}

// WithHmacVerifierSigFunction can be used to override signature location.
// As a query string param for example:
//
//	sigFn := func(c *gin.Context) (string, error) {
//		return c.Query("sig"), nil
//	}
func WithHmacVerifierSigFunction(sigFN StrFromCtx) func(*HmacVerifier) {
	return func(v *HmacVerifier) {
		v.sigFN = sigFN
	}
}

// WithHmacVerifierSigEncoder can be used to override default signature encoder (base64).
func WithHmacVerifierSigEncoder(e func(b []byte) string) func(*HmacVerifier) {
	return func(v *HmacVerifier) {
		v.sigEncoder = e
	}
}

// SignedHandler can be used to construct signed handlers.
// plaintextFN defines the message that should be signed by clients.
// For example:
//
//	func(c *gin.Context) (string, error) {
//		return c.Query("asset") + ":" + c.Query("coin"), nil
//	}
func (v *HmacVerifier) SignedHandler(h gin.HandlerFunc, plaintextFN StrFromCtx) gin.HandlerFunc {
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

func (v *HmacVerifier) verifySignature(msg []byte, sig string) error {
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
