package webauthnif

import (
	"math/rand"
	"time"
)

// ToBufferSource translates given string to BufferSource
func ToBufferSource(str string) (bs BufferSource) {
	//bytes := sha256.Sum256([]byte(str))
	//bs = bytes[:]
	//return
	bs = []byte(str)
	return
}

func GenChallenge() (challenge BufferSource, err error) {
	challenge = make(BufferSource, 32)
	rand.Seed(time.Now().UnixNano())
	_, err = rand.Read(challenge)
	return
}
