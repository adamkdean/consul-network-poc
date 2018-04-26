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
  - [x] Register service with Consul (with state `INITIALIZED`)
  - [x] Create key value pair `stargate/<uuid-v4>`
  - [x] Update key value pair `stargate/<uuid-v4>` `LastActive` every _x_ seconds
  - [ ] Prune services where corresponding key value pair last active is > _x_ seconds old

- [ ] Simple "Gateway" app
  - [x] Consul agent in `docker-compose.yml`
  - [x] Consul agent connects to Consul server
  - [x] Register service with Consul (with state `INITIALIZED`)
  - [x] Create key value pair `gateway/<uuid-v4>`
  - [ ] Have list of arbitrary "apps" in key value pair
  - [x] Update key value pair `gateway/<uuid-v4>` `LastActive` every _x_ seconds
  - [ ] Set state to AWAITING_HOSTS
    - [ ] Listen for connections on port _x_ and update key value pair `server` field with address/port
    - [ ] On connection, perform simple handsake, update key value pair `hosts` field to include host id

- [ ] Simple "Host" app
  - [x] Consul agent in `docker-compose.yml`
  - [x] Consul agent connects to Consul server
  - [x] Register service with Consul (with state `INITIALIZED`)
  - [x] Create key value pair `host/<uuid-v4>`
  - [x] Update key value pair `host/<uuid-v4>` `LastActive` every _x_ seconds
  - [ ] Set state to `AWAITING_GATEWAY`
  - [ ] Get all `gateway` services with tag `AWAITING_HOSTS`
  - [ ] Attempt connection to gateway, set state to `CONNECTING_TO_GATEWAY`
    - [ ] On failure, attempt next, if list exhausted set state to `AWAITING_GATEWAY`, delay, and start again
    - [ ] On success, set state to `GATEWAY_CONNECTED`
      - [ ] Read key value pair for `gateway/<uuid-v4>`, get `apps` field
      - [ ] Pretend we're A-OK and set State to `READY`
