```
#    ___                      _     ___  ___  ___
#   / __\___  _ __  ___ _   _| |   / _ \/___\/ __\
#  / /  / _ \| '_ \/ __| | | | |  / /_)//  // /
# / /__| (_) | | | \__ \ |_| | | / ___/ \_// /___
# \____/\___/|_| |_|___/\__,_|_| \/   \___/\____/
#
# Consul Network proof of concept
# (c) 2018 Adam K Dean
```

# consul-network-poc

Proof of concept network using Consul

## TODO

- [ ] Simple "Stargate" app
  - [ ] Create node "Stargates" if not exist `Consul K/V`
  - [ ] Create node in "Stargates/All" `Consul K/V`
  - [x] Consul server in `docker-compose.yml`
- [ ] Simple "Gateway" app
  - [ ] Create node "Gateways" if not exist `Consul K/V`
  - [ ] Create node in "Gateways/All" `Consul K/V`
  - [ ] Create node in "Gateways/AwaitingHosts" `Consul K/V`
  - [ ] Consul agent in `docker-compose.yml`
- [ ] Simple "Host" app
  - [ ] Create node "Hosts" if not exist `Consul K/V`
  - [ ] Create node in "Hosts/All" `Consul K/V`
  - [ ] Create node in "Hosts/AwaitingGateways" `Consul K/V`
  - [ ] Consul agent in `docker-compose.yml`
