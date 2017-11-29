package main

import (
	"fmt"
	"os"

	"github.com/xytis/go-disco/discovery"
)

func main() {
	var c *discovery.Client
	var err error
	if c, err = discovery.NewFromEnv(); err != nil {
		panic(fmt.Errorf("discovery client creation failed: %v\n", err))
	}

	if l, err := c.DiscoverOnce(os.Args[1]); err != nil {
		fmt.Printf("Discovery failed: %v\n", err)
	} else {
		fmt.Printf("Discovery success: %v\n", l)
	}

	if d, err := c.Discover(os.Args[1]); err != nil {
		fmt.Printf("Discovery failed: %v\n", err)
	} else {
		for i := 0; i < 5; i++ {
			go func(index int, u <-chan discovery.Change) {
				fmt.Printf("[%v] Update received: %v\n", index, <-u)
			}(i, d.Updates())
		}
		fmt.Printf("holding for last update\n")
		<-d.Updates()
		if err := d.Close(); err != nil {
			fmt.Printf("Error while closing state %v\n", err)
		}
	}
	fmt.Printf("main done\n")
}
