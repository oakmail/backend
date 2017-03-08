package filesystem

import (
	"io"
	"net/http"

	"github.com/dchest/uniuri"
	"github.com/ginuerzh/weedo"

	"github.com/oakmail/backend/pkg/config"
)

// Seaweed is an implementation of SeaweedFS connection
type Seaweed struct {
	Client *weedo.Client
}

// NewSeaweed creates a new SeaweedFS client
func NewSeaweed(cfg config.SeaweedConfig) (*Seaweed, error) {
	s := &Seaweed{
		Client: weedo.NewClient(cfg.MasterURL),
	}

	return s, nil
}

// Fetch downloads a file from the storage
func (s *Seaweed) Fetch(name string) (io.ReadCloser, error) {
	public, _, err := s.Client.GetUrl(name)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(public)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

// Upload sends a file to the storage
func (s *Seaweed) Upload(contents io.Reader) (string, int, error) {
	fid, length, err := s.Client.AssignUpload(
		uniuri.NewLen(uniuri.UUIDLen),
		"application/octet-stream",
		contents,
	)
	return fid, int(length), err
}

// Delete removes a file from the server
func (s *Seaweed) Delete(name string) error {
	return s.Client.Delete(name, 1)
}

// Close is a nop for seaweedfs
func (s *Seaweed) Close() error {
	return nil
}
