package dependency

import (
	"context"

	"bucket_organizer/internal/app/repository/bucket"
	"bucket_organizer/internal/app/server"
	"bucket_organizer/internal/app/services"
	"bucket_organizer/pkg/logger"
)

func Inject(ctx context.Context) (*server.Server, error) {
	logger.Debug(ctx, "injecting dependencies")

	bucketRepository := bucket.NewInMemoryRepo()
	bucketService := services.NewBucketService(bucketRepository)

	appServices := services.NewServices(bucketService)

	return server.NewServer(appServices, nil)
}
