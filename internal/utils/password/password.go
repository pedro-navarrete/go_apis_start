// Package password proporciona funciones para hash y verificación de contraseñas
package password

import (
	"golang.org/x/crypto/bcrypt"
)

// defaultCost costo por defecto para bcrypt
const defaultCost = 12

// Hash genera un hash bcrypt de la contraseña
func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), defaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Verify verifica si una contraseña coincide con su hash
func Verify(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
