//    ___                      _     ___  ___  ___
//   / __\___  _ __  ___ _   _| |   / _ \/___\/ __\
//  / /  / _ \| '_ \/ __| | | | |  / /_)//  // /
// / /__| (_) | | | \__ \ |_| | | / ___/ \_// /___
// \____/\___/|_| |_|___/\__,_|_| \/   \___/\____/
//
// Consul Network proof of concept
// (c) 2018 Adam K Dean

package main

import (
	"github.com/adamkdean/consul-network-poc/host/internal/app"
	"os"
)

func main() {
	// create a keepalive channel
	keepalive := make(chan bool)

	// determine consul server location
	addr := os.Getenv("CONSUL_ADDRESS")
	if addr == "" {
		addr = "localhost:8500"
	}

	// create new instance of app and initialize it
	a := app.New()
	a.Initialize(addr)

	// live forever
	<-keepalive
}
