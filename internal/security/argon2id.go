package security

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

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

type usedHashParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
	Salt        []byte
	hash        []byte
}

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
func decodeHash(encodedHash string) (params usedHashParams, err error) {
	// split the hash into it's parameters
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return usedHashParams{}, errors.New("decodeHash(): Received invalid hash")
	}

	// check if the hash is compatible with the used one
	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return usedHashParams{}, err
	}
	if version != argon2.Version {
		return usedHashParams{}, errors.New("decodeHash(): Version of argon2 not supported")
	}

	// extract the other parameters from the string
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &params.Memory, &params.Iterations, &params.Parallelism)
	if err != nil {
		return usedHashParams{}, err
	}

	// decode the salt
	params.Salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return usedHashParams{}, err
	}
	params.SaltLength = uint32(len(params.Salt))

	// decode hash
	params.hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return usedHashParams{}, err
	}
	params.KeyLength = uint32(len(params.hash))

	// done
	return params, nil
}

func ComparePasswordAndHash(password, encodedHash string) (match bool, err error) {
	// decode the given hash
	params, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// hash the password with the same parameters as the given hash
	otherHash := argon2.IDKey([]byte(password), params.Salt, params.Iterations, params.Memory, params.Parallelism, params.KeyLength)

	// securely compare both hashes
	if subtle.ConstantTimeCompare(params.hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}
