---
layout: "appgate"
page_title: "APPGATE: appgate_connector_identity_provider"
sidebar_current: "docs-appgate-resource-connector_identity_provider"
description: |-
   Connector Identity Provider.
---

# appgate_connector_identity_provider

Connector Identity Provider is a builtin default singleton resource, that cannot be deleted.
But we can modifiy the existing one, import the default state from the collective with terraform import.

```bash
$ terraform import 'appgate_connector_identity_provider.connector' connector
```

## Example Usage

```hcl

data "appgate_ip_pool" "ip_v6_pool" {
  ip_pool_name = "default pool v6"
}

data "appgate_ip_pool" "ip_v4_pool" {
  ip_pool_name = "default pool v4"
}

resource "appgate_connector_identity_provider" "connector" {
  ip_pool_v4 = data.appgate_ip_pool.ip_v4_pool.id
  ip_pool_v6 = data.appgate_ip_pool.ip_v6_pool.id
}

```

## Argument Reference

The following arguments are supported:

* `ip_pool_v4`: (Optional) The IPv4 Pool ID the users in this Identity Provider are going to use to allocate IP addresses for the tunnels.
* `ip_pool_v6`: (Optional) The IPv6 Pool ID the users in this Identity Provider are going to use to allocate IP addresses for the tunnels.
* `claim_mappings`: (Optional) The mapping of Identity Provider attributes to claims.
* `on_demand_claim_mappings`: (Optional) The mapping of Identity Provider on demand attributes to claims.


### claim_mappings
The mapping of Identity Provider attributes to claims.

* `attribute_name`:  (Optional) The name of the attribute coming from the Identity Provider.
* `claim_name`:  (Optional) The name of the user claim to be used in AppGate SDP.
* `list`:  (Optional)  default value `false` Whether the claim is expected to be a list and have multiple values or not.
* `encrypt`:  (Optional)  default value `false` Whether the claim should be encrypted or not.
### on_demand_claim_mappings
The mapping of Identity Provider on demand attributes to claims.

* `command`:  (Optional)  Enum values: `fileSize,fileExists,fileCreated,fileUpdated,fileVersion,fileSha512,processRunning,processList,serviceRunning,serviceList,regExists,regQuery,runScript`The name of the command.
* `claim_name`:  (Optional) The name of the device claim to be used in AppGate SDP.
* `parameters`:  (Optional) Depending on the command type, extra parameters to pass to the on-demand claim.
* `platform`:  (Optional)  Enum values: `desktop.windows.all,desktop.macos.all,desktop.linux.all,desktop.all,mobile.android.all,mobile.ios.all,mobile.all,all`The platform(s) to run the on-demand claim.



## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import 'appgate_connector_identity_provider.connector' connector
```
