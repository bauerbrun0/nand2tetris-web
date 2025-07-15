package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

type PasswordHasher struct{}

type PasswordHashParams struct {
	memory      uint32 // amount of memory used by the algorithm in kibibytes
	iterations  uint32 // number of iterations over the memory
	parallelism uint8  // number of threads used by the algorithm
	saltLength  uint32 // length of the random salt
	keyLength   uint32 // Length of the generated password hash
}

var DefaultPasswordHashParams = PasswordHashParams{
	memory:      64 * 1024,
	iterations:  3,
	parallelism: 2,
	saltLength:  16,
	keyLength:   32,
}

func (p *PasswordHasher) GenerateFromPassword(password string, hp PasswordHashParams) (string, error) {
	salt, err := p.generateRandomBytes(hp.saltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, hp.iterations, hp.memory, hp.parallelism, hp.keyLength)

	// base64 encode the salt and hash so it can be stored
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// encoded representation of the hashed password
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, hp.memory, hp.iterations, hp.parallelism, b64Salt, b64Hash)
	return encodedHash, nil
}

func (p *PasswordHasher) ComparePasswordAndHash(password, encodedHash string) (bool, error) {
	hp, salt, hash, err := p.decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey([]byte(password), salt, hp.iterations, hp.memory, hp.parallelism, hp.keyLength)

	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func (p *PasswordHasher) generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (p *PasswordHasher) decodeHash(encodedHash string) (hp *PasswordHashParams, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	hp = &PasswordHashParams{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &hp.memory, &hp.iterations, &hp.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	hp.saltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	hp.keyLength = uint32(len(hash))

	return hp, salt, hash, nil
}
