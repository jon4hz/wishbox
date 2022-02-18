# Wishbox

Generate a wishlist directory based on your netbox inventory.

## How does it work?
When starting wishbox, it queries the netbox api and generates the wishlist endpoints.  
Wishbox will use the devices primary IP to connect to.  
The ssh port is 22 by default. To overwrite this, you can define a netbox service for the host called `ssh`.

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
  # the netbox api token
  token: supersecretapitoken
  # the user for the ssh connection (default is your current system user)
  user: toor
  # forward the local ssh agent?
  forward_agent: yes
  # only list devices which have this role assigned
  filter_role: linux_server
```

## Limitations
- Pagination isn't implemented (yet), so wishbox will return only the first 50 devices