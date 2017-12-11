package discovery

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"
)

const (
	MAX_TIMEOUT = 5 * time.Second
)

type ConsulBackend struct {
	n *url.URL
	c *api.Client
}

func NewConsulBackend(uri *url.URL) (Backend, error) {
	var config *api.Config
	switch uri.Opaque {
	case "default":
		config = api.DefaultConfig()
	case "":
		config = &api.Config{
			Address:   uri.Host,
			Scheme:    "http",
			Transport: cleanhttp.DefaultPooledTransport(),
		}
	default:
		return nil, fmt.Errorf("could not understand url: '%v'", uri)
	}
	if c, err := api.NewClient(config); err != nil {
		return nil, fmt.Errorf("could not create consul backend: %v", err)
	} else {
		b := ConsulBackend{
			n: uri,
			c: c,
		}
		return &b, nil
	}
}

func (b *ConsulBackend) Name() string {
	return b.n.String()
}

func (b *ConsulBackend) Discover(ctx context.Context, name string) (<-chan Change, error) {
	output := make(chan Change)
	c := b.c.Catalog()
	qo := &api.QueryOptions{
		WaitIndex: 0,
	}

	go func() {
		defer close(output)
		var timeout time.Duration
		var channel chan Change
		var change Change
	loop:
		for {
			//Disable sending
			channel = nil
			//TODO: alter consul lib to use ctx
			s, m, err := c.Service(name, "", qo)
			if err != nil {
				output = nil
				timeout = timeout + 1*time.Second
				if timeout > MAX_TIMEOUT {
					timeout = MAX_TIMEOUT
				}
			} else {
				timeout = 0 * time.Second

				qo.WaitIndex = m.LastIndex
				change = Change{
					Index: m.LastIndex,
					List:  make([]*Service, len(s)),
				}
				for i, service := range s {
					addr := service.ServiceAddress
					if addr == "" {
						addr = service.Address
					}
					change.List[i] = &Service{
						Address: addr,
						Port:    service.ServicePort,
					}
				}
				//Enable sending
				channel = output
			}
			select {
			case <-ctx.Done():
				return
			case channel <- change:
				continue loop
			case <-time.After(timeout):
				continue loop
			}
		}
	}()
	return output, nil
}
