---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_dns_policy"
sidebar_current: "docs-appgate-resource-dns-policy"
description: |-
   Create a new DNS Policy.
---

# appgatesdp_dns_policy

Create a new DNS Policy.



## Example Usage

```hcl


resource "appgatesdp_dns_policy" "foobar" {
  name = "tf test dns policy"
  dns_settings {
    domain  = "appgate.com"
    servers = ["8.8.8.8", "1.1.1.1"]
  }
  entitlements = [
    appgatesdp_entitlement.one.id,
    appgatesdp_entitlement.two.id,
  ]
}


```


## Argument Reference

The following arguments are supported:


* `disabled`: (Optional) If true, the Policy will be disregarded during authorization.
* `expression`: (Required) A JavaScript expression that returns boolean. Criteria Scripts may be used by calling them as functions.
* `type`: (Computed) Type of the Policy. It is informational and not enforced.
* `entitlements`: (Optional) List of Entitlement IDs in this Policy.
* `entitlement_links`: (Optional) List of Entitlement tags in this Policy.
* `dns_settings`: (Optional) List of domain names with DNS server IPs that the Client should be using.
* `policy_id`: (Computed) ID of the object.
* `name`: (Required) Name of the object.
* `notes`: (Optional) Notes for the object. Used for documentation purposes.
* `tags`: (Optional) Array of tags.


### entitlements
List of Entitlement IDs in this Policy.

### entitlement_links
List of Entitlement tags in this Policy.


### dns_settings
List of domain names with DNS server IPs that the Client should be using.

* `domain`: (Required) The domain for which the DNS servers should be used by the client.
* `servers`: (Required) undefined


### tags
Array of tags.




## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgatesdp_dns_policy.example d3131f83-10d1-4abc-ac0b-7349538e8300
```
