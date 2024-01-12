---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_connector_identity_provider"
sidebar_current: "docs-appgate-resource-connector_identity_provider"
description: |-
  Import and Update Connector Identity Provider.
---

# appgatesdp_connector_identity_provider

~> **NOTE:** Connector Identity Provider is a builtin default singleton resource, that cannot be deleted. But we can modify the existing one, import the default state from the collective with terraform import.

~> **NOTE:**  The resource documentation is based on the latest available appgate sdp appliance version, which currently is 6.2
Some attributes may not be available if you are running an older version, if you try to use an attribute block that is not permitted in your current version, you will be prompted by an error message.


## Example Usage

```hcl


data "appgatesdp_ip_pool" "ip_v6_pool" {
  ip_pool_name = "default pool v6"
}

data "appgatesdp_ip_pool" "ip_v4_pool" {
  ip_pool_name = "default pool v4"
}

resource "appgatesdp_connector_identity_provider" "connector" {
  ip_pool_v4 = data.appgatesdp_ip_pool.ip_v4_pool.id
  ip_pool_v6 = data.appgatesdp_ip_pool.ip_v6_pool.id
}


```

### Connector Identity Provider Arguments
* `name`: (Required) Name of the object.
* `notes`: (Optional) Notes for the object. Used for documentation purposes.
* `tags`: (Optional) Array of tags.
* `type`: (Computed) The type of the Identity Provider.
* `ip_pool_v4`: (Optional) The IPv4 Pool ID the users in this Identity Provider are going to use to allocate IP addresses for the tunnels.
* `ip_pool_v6`: (Optional) The IPv6 Pool ID the users in this Identity Provider are going to use to allocate IP addresses for the tunnels.
* `claim_mappings`: (Optional) The mapping of Identity Provider attributes to claims.
  * `attribute_name`: (Required) The name of the attribute coming from the Identity Provider.
  * `claim_name`: (Required) The name of the user claim to be used in Appgate SDP.
  * `list`:  (Optional)  default value `false` Whether the claim is expected to be a list and have multiple values or not.
  * `encrypt`:  (Optional)  default value `false` Whether the claim should be encrypt or not.
* `user_scripts`: (Optional) ID of the User Claim Scripts to run during authorization.

## Import
Instances can be imported using the `id`, e.g.

```
$ terraform import appgatesdp_connector_identity_provider.example d3131f83-10d1-4abc-ac0b-7349538e8300
```
