package main

import (
	"fmt"
	"net"
	"time"

	"github.com/criblio/scope-ebpf/internal/ebpf/oom"
	"github.com/criblio/scope-ebpf/internal/prom"
)

type serverCfg struct {
	address string
	debug   bool
}

// scope-ebpf server will start the loader in server mode
func server(cfg serverCfg) {
	if cfg.debug {
		fmt.Println("server configured to listening on:", cfg.address)
	}
	// Setup OOM
	oomEventChan := make(chan string, 25)
	go oom.Setup(oomEventChan)

	oomElement := prom.PromMetricCounter{
		Name:    "oom_kill",
		Counter: 0,
		Unit:    "process",
	}

	for {
		select {
		case oomEvent := <-oomEventChan:
			if cfg.debug {
				fmt.Println("oomEvent happened ", oomEvent)
			}
			oomElement.Valobj = oomEvent
			oomElement.Add()
			msg := oomElement.String()
			conn, err := net.Dial("tcp", cfg.address)
			if err != nil {
				fmt.Println("Failed to connect:", err)
				time.Sleep(time.Second)
				continue
			}
			defer conn.Close()
			_, err = conn.Write([]byte(msg))
			if err != nil {
				fmt.Println("Failed to send message:", err)
				return
			} else {
				if cfg.debug {
					fmt.Println("message sent", msg)
				}
			}

		}
	}
}
