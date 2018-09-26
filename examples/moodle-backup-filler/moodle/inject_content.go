package moodle

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/beevik/etree"

	"moodle-backup-filler/logger"
	"moodle-backup-filler/source"
)

func injectFile(contentHash string, out *tar.Writer) error {
	var reader io.Reader
	var size int64

	if contentHash == "da39a3ee5e6b4b0d3255bfef95601890afd80709" {
		// special handling for the empty file
		reader = &bytes.Buffer{}
		size = int64(0)
	} else {
		var err error
		reader, err = source.GetReader(contentHash)
		if err != nil {
			// Moodle handles restoring backups with missing files
			// relatively gracefully, so warn but don't quit
			logger.Err.WithError(err).Warnf("Unable to read file '%s', skipping", contentHash)
			return nil
		}
		size = reader.(source.ContentReader).Size()
	}

	header := &tar.Header{
		Name:     "files/" + contentHash[:2] + "/" + contentHash,
		Size:     size,
		Mode:     0644,
		ModTime:  time.Now(),
		Typeflag: tar.TypeReg,
	}

	if err := out.WriteHeader(header); err != nil {
		return fmt.Errorf("Failed writing header to output for %s: %v", contentHash, err)
	}

	if _, err := io.Copy(out, reader); err != nil {
		return fmt.Errorf("Failed writing file to output for %s: %v", contentHash, err)
	}

	return nil
}

// ProcessFilesXML reads files.xml from in, adds all files it mentions to
// out, then writes the original files.xml to out as well.
//
// This is essentially the purpose of this software.
func ProcessFilesXML(in io.Reader, out *tar.Writer) error {
	doc := etree.NewDocument()
	_, err := doc.ReadFrom(in)
	if err != nil {
		return err
	}

	filesElement := doc.Root().FindElement("/files")
	if filesElement == nil {
		return fmt.Errorf("files.xml in input is invalid, aborting")
	}

	// need to keep track of files already injected so we can deduplicate
	filesAdded := map[string]bool{}

	for _, fileElement := range filesElement.ChildElements() {
		contentHashElement := fileElement.SelectElement("contenthash")
		contentHash := contentHashElement.Text()

		_, exists := filesAdded[contentHash]
		if !exists {
			if err := injectFile(contentHash, out); err != nil {
				return err
			}

			filesAdded[contentHash] = true
		}
	}

	// prepare for writing to tarfile
	outBytes, err := doc.WriteToBytes()
	if err != nil {
		return fmt.Errorf("files.xml could not be written to bytes buffer")
	}

	// write files.xml header and file to tarfile
	outHeader := &tar.Header{
		Name:     "files.xml",
		Size:     int64(len(outBytes)),
		Mode:     0644,
		ModTime:  time.Now(),
		Typeflag: tar.TypeReg,
	}

	if err := out.WriteHeader(outHeader); err != nil {
		return fmt.Errorf("Failed writing files.xml header to output file: %v", err)
	}

	if _, err := out.Write(outBytes); err != nil {
		return fmt.Errorf("Failed writing files.xml to output file: %v", err)
	}

	return nil
}

// vim: nolist expandtab ts=4 sw=4
