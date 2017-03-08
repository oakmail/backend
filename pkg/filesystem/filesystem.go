package filesystem

import (
	"io"
)

// Filesystem is an interface for anything that allows storing files remotely
type Filesystem interface {
	Fetch(name string) (file io.ReadCloser, err error)
	Upload(contents io.Reader) (fid string, written int, err error)
	Delete(name string) error
	// consider streaming
	Close() error
}
