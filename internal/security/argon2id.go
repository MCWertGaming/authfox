package security

import (
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

// These constants define the hash strength
const (
	// currently 64 * 1024
	hashMemory      = 65536
	hashIterations  = 1
	hashParallelism = 2
	hashSaltLength  = 16
	hashKeyLength   = 32
)

func CreateHash(password string) (string, error) {
	// create a random salt
	salt, err := randomBytes(hashSaltLength)
	if err != nil {
		return "", err
	}
	// create encryption key
	key := argon2.IDKey([]byte(password), []byte(salt), hashIterations, hashMemory, hashParallelism, hashKeyLength)

	// encode values as string
	saltEncoded := base64.RawStdEncoding.EncodeToString(salt)
	keyEncoded := base64.RawStdEncoding.EncodeToString(key)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, hashMemory, hashIterations, hashParallelism, saltEncoded, keyEncoded), nil
}
