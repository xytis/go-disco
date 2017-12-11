package discovery

import (
	"fmt"
)

type Service struct {
	Address string
	Port    int
}

func (s Service) String() string {
	return fmt.Sprintf("[Service %s:%d]", s.Address, s.Port)
}

type Change struct {
	Index uint64
	List  []*Service
}

func (c Change) String() string {
	return fmt.Sprintf("[Change %d %+v]", c.Index, c.List)
}
