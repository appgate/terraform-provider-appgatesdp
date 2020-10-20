terraform {
  required_version = ">= 0.12"
}
provider "appgate" {
  username = "admin"
  password = "admin"
  url      = "https://${var.controller_dns}:444/admin"
  provider = "local"
  insecure = true
}

data "appgate_site" "default_site" {
  site_name = "Default site"
}

resource "appgate_appliance" "new_gateway" {
  name     = replace(var.gateway_dns, ".", "_")
  hostname = var.gateway_dns

  client_interface {
    hostname       = var.gateway_dns
    proxy_protocol = true
    https_port     = 8443
    dtls_port      = 443
    allow_sources {
      address = "0.0.0.0"
      netmask = 0
    }
    allow_sources {
      address = "::"
      netmask = 0
    }
    override_spa_mode = "Disabled"
  }

  peer_interface {
    hostname   = var.gateway_dns
    https_port = "444"

    allow_sources {
      address = "0.0.0.0"
      netmask = 0
    }
    allow_sources {
      address = "::"
      netmask = 0
    }
  }


  admin_interface {
    hostname = var.gateway_dns
    https_ciphers = [
      "ECDHE-RSA-AES256-GCM-SHA384",
      "ECDHE-RSA-AES128-GCM-SHA256"
    ]
  }

  tags = [
    "terraform",
    "api-created"
  ]
  notes = "hello world"
  site  = data.appgate_site.default_site.id


  networking {


    nics {
      enabled = true
      name    = "eth0"
      ipv4 {
        dhcp {
          enabled = true
          dns     = true
          routers = true
          ntp     = true
        }
      }
    }

  }

  # https://sdphelp.appgate.com/adminguide/v5.1/about-appliances.html?anchor=gateway-a
  gateway {
    enabled = true
    vpn {
      weight = 100
      allow_destinations {
        address = "0.0.0.0"
        nic     = "eth0"
      }
    }
  }

}

data "appgate_appliance_seed" "gateway_seed_file" {
  depends_on = [
    appgate_appliance.new_gateway,
  ]
  appliance_id   = appgate_appliance.new_gateway.id
  password       = "cz"
  latest_version = true
}

resource "null_resource" "seed_gateway" {

  depends_on = [
    data.appgate_appliance_seed.gateway_seed_file,
  ]

  connection {
    type        = "ssh"
    user        = "cz"
    timeout     = "25m"
    private_key = file(var.private_key)
    host        = var.gateway_dns
  }


  provisioner "local-exec" {
    command = "echo ${data.appgate_appliance_seed.gateway_seed_file.seed_file} > seed.b64"
  }
  provisioner "remote-exec" {
    inline = [
      "echo ${data.appgate_appliance_seed.gateway_seed_file.seed_file}  > raw.b64",
      "cat raw.b64 | base64 -d  | jq .  >> seed.json",
    ]
  }
}
