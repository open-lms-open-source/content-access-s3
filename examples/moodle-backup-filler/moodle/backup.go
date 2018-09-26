package moodle

import (
	"archive/tar"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/beevik/etree" // for working with XML files
)

// FileHeader represents metadata for files in a Moodle course backup.
type FileHeader struct {
	Name     string
	Size     int64
	Mode     int64
	ModTime  time.Time
	Typeflag byte
}

// BackupReader represents a Moodle course backup on disk.
type BackupReader interface {
	// Next advances to the next entry in the Moodle backup.  The
	// FileHeader.Size determines how many bytes can be read for the next
	// file.  Any remaining data in the current file is automatically
	// discarded.
	//
	// io.EOF is returned at the end of the input.
	Next() (*FileHeader, error)

	// Read reads from the current file in the Moodle backup.  It returns
	// (0, io.EOF) when it reaches the end of that file, until Next is
	// called to advance to the next file.
	Read(b []byte) (int, error)

	// Close closes the file being read.
	Close() error
}

// NewBackupReader creates and populates a new BackupReader object of the
// appropriate type for the file passed as input.
func NewBackupReader(filename string) (BackupReader, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	// determine file type for input file
	header := make([]byte, 512)
	n, err := file.Read(header)
	if err != nil && err != io.EOF {
		return nil, err
	}
	file.Seek(0, 0)

	fileType := http.DetectContentType(header[:n])

	// create an appropriate BackupFile object
	switch fileType {
	case "application/x-gzip":
		return NewTgzBackupReader(file)
	case "application/zip":
		return NewZipBackupReader(file)
	}

	return nil, fmt.Errorf("Unsupported input file type %s", fileType)
}

// ProcessMoodleBackupXML copies the moodle_backup.xml file from in to out,
// while adjusting it to indicate that the backup archive contains files (as
// opposed to being a "fileless" backup).
func ProcessMoodleBackupXML(in io.Reader, out *tar.Writer) error {
	doc := etree.NewDocument()
	_, err := doc.ReadFrom(in)
	if err != nil {
		return err
	}

	element := doc.Root().FindElement("/moodle_backup/information/settings/setting[name='course_files']/value")
	if element == nil {
		// element doesn't exist, create it
		// <setting>
		//   <level>root</level>
		//   <name>course_files</name>
		//   <value>1</value>
		// </setting>
		element = etree.NewElement("setting")
		element.CreateAttr("level", "root")
		element.CreateAttr("name", "course_files")
		element.CreateAttr("value", "1")
		settings := doc.Root().FindElement("/moodle_backup/information/settings")
		if settings == nil {
			return fmt.Errorf("moodle_backup.xml is not valid")
		}
		settings.AddChild(element)
	} else {
		// element exists, update it
		element.SetText("1")
	}

	// generate the updated XML document
	outBytes, err := doc.WriteToBytes()
	if err != nil {
		return fmt.Errorf("moodle_backup.xml could not be regenerated")
	}

	// write the file header and updated moodle_backup.xml to the tar file
	outHeader := &tar.Header{
		Name:     "moodle_backup.xml",
		Size:     int64(len(outBytes)),
		Mode:     0644,
		ModTime:  time.Now(),
		Typeflag: tar.TypeReg,
	}

	if err := out.WriteHeader(outHeader); err != nil {
		return fmt.Errorf("Failed writing moodle_backup.xml header to output file: %v", err)
	}

	if _, err := out.Write(outBytes); err != nil {
		return fmt.Errorf("Failing writing moodle_backup.xml to output file: %v", err)
	}

	return nil
}

// vim: nolist expandtab ts=4 sw=4
