package discovery

import (
	"fmt"
	"sync"
)

type discovery struct {
	m           *sync.RWMutex
	in          chan Change
	state       Change
	out         []chan Change
	stopLoop    chan chan error
	stopBackend chan chan error

	//DEBUG: should not be used in non-debug code
	_id string
}

func (d *discovery) Updates() <-chan Change {
	d.m.Lock()
	defer d.m.Unlock()
	ret := make(chan Change)
	d.out = append(d.out, ret)
	return ret
}

func (d *discovery) Close() error {
	errc := make(chan error)
	d.stopLoop <- errc
	return <-errc

}

func (d *discovery) run() {
	for {
		select {
		case c := <-d.in:
			d.broadcast(c)
		case errc := <-d.stopLoop:
			d.stop(errc)
			return
		}
	}
}

func (d *discovery) broadcast(c Change) {
	d.m.Lock()
	defer d.m.Unlock()
	d.state = c
	for _, out := range d.out {
		out <- c
	}
}

func (d *discovery) stop(errc chan error) {
	d.m.Lock()
	defer d.m.Unlock()
	//Stop backend
	errb := make(chan error)
	d.stopBackend <- errb
	errc <- <-errb //Pass any backend error
	for _, out := range d.out {
		close(out)
	}
	close(d.in)
	d.in = nil
	d.out = nil
	return
}

func NewDiscovery(b Backend, id string) (Discovery, error) {
	d := discovery{
		m:           &sync.RWMutex{},
		in:          make(chan Change),
		state:       Change{},
		out:         make([]chan Change, 0),
		stopLoop:    make(chan chan error),
		stopBackend: make(chan chan error),
		_id:         id,
	}
	if err := b.Discover(id, d.in, d.stopBackend); err != nil {
		return nil, fmt.Errorf("failed to register for '%v' discovery in backend '%v'", id, b.Name())
	}
	go d.run()
	return &d, nil
}
