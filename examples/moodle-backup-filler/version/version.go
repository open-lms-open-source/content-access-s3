package version

import (
	"encoding/json"
	"fmt"
)

// Version and BuildDate contain build information populated by the
// compiler.
var (
	Version   string
	BuildDate string
)

type versionInfo struct {
	Version   string `json:"version"`
	BuildDate string `json:"build_date"`
}

// String renders version information as a string.
func String() string {
	return fmt.Sprintf("%s, built %s", Version, BuildDate)
}

// JSON renders version information in JSON format.
func JSON() string {
	v, _ := json.Marshal(versionInfo{Version, BuildDate})
	return string(v)
}
