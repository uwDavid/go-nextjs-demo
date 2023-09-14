package utils

import "golang.org/x/crypto/bcrypt"

func GenerateHashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func VerifyPassword(password, HashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(HashedPassword), []byte(password))
	return err
}
