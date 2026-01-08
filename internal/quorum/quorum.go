package quorum

import (
	"context"
	"errors"
)

var (
	// ErrInsufficientNodes is returned when not enough nodes are available
	ErrInsufficientNodes = errors.New("insufficient nodes available for quorum")

	// ErrWriteQuorumNotMet is returned when write quorum is not satisfied
	ErrWriteQuorumNotMet = errors.New("write quorum not met")

	// ErrReadQuorumNotMet is returned when read quorum is not satisfied
	ErrReadQuorumNotMet = errors.New("read quorum not met")
)

// Config defines quorum configuration
type Config struct {
	ReplicationFactor int // Total number of replicas
	WriteQuorum       int // Minimum successful writes
	ReadQuorum        int // Minimum successful reads
}

// Result represents the result of an operation
type Result struct {
	NodeID  string
	Success bool
	Error   error
	Data    interface{}
}

// Writer handles quorum writes
type Writer interface {
	// Write performs a quorum write operation
	Write(ctx context.Context, nodeIDs []string, data []byte) ([]Result, error)
}

// Reader handles quorum reads
type Reader interface {
	// Read performs a quorum read operation
	Read(ctx context.Context, nodeIDs []string, objectKey string) ([]byte, []Result, error)
}

// ValidateConfig validates quorum configuration
func ValidateConfig(cfg Config) error {
	if cfg.ReplicationFactor < 1 {
		return errors.New("replication factor must be at least 1")
	}
	if cfg.WriteQuorum < 1 || cfg.WriteQuorum > cfg.ReplicationFactor {
		return errors.New("write quorum must be between 1 and replication factor")
	}
	if cfg.ReadQuorum < 1 || cfg.ReadQuorum > cfg.ReplicationFactor {
		return errors.New("read quorum must be between 1 and replication factor")
	}
	return nil
}
