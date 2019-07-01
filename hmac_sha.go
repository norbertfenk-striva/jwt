package jwt

import (
	"crypto"
	"crypto/hmac"
	"errors"
	"hash"

	"github.com/gbrlsnchs/jwt/v3/internal"
)

var (
	// ErrHMACMissingKey is the error for trying to sign or verify a JWT with an empty key.
	ErrHMACMissingKey = errors.New("jwt: hmacSHA key is empty")
	// ErrHMACVerification is the error for an invalid signature.
	ErrHMACVerification = errors.New("jwt: hmacSHA verification failed")

	_ Algorithm = new(hmacSHA)
)

type hmacSHA struct {
	name string
	key  []byte
	sha  crypto.Hash
	size int
	pool *hashPool
}

func newHMACSHA(name string, key []byte, sha crypto.Hash) *hmacSHA {
	return &hmacSHA{
		name: name, // cache name
		key:  key,
		sha:  sha,
		size: sha.Size(), // cache size
		pool: newHashPool(func() hash.Hash { return hmac.New(sha.New, key) }),
	}
}

// NewHS256 creates a new algorithm using HMAC and SHA-256.
func NewHS256(key []byte) Algorithm {
	return newHMACSHA("HS256", key, crypto.SHA256)
}

// NewHS384 creates a new algorithm using HMAC and SHA-384.
func NewHS384(key []byte) Algorithm {
	return newHMACSHA("HS384", key, crypto.SHA384)
}

// NewHS512 creates a new algorithm using HMAC and SHA-512.
func NewHS512(key []byte) Algorithm {
	return newHMACSHA("HS512", key, crypto.SHA512)
}

// Name returns the algorithm's name.
func (hs *hmacSHA) Name() string {
	return hs.name
}

// Sign signs headerPayload using the HMAC-SHA algorithm.
func (hs *hmacSHA) Sign(headerPayload []byte) ([]byte, error) {
	if string(hs.key) == "" {
		return nil, ErrHMACMissingKey
	}
	return hs.pool.sign(headerPayload)
}

// Size returns the signature byte size.
func (hs *hmacSHA) Size() int {
	return hs.size
}

// Verify verifies a signature based on headerPayload.
func (hs *hmacSHA) Verify(headerPayload, sig []byte) (err error) {
	if sig, err = internal.DecodeToBytes(sig); err != nil {
		return err
	}
	sig2, err := hs.Sign(headerPayload)
	if err != nil {
		return err
	}
	if !hmac.Equal(sig, sig2) {
		return ErrHMACVerification
	}
	return nil
}