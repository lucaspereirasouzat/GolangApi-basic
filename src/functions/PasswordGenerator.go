package functions

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// GeneratePassword faz o hash da senha
func GeneratePassword(passwordSend string) (string, error) {
	password := []byte(passwordSend)

	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ComparePasswords Faz a comparação da senha
func ComparePasswords(passwordSend string, hashedPassword []byte) bool {
	password := []byte(passwordSend)

	// Comparing the password with the hash
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	fmt.Println(err) // nil means it is a match
	if err != nil {
		return false
	} else {
		return true
	}
}
