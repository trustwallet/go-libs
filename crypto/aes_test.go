package crypto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAES(t *testing.T) {
	secret, message := []byte("RfUjXnZr4u7x!A%D*G-KaPdSgVkYp3s5"), "a plain text"

	_, err := AESEncrypt([]byte("not_a_valid_length"), message)
	assert.Error(t, err)

	encryptedMessage, err := AESEncrypt(secret, message)
	assert.NoError(t, err)
	assert.NotEmpty(t, encryptedMessage)
	assert.NotEqual(t, encryptedMessage, message)

	decryptedMessage, err := AESDecrypt(secret, encryptedMessage)
	assert.NoError(t, err)
	assert.Equal(t, decryptedMessage, message)

	failedMessage, err := AESDecrypt([]byte("this_is_an_invalid_secret_key___"), decryptedMessage)
	assert.Error(t, err)
	assert.Empty(t, failedMessage)
}
