package services

import (
	"context"

	"bucket_organizer/internal/app/repository/bucket"
	"bucket_organizer/internal/app/server/dto/response"
	"bucket_organizer/pkg/logger"
)

type BucketService struct {
	bucketRepo bucket.Repository
}

func NewBucketService(bucketRepo bucket.Repository) *BucketService {
	return &BucketService{
		bucketRepo: bucketRepo,
	}
}

func (s *BucketService) InsertObject(ctx context.Context, bucketId, objectId string) (*response.ObjectResponse, error) {
	if err := s.bucketRepo.InsertObject(ctx, bucketId, objectId); err != nil {
		logger.Error(ctx, "error inserting object", err)
		return nil, err
	}
	return &response.ObjectResponse{
		Id: objectId,
	}, nil
}

func (s *BucketService) GetObject(ctx context.Context, bucketId, objectId string) (*response.ObjectResponse, error) {
	o, err := s.bucketRepo.GetObject(ctx, bucketId, objectId)
	if err != nil {
		logger.Error(ctx, "error getting object", err)
		return nil, err
	}
	return &response.ObjectResponse{
		Id: o,
	}, nil
}

func (s *BucketService) RemoveObject(ctx context.Context, bucketId, objectId string) error {
	if err := s.bucketRepo.RemoveObject(ctx, bucketId, objectId); err != nil {
		logger.Error(ctx, "error removing object", err)
		return err
	}
	return nil
}
