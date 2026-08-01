package auth

import "golang.org/x/crypto/bcrypt"

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func checkPassword(hashed_password, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed_password), []byte(password))
	if err != nil {
		return false
	}
	return true
}
