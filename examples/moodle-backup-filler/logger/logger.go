package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Out and Err are global logger objects connected to stdout and stderr
// respectively
var (
	Out *logrus.Logger = &logrus.Logger{}
	Err *logrus.Logger = &logrus.Logger{}
)

func init() {
	Out.Out = os.Stdout
	Out.Formatter = &logrus.TextFormatter{}
	Out.Hooks = make(logrus.LevelHooks)
	Out.Level = logrus.InfoLevel

	// this function runs before parsing any command line arguments so we
	// can write logs while parsing arguments, thus we need to scan the
	// command line arguments for the debug flag to determine the correct
	// log level
	logLevel := logrus.InfoLevel
	for _, arg := range os.Args {
		if arg == "--debug" {
			logLevel = logrus.DebugLevel
		}
	}
	Err.Out = os.Stderr
	Err.Formatter = &logrus.TextFormatter{}
	Err.Hooks = make(logrus.LevelHooks)
	Err.Level = logLevel
}

// UseJSONFormat enables JSON formatted log output.
func UseJSONFormat() {
	switch Out.Formatter.(type) {
	case *logrus.JSONFormatter:
		Out.Formatter = &logrus.JSONFormatter{}
		Out.Debug("Switched error log to JSON format")
	}

	switch Err.Formatter.(type) {
	case *logrus.JSONFormatter:
		Err.Formatter = &logrus.JSONFormatter{}
		Out.Debug("Switched request log to JSON format")
	}
}
