package core

import (
	"math/rand"
	"time"
)

const PasswordLength = 256

type Password [PasswordLength]byte

func init() {
	rand.Seed(time.Now().Unix())
}

func getKey() *Password {
	intArr := rand.Perm(PasswordLength)
	password := &Password{}
	for i, v := range intArr {
		password[i] = byte(v)
		if i == v {
			return getKey()
		}
	}
	return password
}
