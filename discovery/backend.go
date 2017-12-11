package discovery

import (
	"context"
	"fmt"
	"net/url"
)

type Backend interface {
	Name() string
	//Creates a new parralel procedure to watch for changes
	Discover(context.Context, string) (<-chan Change, error)
}

func CreateBackend(descriptor string) (Backend, error) {
	if url, err := url.Parse(descriptor); err != nil {
		return nil, fmt.Errorf("failed to parse url %s", descriptor)
	} else {
		switch url.Scheme {
		case "consul":
			b, e := NewConsulBackend(url)
			return b, e
		case "dns":
			b, e := NewDNSBackend(url)
			return b, e
		default:
			return nil, fmt.Errorf("discovery scheme '%s' not yet supported", url.Scheme)
		}
	}
}
