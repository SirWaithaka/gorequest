package request

import (
	"net/http"

	"github.com/SirWaithaka/gohttp"
)

type Config struct {
	// Endpoint is hostname or fully qualified URI  of the service being called
	Endpoint string

	// Set this to `true` to disable SSL when sending requests. Defaults
	// to `false`
	DisableSSL bool

	// The HTTP client to use when sending requests
	HTTPClient *http.Client

	// The logger writer interface to write logging messages to.
	Logger gohttp.Logger
}
