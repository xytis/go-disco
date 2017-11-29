package discovery

type Service struct {
	Address string
	Port    int
}

type Change struct {
	Index uint64
	List  []*Service
}

type Discovery interface {
	Updates() <-chan Change
	Close() error
}

type Backend interface {
	Name() string
	//Creates a new parralel procedure to watch for changes
	Discover(string, chan<- Change, <-chan chan error) error
}
