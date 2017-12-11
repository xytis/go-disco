package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/xytis/go-disco/discovery"

	_ "net/http/pprof"
)

func main() {
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	var c discovery.Client
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
		defer d.Close()
		go func(u <-chan discovery.Change) {
			fmt.Printf("Update received: %v\n", <-u)
		}(d.Updates())
		signalChan := make(chan os.Signal, 1)
		cleanupDone := make(chan bool)
		signal.Notify(signalChan, os.Interrupt)
		go func() {
			for _ = range signalChan {
				fmt.Println("\nReceived an interrupt, stopping services...\n")
				cleanupDone <- true
			}
		}()
		<-cleanupDone
	}
	fmt.Printf("main done\n")
}
