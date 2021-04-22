---
layout: "appgatesdp"
page_title: "Provisioning appliances"
sidebar_current: "docs-appgatesdp-guide-state_migrate"
description: |-
  Provisioning appliances
---

## Provisioning appliances



### cz-seed

Appliances are provisioned with the [seed file](https://sdphelp.appgate.com/adminguide/v5.3/new-appliance.html?anchor=manual-seeding), we use `cz-seed` to the provision the appliances. `cz-seed` is a built-in program included on every appliance.
This is a non interactive version of [cz-setup](https://sdphelp.appgate.com/adminguide/v5.3/controller-cz-setup.html). 


`cz-seed` has two main cases:
- Initial seeding of the first controller, setup initial network configuration.
- Initial configuration for appliances to setup network configuration before retrieving the seed file.

`cz-seed` can be used with all major cloud providers, such as aws, azure, gcp, openstack, where cloud-init usually is used.



```bash
$ cz-seed --help
usage: cz-seed [-h] [-v] [-i] [-o FILE] [-hh HOSTNAME] [-ch {HOSTNAME|IP}]
               [-ph {HOSTNAME|IP}] [-ah {HOSTNAME|IP}]
               [--profile-hostname HOSTNAME] [-d4 NIC] [-d6 NIC]
               [-a4 NIC IP MASK] [-a6 NIC IP MASK] [-g4 IP] [-g6 IP]
               [-rt IP MASK GATEWAY NIC] [-he HOSTNAME IP] [-ds IP]
               [-dd DOMAIN] [-ns {HOSTNAME|IP}]
               [-rd SELECTOR TEMPLATE DESTINATION] [-rh KEY_FILE PUB_FILE]
               [-hc] [-eg] [-et] [-ed] [-es] [-eh] [-gn NIC] [-el] [-dc]
               [-pw USERNAME PASSWORD] [-np USERNAME] [-ak USERNAME FILE]
               [-ap ADMIN_PASSWORD] [-nr] [-ll SERVICENAME LOGLEVEL]

optional arguments:
  -h, --help            show this help message and exit
  -v, --version         print version and exit
  -i, --interactive     use questions rather than options
  -o FILE, --output FILE
  -hh HOSTNAME, --hostname HOSTNAME
  -ch {HOSTNAME|IP}, --client-hostname {HOSTNAME|IP}
  -ph {HOSTNAME|IP}, --peer-hostname {HOSTNAME|IP}, --appliance-hostname {HOSTNAME|IP}
  -ah {HOSTNAME|IP}, --admin-hostname {HOSTNAME|IP}
  --profile-hostname HOSTNAME
  -d4 NIC, --dhcp-ipv4 NIC
  -d6 NIC, --dhcp-ipv6 NIC
  -a4 NIC IP MASK, --address-ipv4 NIC IP MASK
  -a6 NIC IP MASK, --address-ipv6 NIC IP MASK
  -g4 IP, --default-gateway-ipv4 IP
  -g6 IP, --default-gateway-ipv6 IP
  -rt IP MASK GATEWAY NIC, --route IP MASK GATEWAY NIC
  -he HOSTNAME IP, --hosts-entry HOSTNAME IP
  -ds IP, --dns-server IP
  -dd DOMAIN, --dns-domain DOMAIN
  -ns {HOSTNAME|IP}, --ntp-server {HOSTNAME|IP}
  -rd SELECTOR TEMPLATE DESTINATION, --rsyslog-destination SELECTOR TEMPLATE DESTINATION
  -rh KEY_FILE PUB_FILE, --rsa-hostkey KEY_FILE PUB_FILE
  -hc, --enable-healthcheck
  -eg, --enable-gateway
  -et, --enable-tls
  -ed, --enable-dtls
  -es, --enable-snat
  -eh, --enable-state-sharing
  -gn NIC, --gateway-nic NIC
  -el, --enable-logserver
  -dc, --disable-controller
  -pw USERNAME PASSWORD, --password USERNAME PASSWORD
  -np USERNAME, --no-password USERNAME
  -ak USERNAME FILE, --authorized-keys USERNAME FILE
  -ap ADMIN_PASSWORD, --admin-password ADMIN_PASSWORD
  -nr, --no-registration
  -ll SERVICENAME LOGLEVEL, --log-level SERVICENAME LOGLEVEL
                        set the log level for the service specified

These options can be specified multiple times:
    --dhcp-ipv4
    --dhcp-ipv6
    --address-ipv4
    --address-ipv6
    --route
    --hosts-entry
    --dns-server
    --dns-domain
    --ntp-server
    --rsyslog-destination
    --gateway-nic
    --password
    --authorized-keys
    --log-level

```

### seed file

The seed file is used on a appliance to join the collective. We need to send the seed file to the appliance after we have created an inactive appliance.
- https://sdphelp.appgate.com/adminguide/v5.3/new-appliance.html?anchor=manual-seeding


### Examples


#### Provision with seed file

First we need to create an inactive [appliance configuration](../r/appliance.markdown) 
```hcl
resource "appgatesdp_appliance" "new_gateway" {
  // ...
}

```

Second part, is to take the [seed file](../d/appgate_appliance_seed.markdown) from the inactive appliance we created and provision the instance.

```hcl

data "appgatesdp_appliance_seed" "gateway_seed_file" {
  depends_on = [
    appgatesdp_appliance.new_gateway,
  ]
  appliance_id   = appgatesdp_appliance.new_gateway.id
  password       = "cz"
  latest_version = true
}

resource "null_resource" "seed_gateway" {
  depends_on = [
    data.appgatesdp_appliance_seed.gateway_seed_file,
  ]

  connection {
    type        = "ssh"
    user        = "cz"
    timeout     = "25m"
    private_key = file(var.private_key)
    host        = var.gateway_dns
  }


  provisioner "local-exec" {
    command = "echo ${data.appgatesdp_appliance_seed.gateway_seed_file.seed_file} > seed.b64"
  }
  provisioner "remote-exec" {
    inline = [
      # any file named seed.json in the home directory will be picked up and treated as configuration file by the configuration daemon.
      # if the configuration is invalid or corrupt, a log entry will be written to journalctl
      #
      # the data.appgatesdp_appliance_seed.gateway_seed_file.seed_file is base 64 encoded
      # so we need to decode it before we write.
      "echo ${data.appgatesdp_appliance_seed.gateway_seed_file.seed_file}  > raw.b64",
      "cat raw.b64 | base64 -d  | jq .  >> seed.json",
    ]
  }
}
```

#### Provision with cz-seed

Depending on your environment, you could use userdata to provision the instance on startup, for example in aws we want to setup an initial controller with `cz-seed`


```hcl

locals {
  controller_user_data = <<-EOF
#!/bin/bash
PUBLIC_HOSTNAME=`curl --silent http://169.254.169.254/latest/meta-data/public-hostname`
# seed the first controller, and enable admin interface on :8443
cz-seed \
    --password cz cz \
    --dhcp-ipv4 eth0 \
    --enable-logserver \
    --no-registration \
    --hostname "$PUBLIC_HOSTNAME" \
    --admin-password ${var.admin_login_password} \
    | jq '.remote.adminInterface.hostname = .remote.peerInterface.hostname | .remote.adminInterface.allowSources = .remote.peerInterface.allowSources' >> /home/cz/seed.json
EOF
}

resource "aws_instance" "appgatesdp_controller" {
   // ...
   user_data_base64 = base64encode(local.controller_user_data)
}
```

### Documentation

For more documentation how to configure Appgate, please see 
- https://sdphelp.appgate.com/ - the official manual
- https://github.com/appgate/sdp-tf-reference-architecture - reference project
