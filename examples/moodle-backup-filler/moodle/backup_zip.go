package moodle

import (
	"os"

	"moodle-backup-filler/logger"
)

// ZipBackupReader implements the BackupReader interface for zip formatted
// Moodle course backups.
type ZipBackupReader struct {
	file *os.File
}

// NewZipBackupReader returns a ZipBackupReader object initialised with the
// configured input file.  Returns an error if the file is not a zip file or
// the contents of the file don't look like a Moodle course backup.
func NewZipBackupReader(file *os.File) (BackupReader, error) {
	return &ZipBackupReader{
		file: file,
	}, nil
}

// Next advances to the next entry in a zip formatted Moodle backup.
func (br *ZipBackupReader) Next() (*FileHeader, error) {
	// TODO: implement this
	logger.Err.Fatal("Zip formatted Moodle backups are not currently supported")

	return nil, nil
}

// Read reads from the current file in a zip formatted Moodle backup.
func (br *ZipBackupReader) Read(b []byte) (int, error) {
	// TODO: implement this
	logger.Err.Fatal("Zip formatted Moodle backups are not currently supported")

	return 0, nil
}

// Close closes the open file.
func (br *ZipBackupReader) Close() error {
	return br.file.Close()
}

// vim: nolist expandtab ts=4 sw=4
