package users

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "password"

	response := hashPassword(password)

	if response == password {
		t.Log("password not properly hashed")
		t.FailNow()
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "password"

	hash := hashPassword(password)

	response := checkPasswordHash(password, hash)

	if response != true {
		t.Log("password check failed")
		t.FailNow()
	}
}
