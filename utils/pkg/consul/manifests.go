//    ___                      _     ___  ___  ___
//   / __\___  _ __  ___ _   _| |   / _ \/___\/ __\
//  / /  / _ \| '_ \/ __| | | | |  / /_)//  // /
// / /__| (_) | | | \__ \ |_| | | / ___/ \_// /___
// \____/\___/|_| |_|___/\__,_|_| \/   \___/\____/
//
// Consul Network proof of concept
// (c) 2018 Adam K Dean

// Package consul provides helper methods to work with Consul
// in context of the DADI Cloud decentralized network.
package consul

// ServiceManifest holds all potential fields that a
// network service may utilise.
type ServiceManifest struct {
	ID         string   `json:"id,omitempty"`
	Type       string   `json:"type,omitempty"`
	Address    string   `json:"address,omitempty"`
	LastActive int64    `json:"last_active,omitempty"`
	Apps       []*App   `json:"apps,omitempty"`
	Hosts      []string `json:"hosts,omitempty"`
}

// App holds basic information about a User app.
type App struct {
	User  string `json:"user,omitempty"`
	Name  string `json:"name,omitempty"`
	Image string `json:"image,omitempty"`
}
