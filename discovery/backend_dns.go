package discovery

import (
	"context"
	"net"
	"net/url"
	"time"
)

type DNSBackend struct {
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

	b := &DNSBackend{
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

func (b *DNSBackend) Discover(ctx context.Context, name string) (<-chan Change, error) {
	output := make(chan Change)

	go func() {
		defer close(output)
		var index uint64
		var channel chan Change
		var change Change
	loop:
		for {
			//Disable sending
			channel = nil
			_, addrs, err := b.c.LookupSRV(ctx, "", "", b.nameToSRV(name))
			if err == nil {
				index = index + 1
				change = Change{
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
				//Enable sending
				channel = output
			}
			select {
			case <-ctx.Done():
				return
			case channel <- change:
				continue loop
			case <-time.After(5 * time.Second):
				continue loop
			}
		}
	}()
	return output, nil
}
