package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(p string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(p), 14)
	return string(bytes), err
}

func CheckPassword(hash, p string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(p))
}
