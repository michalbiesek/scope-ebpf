package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var cfg serverCfg

	if os.Geteuid() != 0 {
		fmt.Println("This binary must be run with sudo for elevated privileges.")
		return
	}

	// fmt.Println(os.Args[0], "started, PID:", os.Getpid())
	flag.StringVar(&cfg.address, "scopepromserver", "", "Scope Prometheus server")
	flag.BoolVar(&cfg.debug, "debug", false, "Enable debug message")
	flag.Parse()

	if cfg.address != "" {
		server(cfg)
	} else {
		loader()
	}
}
