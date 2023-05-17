package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/criblio/scope-ebpf/internal/ebpf/sigdel"
)

const timeout = 60 * time.Second

// scope-ebpf server will start the loader in loader mode
func loader() {
	// Setup Sigdel
	sd, err := sigdel.Setup()
	if err != nil {
		fmt.Printf("sigdel.Setup failed %v", err)
		return
	}
	defer sd.Teardown()
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGUSR1)

	// Create a channel to implement the timeout
	timeoutChan := time.After(timeout)

	// Teardown procedure
	for {
		select {
		case stopSig := <-stopChan:
			fmt.Println("\nReceived signal:", stopSig.String())
			fmt.Println("\nExiting")
			os.Exit(0)
		case <-timeoutChan:
			fmt.Printf("\nTimeout %v reached. Exiting...", timeout)
			os.Exit(1)
		}
	}
}
