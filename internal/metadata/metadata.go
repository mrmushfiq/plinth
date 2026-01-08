package metadata

import (
	"context"
	"time"
)

// ObjectState represents the state of an object
type ObjectState string

const (
	ObjectStatePending    ObjectState = "pending"
	ObjectStateCommitted  ObjectState = "committed"
	ObjectStateTombstoned ObjectState = "tombstoned"
)

// Object represents an object in the metadata store
type Object struct {
	ID             string
	BucketName     string
	ObjectKey      string
	VersionID      string
	IsLatest       bool
	IsDeleteMarker bool
	SizeBytes      int64
	ETag           string
	ContentType    string
	Placement      []string // Node IDs where replicas exist
	State          ObjectState
	Metadata       map[string]string
	Tags           map[string]string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Bucket represents a bucket in the metadata store
type Bucket struct {
	ID                string
	Name              string
	VersioningEnabled bool
	Region            string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// Service defines the interface for metadata operations
type Service interface {
	// Bucket operations
	CreateBucket(ctx context.Context, name string) (*Bucket, error)
	GetBucket(ctx context.Context, name string) (*Bucket, error)
	DeleteBucket(ctx context.Context, name string) error
	ListBuckets(ctx context.Context) ([]*Bucket, error)

	// Object operations
	CreateObject(ctx context.Context, obj *Object) error
	GetObject(ctx context.Context, bucketName, objectKey string) (*Object, error)
	GetObjectVersion(ctx context.Context, bucketName, objectKey, versionID string) (*Object, error)
	DeleteObject(ctx context.Context, bucketName, objectKey string) error
	ListObjects(ctx context.Context, bucketName, prefix string, limit int) ([]*Object, error)

	// Placement operations
	UpdateObjectPlacement(ctx context.Context, objectID string, nodeIDs []string) error

	// Repair operations
	FindUnderReplicatedObjects(ctx context.Context, replicationFactor int) ([]*Object, error)
}
