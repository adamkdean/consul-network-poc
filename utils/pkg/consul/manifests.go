//    ___                      _     ___  ___  ___
//   / __\___  _ __  ___ _   _| |   / _ \/___\/ __\
//  / /  / _ \| '_ \/ __| | | | |  / /_)//  // /
// / /__| (_) | | | \__ \ |_| | | / ___/ \_// /___
// \____/\___/|_| |_|___/\__,_|_| \/   \___/\____/
//
// Consul Network proof of concept
// (c) 2018 Adam K Dean

package consul

// StargateManifest contains the data fields
// pertaining to a network Stargate
type StargateManifest struct {
	ID, Service string
	LastActive  int64
}

// GatewayManifest contains the data fields
// pertaining to a network Gateway
type GatewayManifest struct {
	ID, Service string
	LastActive  int64
	Address     string
	Port        int
	Apps        []*GatewayApp
	Hosts       []string
}

// GatewayApp holds basic information about a User app
type GatewayApp struct {
	User, Name, Image string
}

// HostManifest contains the data fields
// pertaining to a network Host
type HostManifest struct {
	ID, Service string
	LastActive  int64
}
