package moodle

import (
	"archive/tar"
	"compress/gzip"
	"os"
)

// TgzBackupReader implements the BackupReader interface for tar.gz
// formatted Moodle course backups.
type TgzBackupReader struct {
	reader *tar.Reader
	file   *os.File
}

// NewTgzBackupReader returns a TgzBackupReader object initialised with the
// configured input file.  Returns an error if the file is not a gzipped tar
// file or the contents of the file don't look like a Moodle course backup.
func NewTgzBackupReader(file *os.File) (BackupReader, error) {
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	tarReader := tar.NewReader(gzipReader)

	return &TgzBackupReader{
		file:   file,
		reader: tarReader,
	}, nil
}

// Next advances to the next entry in a tar formatted Moodle backup.
func (br *TgzBackupReader) Next() (header *FileHeader, err error) {
	tarHeader, tarErr := br.reader.Next()
	if tarErr == nil {
		header = &FileHeader{
			Name:     tarHeader.Name,
			Size:     tarHeader.Size,
			Mode:     tarHeader.Mode,
			ModTime:  tarHeader.ModTime,
			Typeflag: tarHeader.Typeflag,
		}
	}
	err = tarErr

	return
}

// Read reads from the current file in a tar formatted Moodle backup.
func (br *TgzBackupReader) Read(b []byte) (int, error) {
	return br.reader.Read(b)
}

// Close closes the open file.
func (br *TgzBackupReader) Close() error {
	return br.file.Close()
}

// vim: nolist expandtab ts=4 sw=4
