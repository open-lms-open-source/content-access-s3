// moodle-backup-filler
//
// Populate fileless course backups for archiving or import to a Moodle
// instance that doesn't have access to the Blackboard Open LMS Enterprise
// file store.
package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"moodle-backup-filler/config"
	"moodle-backup-filler/logger"
	"moodle-backup-filler/moodle"
)

func main() {
	if config.Config.SourceBackupFile != "" {
		// hydrate a single course backup
		hydrate(config.Config.SourceBackupFile, config.Config.DestBackupFile)
	} else {
		// hydrate a directory full of course backups
		files, err := ioutil.ReadDir(config.Config.SourceBackupDir)
		if err != nil {
			logger.Err.WithError(err).Fatalf("Unable to read directory %s", config.Config.SourceBackupDir)
		}

		for _, file := range files {
			filename := file.Name()
			if filename[0] == '.' {
				// skip hidden files
				continue
			}

			source := filepath.Join(config.Config.SourceBackupDir, filename)
			dest := filepath.Join(config.Config.DestBackupDir, filename)

			if _, err := os.Stat(dest); !os.IsNotExist(err) {
				logger.Err.Infof("Processed backup already exists for '%s', skipping", filename)
				continue
			}

			logger.Err.Infof("Processing %s", filename)
			hydrate(source, dest)
		}
	}

	os.Exit(0)
}

func hydrate(source, dest string) {
	// input file setup
	in, err := moodle.NewBackupReader(source)
	if err != nil {
		logger.Err.WithError(err).Fatal("Unable to read original backup file")
	}
	defer in.Close()

	// output file setup
	out, err := os.Create(dest)
	if err != nil {
		logger.Err.WithError(err).Fatal("Unable to write new backup file")
	}
	gzWriter := gzip.NewWriter(out)
	tarWriter := tar.NewWriter(gzWriter)

	defer func() {
		if err := tarWriter.Close(); err != nil {
			logger.Err.WithError(err).Error("Error closing tar writer")
		}
		if err := gzWriter.Close(); err != nil {
			logger.Err.WithError(err).Error("Error closing gzip writer")
		}
		if err := out.Close(); err != nil {
			logger.Err.WithError(err).Error("Error closing output file")
		}
	}()

	// process the backup
	for {
		// read from input file
		inHeader, err := in.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Err.WithError(err).Fatal("Error reading from input")
		}

		switch inHeader.Name {
		case ".ARCHIVE_INDEX":
			// The .ARCHIVE_INDEX file must be the first file in the
			// archive, but needs to be modified to reflect the entire
			// content of the backup which isn't known until the archive is
			// written.  It's optional, so rather than write the entire
			// backup to a temporary file/directory just to create this
			// file, we leave it out.  The downside of not having an
			// .ARCHIVE_INDEX file is that Moodle's list_files in slower and
			// progress reporting isn't as good.
			continue
		case "files.xml":
			// Inject files listed in files.xml from content source.
			if err := moodle.ProcessFilesXML(in, tarWriter); err != nil {
				logger.Err.WithError(err).Fatal("Failed to process files.xml")
			}
		case "moodle_backup.xml":
			// Fileless backups are marked as such in moodle_backup.xml, so
			// we change that to indicate files are included.
			if err := moodle.ProcessMoodleBackupXML(in, tarWriter); err != nil {
				logger.Err.WithError(err).Fatal("Failed to update moodle_backup.xml")
			}
		default:
			// Copy all other files from input to output.
			outHeader := &tar.Header{
				Name:     inHeader.Name,
				Size:     inHeader.Size,
				Mode:     inHeader.Mode,
				ModTime:  inHeader.ModTime,
				Typeflag: inHeader.Typeflag,
			}

			if err := tarWriter.WriteHeader(outHeader); err != nil {
				logger.Err.WithError(err).Fatal("Failed writing file header to ouput file")
			}

			if _, err := io.Copy(tarWriter, in); err != nil {
				logger.Err.WithError(err).Fatal("Failing writing file content to output file")
			}
		}
	}
}

// vim: nolist expandtab ts=4 sw=4
