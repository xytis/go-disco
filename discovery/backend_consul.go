package discovery

import (
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
	fmt.Printf("uri:opaque %v\nuri:host %v\n", uri.Opaque, uri.Host)
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

func (b *ConsulBackend) Discover(name string, ch chan<- Change, done <-chan chan error) error {
	c := b.c.Catalog()
	qo := &api.QueryOptions{
		WaitIndex: 0,
	}
	go func() {
		var lerr error
		var output chan<- Change
		var timeout time.Duration
		var change Change
	loop:
		for {
			fmt.Printf("looking up a service\n")
			s, m, err := c.Service(name, "", qo)
			if err != nil {
				fmt.Printf("err: %v\n", err)
				lerr = err
				output = nil
				timeout = timeout + 1*time.Second
				if timeout > MAX_TIMEOUT {
					timeout = MAX_TIMEOUT
				}
			} else {
				output = ch
				timeout = 0 * time.Second

				qo.WaitIndex = m.LastIndex
				change = Change{
					Index: m.LastIndex,
					List:  make([]*Service, len(s)),
				}
				for i, service := range s {
					change.List[i] = &Service{
						Address: service.ServiceAddress,
						Port:    service.ServicePort,
					}
				}
			}
			for {
				fmt.Printf("trying to put output\n")
				select {
				case errc := <-done:
					fmt.Printf("closed\n")
					errc <- lerr
					return
				case output <- change:
					fmt.Printf("changes delivered\n")
					continue loop
				case <-time.After(timeout):
					fmt.Printf("timeout reached: %v\n", timeout.Seconds())
					continue loop
				}
			}
		}
	}()
	return nil
}
