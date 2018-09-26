package source

import (
	"io"
	"strings"

	"moodle-backup-filler/config"
)

// ContentReader implements the io.Reader interface for a single file.
type ContentReader interface {
	// Size returns the size of file being read.
	Size() int64

	io.ReadCloser
}

// GetReader returns a ContentReader of the configured type for the file
// with hash contentHash.
func GetReader(contentHash string) (ContentReader, error) {
	if strings.HasPrefix(config.Config.ContentBase, "s3://") {
		return NewS3ContentReader(contentHash)
	} else if strings.HasPrefix(config.Config.ContentBase, "http://") {
		return NewHTTPContentReader(contentHash)
	}

	return NewLocalContentReader(contentHash)
}

// vim: nolist expandtab ts=4 sw=4
