package discovery

import (
	"fmt"
	"os"
)

type Client struct {
	b Backend
}

func NewFromEnv() (*Client, error) {

	uri := os.Getenv("DISCOVERY_BACKEND")

	if backend, err := CreateBackend(uri); err == nil {
		c := Client{
			b: backend,
		}
		c.init()
		return &c, nil
	} else {
		return nil, fmt.Errorf("could not create backend %v: %v", uri, err)
	}
}

func (c *Client) DiscoverOnce(name string) ([]*Service, error) {
	if d, err := NewDiscovery(c.b, name); err != nil {
		return nil, err
	} else {
		change := <-d.Updates()
		return change.List, nil
	}
}

func (c *Client) Discover(name string) (d Discovery, err error) {
	d, err = NewDiscovery(c.b, name)
	return
}

func (c *Client) init() {

}
