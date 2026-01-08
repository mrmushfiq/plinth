package placement

import (
	"context"
)

// Node represents a storage node
type Node struct {
	ID       string
	Address  string
	Tier     string // hot, warm, cold
	Capacity int64
	Used     int64
	Status   string // healthy, degraded, offline
}

// Controller manages object placement across nodes
type Controller interface {
	// GetNodes returns nodes for storing an object
	GetNodes(ctx context.Context, objectKey string, replicationFactor int) ([]Node, error)

	// GetNode returns a specific node for reading
	GetNode(ctx context.Context, nodeID string) (*Node, error)

	// ListNodes returns all available nodes
	ListNodes(ctx context.Context) ([]Node, error)

	// AddNode registers a new node
	AddNode(ctx context.Context, node Node) error

	// RemoveNode removes a node from the cluster
	RemoveNode(ctx context.Context, nodeID string) error

	// UpdateNodeHealth updates node health status
	UpdateNodeHealth(ctx context.Context, nodeID string, status string, capacity, used int64) error
}
