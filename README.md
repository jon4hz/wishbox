# Wishbox
[![goreleaser](https://github.com/jon4hz/wishbox/actions/workflows/goreleaser.yml/badge.svg)](https://github.com/jon4hz/wishbox/actions/workflows/goreleaser.yml)
[![lint](https://github.com/jon4hz/wishbox/actions/workflows/lint.yml/badge.svg)](https://github.com/jon4hz/wishbox/actions/workflows/lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/jon4hz/wishbox)](https://goreportcard.com/report/github.com/jon4hz/wishbox)


Generate a wishlist directory based on your [netbox](https://github.com/netbox-community/netbox) inventory.

## How does it work?
When starting wishbox, it queries the netbox api and generates the wishlist endpoints.  
Wishbox will use the devices primary IP to connect to.  
The ssh port is 22 by default. To overwrite this, you can define a netbox service for the host called `ssh`.

## Installation
### Docker-compose
```yaml
---
version: "3.8"
services:
  wishbox:
    image: ghcr.io/jon4hz/wishbox:latest
    restart: unless-stopped
    volumes:
      - ./config.yml:/app/config.yml
      - .wishlist:/app/.wishlist
    ports:
      - "22:2223"
```
### Build from source
```
git clone https://github.com/jon4hz/wishbox.git
go build .
./wishbox
```

## Configuration
The configuration is loaded from the `./config.yml` file by default.

```yaml
# the ip the wishlist server listens on
listen: 127.0.0.1

# the port wishlist uses
port: 2223

# netbox configuration
netbox:
  # your netbox host
  host: my.netbox.net
  # set to true to disable tls validation
  ignore_tls: false
  # the netbox api token
  token: supersecretapitoken
  # the user for the ssh connection (default is your current system user)
  user: toor
  # forward the ssh agent?
  forward_agent: yes
  # only list devices which have this role assigned
  filter_role: linux_server
  # list only devices that are active inside netbox
  only_active: yes
```

## Limitations
- The inventory is only generated when starting wishbox
- Pagination isn't implemented (yet), so wishbox will return only the first 50 devices