package bucket

import (
	"context"

	"bucket_organizer/internal/pkg/types"
	"bucket_organizer/pkg/logger"
)

type InMemoryRepo struct {
	cache map[string]map[string]bool
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		cache: make(map[string]map[string]bool),
	}
}

func (r *InMemoryRepo) InsertObject(ctx context.Context, bucketId, objectId string) error {
	if r.cache[bucketId] == nil {
		r.cache[bucketId] = make(map[string]bool)
	}
	r.cache[bucketId][objectId] = true
	return nil
}

func (r *InMemoryRepo) GetObject(ctx context.Context, bucketId, objectId string) (string, error) {
	l, ok := r.cache[bucketId]
	if !ok {
		logger.Error(ctx, "bucket not found")
		return "", types.ErrNoBucketFound
	}
	if _, ok := l[objectId]; ok {
		logger.Debug(ctx, "found object", logger.NewLogValue("object", objectId))
		return objectId, nil
	}
	logger.Error(ctx, "object not found")
	return "", types.ErrNoObjectFound
}

func (r *InMemoryRepo) RemoveObject(ctx context.Context, bucketId, objectId string) error {
	l, ok := r.cache[bucketId]
	if !ok {
		logger.Error(ctx, "bucket not found")
		return types.ErrNoBucketFound
	}
	if _, ok := l[objectId]; ok {
		logger.Debug(ctx, "found object, deleting...", logger.NewLogValue("object", objectId))
		delete(l, objectId)
		return nil
	}
	logger.Error(ctx, "object not found")
	return types.ErrNoObjectFound
}
