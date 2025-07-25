package handler

import (
	"net/http"

	"bucket_organizer/internal/app/server/httputils"
	"bucket_organizer/internal/app/services"
	"bucket_organizer/internal/pkg/types"
)

func UploadObject(bs *services.BucketService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		bucketId := r.PathValue("bucketId")
		objectId := r.PathValue("objectId")
		object, err := bs.InsertObject(ctx, bucketId, objectId)
		if err != nil {
			types.SetErrorInRequestContext(r, err, "error while inserting object")
			return
		}
		_ = httputils.Respond(w, r, http.StatusCreated, object)
	}
}

func GetObject(bs *services.BucketService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		bucketId := r.PathValue("bucketId")
		objectId := r.PathValue("objectId")
		object, err := bs.GetObject(ctx, bucketId, objectId)
		if err != nil {
			types.SetErrorInRequestContext(r, err, "error while getting object")
			return
		}
		_ = httputils.Respond(w, r, http.StatusOK, object)
	}
}

func DeleteObject(bs *services.BucketService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		bucketId := r.PathValue("bucketId")
		objectId := r.PathValue("objectId")
		if err := bs.RemoveObject(ctx, bucketId, objectId); err != nil {
			types.SetErrorInRequestContext(r, err, "error while deleting object")
			return
		}
		_ = httputils.Respond(w, r, http.StatusOK, "")
	}
}
