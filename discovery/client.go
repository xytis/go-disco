package discovery

import (
	"context"
	"fmt"
	"os"
)

type Client interface {
	DiscoverOnce(string) ([]*Service, error)
	Discover(string) (Discovery, error)
	Close() error
}

type client struct {
	ctx    context.Context
	cancel context.CancelFunc
	b      Backend
}

func NewFromEnv() (Client, error) {

	uri := os.Getenv("DISCOVERY_BACKEND")
	ctx, cancel := context.WithCancel(context.TODO())

	if backend, err := CreateBackend(uri); err == nil {
		c := client{
			ctx:    ctx,
			cancel: cancel,
			b:      backend,
		}
		return &c, nil
	} else {
		return nil, fmt.Errorf("could not create backend %v: %v", uri, err)
	}
}

func (c *client) DiscoverOnce(name string) ([]*Service, error) {
	if d, err := NewDiscovery(c.ctx, c.b, name); err != nil {
		return nil, err
	} else {
		change := <-d.Updates()
		return change.List, nil
	}
}

func (c *client) Discover(name string) (d Discovery, err error) {
	d, err = NewDiscovery(c.ctx, c.b, name)
	return
}

func (c *client) Close() error {
	return nil

}
