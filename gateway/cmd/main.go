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
	"github.com/adamkdean/consul-network-poc/gateway/internal/app"
	"os"
)

func main() {
	// create a keepalive channel
	keepalive := make(chan bool)

	// determine consul server location
	consulAddr := os.Getenv("CONSUL_ADDRESS")
	if consulAddr == "" {
		consulAddr = "localhost:8500"
	}

	// determine listen address
	listenAddr := os.Getenv("LISTEN_ADDRESS")
	if listenAddr == "" {
		listenAddr = "0.0.0.0:8000"
	}

	// create new instance of app and initialize it
	a := app.New()
	a.Initialize(consulAddr, listenAddr)

	// live forever
	<-keepalive
}
