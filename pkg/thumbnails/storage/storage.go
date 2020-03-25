package storage

import "github.com/owncloud/ocis-thumbnails/pkg/thumbnails/resolutions"

// Request combines different attributes needed for storage operations.
type Request struct {
	ETag       string
	Types      []string
	Resolution resolutions.Resolution
}

// Storage defines the interface for a thumbnail store.
type Storage interface {
	Get(string) []byte
	Set(string, []byte) error
	BuildKey(Request) string
}
