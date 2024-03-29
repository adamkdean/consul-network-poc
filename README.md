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

- [x] Simple "Stargate" app
  - [x] Consul server in `docker-compose.yml`
  - [x] Integrate state machine
  - [x] Update service state on state change
  - [x] Register service with Consul (with state `INITIALIZING`)
  - [x] Create key value pair `stargate/<uuid-v4>`
  - [x] Update key value pair `stargate/<uuid-v4>` `LastActive` every _x_ seconds
  - [x] Prune services where corresponding key value pair last active is > _x_ seconds old

- [x] Simple "Gateway" app
  - [x] Consul agent in `docker-compose.yml`
  - [x] Consul agent connects to Consul server
  - [x] Integrate state machine
  - [x] Update service state on state change
  - [x] Register service with Consul (with state `INITIALIZING`)
  - [x] Create key value pair `gateway/<uuid-v4>`
  - [x] Have list of arbitrary "apps" in key value pair
  - [x] Update key value pair `gateway/<uuid-v4>` `LastActive` every _x_ seconds
  - [x] Set state to `AWAITING_HOSTS`
    - [x] Listen for connections on port _x_ and update key value pair `server` field with address/port
    - [x] On connection, perform simple handsake, update key value pair `hosts` field to include host id

- [x] Simple "Host" app
  - [x] Consul agent in `docker-compose.yml`
  - [x] Consul agent connects to Consul server
  - [x] Integrate state machine
  - [x] Update service state on state change
  - [x] Register service with Consul (with state `INITIALIZING`)
  - [x] Create key value pair `host/<uuid-v4>`
  - [x] Update key value pair `host/<uuid-v4>` `LastActive` every _x_ seconds
  - [x] Set state to `SEARCHING_FOR_GATEWAY`
  - [x] Get all `gateway` services with tag `AWAITING_HOSTS`
  - [x] Attempt connection to gateway, set state to `CONNECTING_TO_GATEWAY`
    - [x] On failure, attempt next, if list exhausted set state to `WAITING_BEFORE_RETRY`, delay, and start again
    - [x] On success, set state to `GATEWAY_CONNECTED`
