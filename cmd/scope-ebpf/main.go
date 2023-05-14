package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/criblio/scope-ebpf/internal/ebpf/oom"
	"github.com/criblio/scope-ebpf/internal/prom"
)

func main() {
	var scopepromserver string

	if os.Geteuid() != 0 {
		fmt.Println("This binary must be run with sudo for elevated privileges.")
		return
	}
	// fmt.Println(os.Args[0], "started, PID:", os.Getpid())
	flag.StringVar(&scopepromserver, "scopepromserver", "", "Scope Prometheus server")
	flag.Parse()

	oomEventChan := make(chan string, 25)
	go oom.SetupReadOOM(oomEventChan)

	oomElement := prom.PromMetricCounter{
		Name:    "oom_kill",
		Counter: 0,
		Unit:    "process",
	}
	for {
		select {
		case oomEvent := <-oomEventChan:
			oomElement.Valobj = oomEvent
			oomElement.Add()
			msg := oomElement.String()
			fmt.Println("Prepare message", msg)
			conn, err := net.Dial("tcp", scopepromserver)
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
			}
			fmt.Println("Send message", msg)
		}
	}

}
