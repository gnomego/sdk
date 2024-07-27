package crypto_test

import (
	"testing"

	"github.com/gnomeco/crypto"
	assert2 "github.com/stretchr/testify/assert"
)

func TestAes256CBC(t *testing.T) {
	cipher := crypto.NewAes256CBC()

	plaintext := []byte("Hello, World!")
	key := []byte("0123456789abcdef0123456789abcdef")
	encrypted, err := cipher.Encrypt(key, plaintext)

	if err != nil {
		t.Fatalf("Failed to encrypt plaintext: %v", err)
	}

	decrypted, err := cipher.Decrypt(key, encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt plaintext: %v", err)
	}

	assert := assert2.New(t)
	assert.Equal(plaintext, decrypted)
}
