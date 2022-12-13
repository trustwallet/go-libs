package crypto

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIBOgIBAAJBAK1ASa283Iotdl+Sbp5IRNjumvuTs/r0ZSt1S/8dqe08WN2GiDXn
f+U1UOJPDp5qN7d+AoQSMUg2bHXeLjrxxCUCAwEAAQJAcYfJQGKcmqfEBEju2CY/
h3CEewuFS5RPn7TTwi/sJJrtEkeha4CYgGJJusAr8K3J0O8EBnMtEz+KltYDWd6i
AQIhANSWLwXtb0lUqemqoslj3RKirsHac30IyyiJ45NQWp5BAiEA0KGuouUQdNbL
vso31iilbUnJJ54k1C8hREoEAqx9NOUCIQC5INByaQKw6XnOczqwBrdOsz1cs9A+
4pmJBAubDi7cAQIgOIFx4SCVQm/iovv1/4TmuSDg4GAOrYFOS0aYq3i4OJkCIAQw
PklhQYvKRwjm1jiktUyTyRHIDSVSmveZ/8N6zJSW
-----END RSA PRIVATE KEY-----
`

func testCompareKeys(t *testing.T, exp string, act *rsa.PrivateKey) {
	var buf bytes.Buffer
	assert.NoError(t, pem.Encode(&buf, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(act)}))

	assert.Equal(t, exp, buf.String())
}

func TestGetRSAPrivateKey(t *testing.T) {
	key, err := GetRSAPrivateKey(strings.NewReader(privateKey))
	assert.NoError(t, err)

	testCompareKeys(t, privateKey, key)
}

func TestHMACSHA256(t *testing.T) {
	res, err := HMACSHA256([]byte("test"), "e9a9b09e-6dfb-455e-8c27-7b206bec08a1")
	assert.NoError(t, err)
	assert.Equal(
		t,
		"9e99537c0a09c501bb348bc12743707beee35eba0b1bd885de15f91bc9311047",
		fmt.Sprintf("%x", string(res)),
	)
}

func TestSHA256WithRSA(t *testing.T) {
	key, err := GetRSAPrivateKey(strings.NewReader(privateKey))
	assert.NoError(t, err)

	res, err := SHA256WithRSA([]byte("test"), key)
	assert.NoError(t, err)
	assert.Equal(
		t,
		"dUiTbTPRbMhL0GyTuAE+BAbSxfEwdbWdzQuF2r3esVKg0CMtEa2btCN7O0eQezQFDRIQVXmhKRccqWPQw/Zjbw==",
		base64.StdEncoding.EncodeToString(res),
	)
}

func TestGetRSAPrivateKeyFromFile(t *testing.T) {
	f, err := os.CreateTemp("", "test")
	assert.NoError(t, err)
	defer os.Remove(f.Name())

	_, err = f.Write([]byte(privateKey))
	assert.NoError(t, err)
	assert.NoError(t, f.Sync())

	key, err := GetRSAPrivateKeyFromFile(f.Name())
	assert.NoError(t, err)

	testCompareKeys(t, privateKey, key)
}

func TestGetRSAPrivateKeyFromString(t *testing.T) {
	key, err := GetRSAPrivateKeyFromString(privateKey)
	assert.NoError(t, err)

	testCompareKeys(t, privateKey, key)
}
