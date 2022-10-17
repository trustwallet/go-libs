package crypto

import (
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// Signer is an interface that can be used to sign messages
type Signer interface {
	Sign(msg []byte) ([]byte, error)
}

// SignFunc is a wrapper of sign functions to implement Signer interface
type SignFunc func(msg []byte) ([]byte, error)

func (sf SignFunc) Sign(msg []byte) ([]byte, error) {
	return sf(msg)
}

// NewSHA256WithRSASigner returns Signer instance which signs msg using SHA256WithRSA and private key
func NewSHA256WithRSASigner(privateKey *rsa.PrivateKey) SignFunc {
	return func(msg []byte) ([]byte, error) {
		return SHA256WithRSA(msg, privateKey)
	}
}

// NewHMACSHA256Signer returns Signer instance which signs msg using HMACSHA256S and key
func NewHMACSHA256Signer(key string) SignFunc {
	return func(msg []byte) ([]byte, error) {
		return HMACSHA256(msg, key)
	}
}

// SHA256WithRSA signs SHA256 hash of the message with RSA privateKey
func SHA256WithRSA(msg []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	if privateKey == nil {
		return nil, errors.New("private key is empty")
	}

	h := sha256.New()
	_, err := h.Write(msg)
	if err != nil {
		return nil, fmt.Errorf("write bytes: %v", err)
	}

	res, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, h.Sum(nil))
	if err != nil {
		return nil, fmt.Errorf("SignPKCS1v15: %v", err)
	}

	return res, nil
}

// HMACSHA256 signs message with HMAC SHA256 using key
func HMACSHA256(msg []byte, key string) ([]byte, error) {
	h := hmac.New(sha256.New, []byte(key))
	_, err := h.Write(msg)
	if err != nil {
		return nil, fmt.Errorf("hmac write: %v", err)
	}
	return h.Sum(nil), nil
}

// GetRSAPrivateKey reads RSA private key from the reader
func GetRSAPrivateKey(reader io.Reader) (*rsa.PrivateKey, error) {
	bs, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("read bytes: %v", err)
	}

	privPem, _ := pem.Decode(bs)
	if privPem == nil {
		return nil, errors.New("decoded key is empty")
	}

	if privPem.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("key type is not RSA private key")
	}

	var parsedKey interface{}
	parsedKey, err = x509.ParsePKCS1PrivateKey(privPem.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse PKCS1 private key: %v", err)
	}

	privateKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("parsed key is not RSA private key")
	}

	return privateKey, nil
}

// GetRSAPrivateKeyFromFile reads RSA private key from file
func GetRSAPrivateKeyFromFile(fileName string) (*rsa.PrivateKey, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("open file: %v", err)
	}

	return GetRSAPrivateKey(file)
}

// GetRSAPrivateKeyFromString reads RSA private key from string
func GetRSAPrivateKeyFromString(s string) (*rsa.PrivateKey, error) {
	return GetRSAPrivateKey(strings.NewReader(s))
}
