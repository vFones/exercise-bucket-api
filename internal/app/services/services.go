package services

type Services struct {
	BucketService *BucketService
}

func NewServices(bs *BucketService) *Services {
	return &Services{
		BucketService: bs,
	}
}
