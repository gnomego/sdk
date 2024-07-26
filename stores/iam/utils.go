package iam

import "golang.org/x/crypto/bcrypt"

func HashSecret(secret string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(secret), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ValidateSecret(secret string, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(secret))
	return err
}
