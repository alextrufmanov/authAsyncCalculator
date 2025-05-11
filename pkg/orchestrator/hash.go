package orchestrator

import "golang.org/x/crypto/bcrypt"

func Hash(login string, password string) (string, bool) {
	// return login + password, true
	bytes, err := bcrypt.GenerateFromPassword([]byte(login+password), bcrypt.DefaultCost)
	if err != nil {
		return "", false
	}
	return string(bytes), true
}

func TestHash(login string, password string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(login+password)) == nil
}
