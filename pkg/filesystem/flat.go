package filesystem

import (
	"io"
	"os"
	"path/filepath"

	"github.com/dchest/uniuri"

	"github.com/oakmail/backend/pkg/config"
)

// Flat is a local flatfile implementation of Filesystem
type Flat struct {
	basePath string
}

// NewFlat returns a new flat filesystem
func NewFlat(cfg config.FlatConfig) (*Flat, error) {
	return &Flat{
		basePath: cfg.Path,
	}, nil
}

// Fetch loads up a file by name
func (f *Flat) Fetch(name string) (io.ReadCloser, error) {
	path := filepath.Join(f.basePath, name)

	file, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// Upload generates a new name for a file and saves it
func (f *Flat) Upload(contents io.Reader) (string, int, error) {
	name := uniuri.NewLen(uniuri.UUIDLen)
	path := filepath.Join(f.basePath, name)

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return "", 0, err
	}

	written, err := io.Copy(file, contents)
	if err != nil {
		return "", 0, err
	}

	return name, int(written), nil
}

// Delete removes the file by name
func (f *Flat) Delete(name string) error {
	path := filepath.Join(f.basePath, name)
	return os.Remove(path)
}

// Close is a nop for flatfs
// nop until streaming is implemented
func (f *Flat) Close() error {
	return nil
}
