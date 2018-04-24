#    ___                      _     ___  ___  ___
#   / __\___  _ __  ___ _   _| |   / _ \/___\/ __\
#  / /  / _ \| '_ \/ __| | | | |  / /_)//  // /
# / /__| (_) | | | \__ \ |_| | | / ___/ \_// /___
# \____/\___/|_| |_|___/\__,_|_| \/   \___/\____/
#
# Consul Network proof of concept
# (c) 2018 Adam K Dean

#
# Step 1: Build
#
FROM library/golang:1.10 as build

ARG BASE=/go/src/github.com/adamkdean/consul-network-poc/stargate
ARG TARGET=$BASE/bin/stargate

# install the latest version of dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# copy local files
WORKDIR ${BASE}
COPY . .

RUN make install
RUN make build
RUN mv ${TARGET} /app

#
# Step 2: Run
#
FROM alpine:latest
COPY --from=build /app app
CMD ["./app"]
