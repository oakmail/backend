package passwords

import (
	"github.com/pzduniak/go-argon2"
)

const (
	iterations  = 10
	memory      = 15
	parallelism = 8
	hashLen     = 32
	mode        = argon2.ModeArgon2i
	version     = argon2.Version13
)
