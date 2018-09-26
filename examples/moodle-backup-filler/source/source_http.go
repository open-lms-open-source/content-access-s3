package source

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"moodle-backup-filler/config"
)

// httpClient is largely the same as the default http.Client, but has a
// couple of timeout tweaks (ResponseHeaderTimeout and
// ExpectContinueTimeout).
var httpClient = &http.Client{
	Timeout: time.Hour * 2,
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
	},
}

// HTTPContentReader implements the ContentReader interface for files
// accessible through HTTP.
type HTTPContentReader struct {
	reader io.ReadCloser
	size   int64
}

// NewHTTPContentReader returns a ContentReader for the given contentHash,
// which reads the file from an HTTP endpoint.
func NewHTTPContentReader(contentHash string) (*HTTPContentReader, error) {
	req, err := http.NewRequest("GET", config.Config.ContentBase+"/"+contentHash, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("Received status code '%d' while reading '%s'", resp.StatusCode, contentHash)
	}

	return &HTTPContentReader{
		reader: resp.Body,
		size:   resp.ContentLength,
	}, nil
}

// Size returns the size of the currently open file.
func (cr *HTTPContentReader) Size() int64 {
	return cr.size
}

// Read reads bytes from the currently open file.
func (cr *HTTPContentReader) Read(b []byte) (int, error) {
	return cr.reader.Read(b)
}

// Close closes the currently open file.
func (cr *HTTPContentReader) Close() error {
	return cr.reader.Close()
}

// vim: nolist expandtab ts=4 sw=4
