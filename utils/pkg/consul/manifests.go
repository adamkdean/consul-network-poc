//    ___                      _     ___  ___  ___
//   / __\___  _ __  ___ _   _| |   / _ \/___\/ __\
//  / /  / _ \| '_ \/ __| | | | |  / /_)//  // /
// / /__| (_) | | | \__ \ |_| | | / ___/ \_// /___
// \____/\___/|_| |_|___/\__,_|_| \/   \___/\____/
//
// Consul Network proof of concept
// (c) 2018 Adam K Dean

package consul

// BasicManifest simply contains an ID, a Service type,
// and an LastActive int64 unix timestamp
type BasicManifest struct {
	ID, Service string
	LastActive  int64
}

// GatewayManifest extends BasicManifest, adding
// an array of GatewayApps
type GatewayManifest struct {
	ID, Service string
	LastActive  int64
	Apps        []*GatewayApp
}

// GatewayApp holds basic information about a User app
type GatewayApp struct {
	User, Name, Image string
}
