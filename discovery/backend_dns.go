package discovery

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"time"
)

type DNSBackend struct {
	ctx    context.Context
	cancel context.CancelFunc

	n      *url.URL
	c      *net.Resolver
	ns     string
	suffix string
}

func NewDNSBackend(uri *url.URL) (Backend, error) {
	//To imlement:
	// * Multiple nameserver support
	// * Options:
	//     timeout=[default 5]
	//     protocol=[default udp]

	ctx, cancel := context.WithCancel(context.TODO())
	b := &DNSBackend{
		ctx:    ctx,
		cancel: cancel,
		n:      uri,
		c:      new(net.Resolver),
		ns:     uri.Host,
		suffix: uri.Query().Get("suffix"),
	}
	return b, nil
}

func (b *DNSBackend) nameToSRV(name string) string {
	if b.suffix == "" {
		return name
	}
	return name + "." + b.suffix
}

func (b *DNSBackend) Name() string {
	return b.n.String()
}

func (b *DNSBackend) Discover(name string, ch chan<- Change, done <-chan chan error) error {
	go func() {
		var index uint64
	loop:
		for {
			_, addrs, err := b.c.LookupSRV(b.ctx, "", "", b.nameToSRV(name))
			if err == nil {
				index = index + 1
				change := Change{
					Index: index,
					List:  make([]*Service, len(addrs)),
				}
				//Optionaly unsort SRV record
				for i, srv := range addrs {
					change.List[i] = &Service{
						Address: srv.Target,
						Port:    int(srv.Port),
					}
				}
				select {
				case ch <- change:
				case errc := <-done:
					fmt.Printf("closed\n")
					errc <- nil
					return
				}

			}
			select {
			case <-b.ctx.Done():
				return
			case <-time.After(5 * time.Second):
				continue loop
			}
		}
	}()
	return nil
}
