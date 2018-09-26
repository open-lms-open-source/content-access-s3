package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"  // TOML format config file
	"github.com/alexflint/go-arg" // command line options

	"moodle-backup-filler/logger"
	"moodle-backup-filler/version"
)

// args and Config represent active configuration from the command line and
// configuration file.
var (
	args   cliArgs
	Config TOMLConfig
)

// cliArgs defines the command line options and is used by alexflint/go-arg.
type cliArgs struct {
	Debug      bool
	ConfigFile string `arg:"--config"`

	SourceBackupFile string `arg:"--source"`
	DestBackupFile   string `arg:"--dest"`

	SourceBackupDir string `arg:"--sourcedir"`
	DestBackupDir   string `arg:"--destdir"`

	ContentBase string
}

func (cliArgs) Version() string {
	return fmt.Sprintf(version.String())
}

// Application configuration from the TOML configuration file.
type (
	// TOMLConfig represents application configuration and is populated by
	// BurntSushi/toml.
	TOMLConfig struct {
		Debug bool

		// Fileless Moodle course backup to be used as input and the output
		// file to which the hydrated backup will be written.  If not fully
		// pathed, will be prefixed with SourceBackupDir and DestBackupDir.
		SourceBackupFile string `toml:"source_backup_file"`
		DestBackupFile   string `toml:"destination_backup_file"`

		// Directory containing input files and directory to which output
		// files will be written.  If SourceBackupFile is not provided, all
		// backups in SourceBackupDir will be hydrated and written to
		// DestBackupDir.
		SourceBackupDir string `toml:"source_backup_directory"`
		DestBackupDir   string `toml:"destination_backup_directory"`

		// base URL or path for Moodle content directory
		ContentBase string `toml:"content_base"`

		// S3Region is the name of the region in which the S3 bucket exists
		S3Region string `toml:"s3_region"`

		// S3Bucket is the name of the bucket from which to read files
		S3Bucket string `toml:"-"`

		// S3AssumeRoleARN is the ARN of a role that provides read access to
		// the bucket
		S3AssumeRoleARN string `toml:"s3_assume_role_arn"`
	}
)

func init() {
	if os.Getenv("GO_TEST") != "" {
		// we're running under `go test`; ignore command line arguments.
		return
	}

	arg.MustParse(&args)

	if args.ConfigFile != "" {
		// parse configuration file
		logger.Err.Infof("Parsing configuration file '%s'", args.ConfigFile)
		if _, err := toml.DecodeFile(args.ConfigFile, &Config); err != nil {
			logger.Err.WithError(err).Fatalf("Parse failed")
		}
	}

	mutateConfig()

	if err := validateConfig(); err != nil {
		logger.Err.WithError(err).Fatalf("Config failed validation checks")
	}

	showConfig()
}

func mutateConfig() {
	// override config file with options provided on command line
	if args.Debug {
		Config.Debug = true
	}

	if args.SourceBackupFile != "" {
		Config.SourceBackupFile = args.SourceBackupFile
	}
	if args.DestBackupFile != "" {
		Config.DestBackupFile = args.DestBackupFile
	}

	if args.SourceBackupDir != "" {
		Config.SourceBackupDir = args.SourceBackupDir
	}
	if args.DestBackupDir != "" {
		Config.DestBackupDir = args.DestBackupDir
	}

	if args.ContentBase != "" {
		Config.ContentBase = args.ContentBase
	}
	if strings.HasPrefix(Config.ContentBase, "s3://") { // S3ContentReader
		if Config.S3Bucket == "" && len(Config.ContentBase) > 5 {
			Config.S3Bucket = Config.ContentBase[5:]
		}
	} else if strings.HasPrefix(Config.ContentBase, "http://") { // HTTPContentReader
		// ensure there's a trailing slash to avoid checking for slash when
		// generating URLs
		if !strings.HasSuffix(Config.ContentBase, "/") {
			Config.ContentBase = Config.ContentBase + "/"
		}
	} else { // LocalContentReader
	}

	// make SourceBackupFile absolute if possible with the given
	// configuration
	if Config.SourceBackupFile != "" && !filepath.IsAbs(Config.SourceBackupFile) {
		Config.SourceBackupFile = filepath.Join(Config.SourceBackupDir, Config.SourceBackupFile)
	}

	// make DestBackupFile absolute if possible with the given configuration
	if Config.DestBackupFile != "" && !filepath.IsAbs(Config.DestBackupFile) {
		Config.DestBackupFile = filepath.Join(Config.DestBackupDir, Config.DestBackupFile)
	}
}

func validateConfig() error {
	if Config.SourceBackupDir != "" {
		// confirm that source directory is valid
		fileInfo, err := os.Stat(Config.SourceBackupDir)
		if err != nil {
			return err
		}
		if !fileInfo.IsDir() {
			return fmt.Errorf("sourcedir '%s' is not a directory", Config.SourceBackupDir)
		}
	}
	if Config.DestBackupDir != "" {
		// confirm that destination directory is valid
		fileInfo, err := os.Stat(Config.DestBackupDir)
		if err != nil {
			return err
		}
		if !fileInfo.IsDir() {
			return fmt.Errorf("destdir '%s' is not a directory", Config.DestBackupDir)
		}
	}

	if Config.SourceBackupFile == "" && Config.SourceBackupDir == "" {
		return fmt.Errorf("At least one of source and sourcedir must be specified")
	}
	if Config.SourceBackupDir != "" && Config.SourceBackupFile == "" && Config.DestBackupDir == "" {
		return fmt.Errorf("Destination directory must be provided when filling multiple files")
	}

	if Config.SourceBackupFile != "" {
		// confirm that source file is valid
		fileInfo, err := os.Stat(Config.SourceBackupFile)
		if err != nil {
			return err
		}
		if !fileInfo.Mode().IsRegular() {
			return fmt.Errorf("source '%s' is not a file", Config.SourceBackupFile)
		}
	}

	if Config.ContentBase == "" {
		return fmt.Errorf("contentbase is required")
	}
	if strings.HasPrefix(Config.ContentBase, "s3://") {
		// S3ContentReader requires a valid S3 bucket URL
		// TODO: implement validation
	} else if strings.HasPrefix(Config.ContentBase, "http://") {
		// HTTPContentReader requires a valid HTTP URL
		// TODO: implement validation
	} else {
		// LocalContentReader requires an existing local directory
		fileInfo, err := os.Stat(Config.ContentBase)
		if err != nil {
			return err
		}
		if !fileInfo.IsDir() {
			return fmt.Errorf("contentbase '%s' is not a directory", Config.ContentBase)
		}
	}

	return nil
}

func showConfig() {
	logger.Err.Debugf("Debug: %v", Config.Debug)
	logger.Err.Debugf("SourceBackupFile: %v", Config.SourceBackupFile)
	logger.Err.Debugf("DestBackupFile: %v", Config.DestBackupFile)
	logger.Err.Debugf("SourceBackupDir: %v", Config.SourceBackupDir)
	logger.Err.Debugf("DestBackupDir: %v", Config.DestBackupDir)
	logger.Err.Debugf("ContentBase: %v", Config.ContentBase)
	logger.Err.Debugf("S3Region: %v", Config.S3Region)
	logger.Err.Debugf("S3Bucket: %v", Config.S3Bucket)
	logger.Err.Debugf("S3AssumeRoleARN: %v", Config.S3AssumeRoleARN)
}

// vim: nolist noexpandtab ts=4 sw=4
