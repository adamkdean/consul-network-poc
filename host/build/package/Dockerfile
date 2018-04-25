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

ARG BASE=/go/src/github.com/adamkdean/consul-network-poc/host
ARG TARGET=$BASE/bin/host

# install the latest version of dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# copy local files
WORKDIR ${BASE}
COPY . .

# set the correct permissions on the deploy key
RUN chmod 0600 build/keys/deploy_key

# add deploy key so that we can install private pkgs
# we have to do this in the same session hence one long cmd
RUN mkdir /root/.ssh && \
    eval `ssh-agent -s` && \
    ssh-keyscan -t rsa github.com >> /root/.ssh/known_hosts && \
    ssh-add build/keys/deploy_key && \
    make install

# once we have our deps, let's build
RUN make build

# finally, put the build in an easy to access place
RUN mv ${TARGET} /app

#
# Step 2: Run
#
FROM alpine:latest
COPY --from=build /app app
CMD ["./app"]
