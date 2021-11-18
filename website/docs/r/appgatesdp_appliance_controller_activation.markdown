---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_appliance"
sidebar_current: "docs-appgate-resource-appliance_controller_activation"
description: |-
   Activate an controller
---


# appgatesdp_appliance_controller_activation

Activate controller functionality on a appliance.


~> **NOTE:**  The resource is only available in >= 5.4 appliances. When destroying this resource, it will disable the controller function on the appliance, and the appliance will remain.


!> **Warning:** You won't be able to use this in a safe way on more then 1 appliance at the time.


## Example Usage

```hcl
resource "appgatesdp_appliance_controller_activation" "activate_second_controller" {
  appliance_id = appgatesdp_appliance.second_controller.id
  controller {
    enabled = true
  }
  admin_interface {
    hostname   = "second_controller.com"
    https_port = 8443
    https_ciphers = [
      "ECDHE-RSA-AES256-GCM-SHA384",
      "ECDHE-RSA-AES128-GCM-SHA256"
    ]
  }
}

```


## Argument Reference

The following arguments are supported:

* `controller`: (Required) Controller settings.
* `admin_interface`: (Required) The details of the admin connection interface.


### controller
Controller settings.

* `enabled`:  (Optional)  default value `false` Whether the Controller is enabled on this appliance or not. Cannot be enabled on an inactive Appliance since some checks need to be done first.


### admin_interface
The details of the admin connection interface. Required on Controllers and LogServers.

* `hostname`: (Required) Hostname to connect to the admin interface. This hostname will be used to validate the appliance certificate. Example: appgate.company.com.
* `https_port`:  (Optional)  default value `8443` Port to connect for admin services.
* `https_ciphers`: (Required)  default value `ECDHE-RSA-AES256-GCM-SHA384,ECDHE-RSA-AES128-GCM-SHA256` The type of TLS ciphers to allow. See: https://www.openssl.org/docs/man1.0.2/apps/ciphers.html for all supported ciphers.
* `allow_sources`:  (Optional) Source configuration to allow via iptables.


