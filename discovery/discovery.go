package discovery

import (
	"context"
	"fmt"
	"sync"
)

type Discovery interface {
	Updates() <-chan Change
	Close()
}

type discovery struct {
	ctx  context.Context
	done context.CancelFunc

	m     *sync.RWMutex
	in    <-chan Change
	state Change
	out   []chan Change

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

func (d *discovery) Close() {
	d.done()
}

func (d *discovery) run() {
	for {
		select {
		case c, ok := <-d.in:
			if !ok {
				return
			}
			d.broadcast(c)
		case <-d.ctx.Done():
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

func NewDiscovery(ctx context.Context, b Backend, id string) (Discovery, error) {
	nctx, cancel := context.WithCancel(ctx)
	if in, err := b.Discover(nctx, id); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to register for '%v' discovery in backend '%v'", id, b.Name())
	} else {
		d := discovery{
			ctx:   nctx,
			done:  cancel,
			m:     &sync.RWMutex{},
			in:    in,
			state: Change{},
			out:   make([]chan Change, 0),
			_id:   id,
		}
		go d.run()
		return &d, nil
	}
}
