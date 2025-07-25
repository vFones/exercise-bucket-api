package server

import (
	"net/http"
	"net/http/pprof"

	"bucket_organizer/internal/app/server/handler"
	"bucket_organizer/internal/app/server/middleware"
)

func middlewares(handler http.Handler) http.Handler {
	return middleware.Logging(middleware.ErrorResponder(handler))
}

func (s *Server) setupRoutes() {
	s.router.Handle("PUT /objects/{bucketId}/{objectId}", middlewares(handler.UploadObject(s.services.BucketService)))
	s.router.Handle("GET /objects/{bucketId}/{objectId}", middlewares(handler.GetObject(s.services.BucketService)))
	s.router.Handle("DELETE /objects/{bucketId}/{objectId}", middlewares(handler.DeleteObject(s.services.BucketService)))

	s.router.Handle("GET /debug/pprof/", middlewares(http.HandlerFunc(pprof.Index)))
	s.router.Handle("GET /debug/pprof/cmdline", middlewares(http.HandlerFunc(pprof.Cmdline)))
	s.router.Handle("GET /debug/pprof/profile", middlewares(http.HandlerFunc(pprof.Profile)))
	s.router.Handle("GET /debug/pprof/symbol", middlewares(http.HandlerFunc(pprof.Symbol)))
	s.router.Handle("GET /debug/pprof/trace", middlewares(http.HandlerFunc(pprof.Trace)))
	s.router.Handle("GET /debug/pprof/{cmd}", middlewares(http.HandlerFunc(pprof.Index)))
}
