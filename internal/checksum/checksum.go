package checksum

import (
	"hash"
	"io"

	"github.com/cespare/xxhash/v2"
)

// Algorithm represents a checksum algorithm
type Algorithm string

const (
	// XXHash is the default fast checksum algorithm
	XXHash Algorithm = "xxhash"
)

// Calculator provides checksum calculation functionality
type Calculator interface {
	// Calculate computes checksum from a reader
	Calculate(r io.Reader) (string, error)

	// Verify checks if data matches the expected checksum
	Verify(r io.Reader, expected string) (bool, error)
}

// xxHashCalculator implements Calculator using xxHash
type xxHashCalculator struct{}

// NewCalculator creates a new checksum calculator
func NewCalculator(algo Algorithm) Calculator {
	switch algo {
	case XXHash:
		return &xxHashCalculator{}
	default:
		return &xxHashCalculator{}
	}
}

// Calculate computes xxHash checksum
func (c *xxHashCalculator) Calculate(r io.Reader) (string, error) {
	h := xxhash.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return string(h.Sum(nil)), nil
}

// Verify checks if data matches expected xxHash checksum
func (c *xxHashCalculator) Verify(r io.Reader, expected string) (bool, error) {
	actual, err := c.Calculate(r)
	if err != nil {
		return false, err
	}
	return actual == expected, nil
}

// NewHash creates a new hash instance for the algorithm
func NewHash(algo Algorithm) hash.Hash {
	switch algo {
	case XXHash:
		return xxhash.New()
	default:
		return xxhash.New()
	}
}
