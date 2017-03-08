package passwords

import (
	"crypto/rand"

	"github.com/pzduniak/go-argon2"
)

var ctx *argon2.Context

func init() {
	ctx = &argon2.Context{
		Iterations:  iterations,
		Memory:      1 << memory,
		Parallelism: parallelism,
		HashLen:     hashLen,
		Mode:        mode,
		Version:     version,
	}

	if _, err := argon2.HashEncoded(ctx, []byte("sample password"), []byte("sample salt")); err != nil {
		panic(err)
	}
}

// Hash returns a hash with a random salt.
func Hash(input string) string {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		panic(err)
	}

	hash, err := argon2.HashEncoded(ctx, []byte(input), salt)
	if err != nil {
		panic(err)
	}

	return string(hash)
}

// Verify checks if the hash and password match.
func Verify(input string, password string) bool {
	result, err := argon2.VerifyEncoded(input, []byte(password))
	if err != nil {
		panic(err)
	}

	return result
}
