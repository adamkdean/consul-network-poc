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
	// we would get this from some sort of config
	consulAddr := "localhost:8500"

	a := app.New()
	a.Init(consulAddr)
}
