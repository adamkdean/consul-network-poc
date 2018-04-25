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
	"github.com/adamkdean/consul-network-poc/stargate/internal/app"
)

func main() {
	// create a keepalive channel
	keepalive := make(chan bool)

	// create new instance of app and initialize it
	// in real life, we'd get the address from config
	a := app.New()
	a.Initialize("localhost:8500")

	// live forever
	<-keepalive
}
