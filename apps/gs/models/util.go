package models

import "golang.org/x/crypto/bcrypt"

var defaultBcryptHasher *BcryptSecretHasher
var defaultSecretHasher SecretHasher

type SecretHasher interface {
	HashSecretString(secret string) (string, error)
	HashSecretBytes(secret []byte) ([]byte, error)
	ValidateSecretString(secret string, hash string) error
	ValidateSecretBytes(secret []byte, hash []byte) error
}

type BcryptSecretHasher struct {
	Cost int
}

func DefaultBcryptSecretHasher() *BcryptSecretHasher {
	if defaultBcryptHasher == nil {
		defaultBcryptHasher = NewBcryptSecretHasher(14)
	}

	return defaultBcryptHasher
}

func DefaultSecretHasher() SecretHasher {
	if defaultSecretHasher == nil {
		defaultSecretHasher = DefaultBcryptSecretHasher()
	}

	return defaultSecretHasher
}

func NewBcryptSecretHasher(cost int) *BcryptSecretHasher {
	return &BcryptSecretHasher{
		Cost: cost,
	}
}

func (h *BcryptSecretHasher) HashSecretString(secret string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(secret), h.Cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (h *BcryptSecretHasher) HashSecretBytes(secret []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(secret, h.Cost)
}

func (h *BcryptSecretHasher) ValidateSecretString(secret string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(secret))
}

func (h *BcryptSecretHasher) ValidateSecretBytes(secret []byte, hash []byte) error {
	return bcrypt.CompareHashAndPassword(hash, secret)
}
