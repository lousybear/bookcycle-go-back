package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(hashedPass, plainPass string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(plainPass))
}
