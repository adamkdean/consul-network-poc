//    ___                      _     ___  ___  ___
//   / __\___  _ __  ___ _   _| |   / _ \/___\/ __\
//  / /  / _ \| '_ \/ __| | | | |  / /_)//  // /
// / /__| (_) | | | \__ \ |_| | | / ___/ \_// /___
// \____/\___/|_| |_|___/\__,_|_| \/   \___/\____/
//
// Consul Network proof of concept
// (c) 2018 Adam K Dean

// Package state holds string constants used in
// the DADI Cloud decentralized network.
package state

// State constants.
const (
	Initializing        = "INITIALIZING"
	AwaitingHosts       = "AWAITING_HOSTS"
	WaitingBeforeRetry  = "WAITING_BEFORE_RETRY"
	SearchingForGateway = "SEARCHING_FOR_GATEWAY"
	ConnectingToGateway = "CONNECTING_TO_GATEWAY"
	GatewayConnected    = "GATEWAY_CONNECTED"
	Ready               = "READY"
	Error               = "ERROR"
)
