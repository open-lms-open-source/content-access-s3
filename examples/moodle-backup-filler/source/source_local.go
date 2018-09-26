package source

import (
	"os"
	"path/filepath"

	"moodle-backup-filler/config"
)

// LocalContentReader implements the ContentSource interface for files
// stored on local disk using the standard Moodle data directory layout.
type LocalContentReader struct {
	file *os.File
	size int64
}

// NewLocalContentReader returns a ContentReader for the given contentHash,
// which reads the file from local disk.
func NewLocalContentReader(contentHash string) (*LocalContentReader, error) {
	filePath := filepath.Join(config.Config.ContentBase, contentHash[:2], contentHash[2:4], contentHash)

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return &LocalContentReader{
		file: file,
		size: fileInfo.Size(),
	}, nil
}

// Size returns the size of the currently open file.
func (cr *LocalContentReader) Size() int64 {
	return cr.size
}

// Read reads bytes from the currently open file.
func (cr *LocalContentReader) Read(b []byte) (int, error) {
	return cr.file.Read(b)
}

// Close closes the currently open file.
func (cr *LocalContentReader) Close() error {
	return cr.file.Close()
}

// vim: nolist expandtab ts=4 sw=4
