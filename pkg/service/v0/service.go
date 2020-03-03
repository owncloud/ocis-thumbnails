package svc

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/owncloud/ocis-thumbnails/pkg/config"
	"github.com/owncloud/ocis-thumbnails/pkg/thumbnails"
	"github.com/owncloud/ocis-thumbnails/pkg/thumbnails/cache"
)

// Service defines the extension handlers.
type Service interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	Thumbnails(http.ResponseWriter, *http.Request)
}

// NewService returns a service implementation for Service.
func NewService(opts ...Option) Service {
	options := newOptions(opts...)

	m := chi.NewMux()
	m.Use(options.Middleware...)

	svc := Thumbnails{
		config: options.Config,
		mux:    m,
		manager: thumbnails.SimpleManager{
			Cache: cache.NewInMemoryCache(),
		},
	}

	m.Route(options.Config.HTTP.Root, func(r chi.Router) {
		r.Get("/thumbnails", svc.Thumbnails)
	})

	return svc
}

// Thumbnails defines implements the business logic for Service.
type Thumbnails struct {
	config  *config.Config
	mux     *chi.Mux
	manager thumbnails.Manager
}

// ServeHTTP implements the Service interface.
func (g Thumbnails) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.mux.ServeHTTP(w, r)
}

// Thumbnails provides the endpoint to retrieve a thumbnail for an image
func (g Thumbnails) Thumbnails(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	width, _ := strconv.Atoi(query.Get("w"))
	height, _ := strconv.Atoi(query.Get("h"))
	fileType := query.Get("type")
	fileID := query.Get("file_id")

	encoder := thumbnails.EncoderForType(fileType)
	if encoder == nil {
		// TODO: better error responses
		w.Write([]byte("can't encode that"))
		return
	}
	ctx := thumbnails.ThumbnailContext{
		Width:     width,
		Height:    height,
		ImagePath: fileID,
		Encoder:   encoder,
	}
	thumbnail, err := g.manager.Get(ctx)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	w.Write(thumbnail)
}
