package slim

import (
	"net/http"
	"strings"

	"k8s.io/client-go/rest"
)

// Client is a lightweight Kubernetes client that deserializes API
// responses directly into user-provided Go structs.
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewForConfig creates a new Client from a *rest.Config.
func NewForConfig(config *rest.Config) (*Client, error) {
	transport, err := rest.TransportFor(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		httpClient: &http.Client{Transport: transport},
		baseURL:    strings.TrimRight(config.Host, "/"),
	}, nil
}
