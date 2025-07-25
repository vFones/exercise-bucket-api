package bucket

import "context"

type Repository interface {
	InsertObject(ctx context.Context, bucketId, object string) error
	GetObject(ctx context.Context, bucketId, object string) (string, error)
	RemoveObject(ctx context.Context, bucketId, object string) error
}
