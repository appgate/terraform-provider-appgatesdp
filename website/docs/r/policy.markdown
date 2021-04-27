---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_policy"
sidebar_current: "docs-appgate-resource-policy"
description: |-
   Create a new Policy.
---

# appgatesdp_policy

Create a new Policy.

## Example Usage

```hcl

resource "appgatesdp_policy" "basic_policy" {
  name  = "terraform policy"
  notes = "terraform policy notes"
  tags = [
    "terraform",
    "api-created"
  ]
  disabled = false

  expression = <<-EOF
var result = false;
/*claims.user.groups*/
if(claims.user.groups && claims.user.groups.indexOf("developers") >= 0) {
  return true;
}
/*end claims.user.groups*/
/*criteriaScript*/
if (admins(claims)) {
  return true;
}
/*end criteriaScript*/
return result;
EOF
}

```

Example of a looping policy for multiple claims of a single type

```hcl

locals {
  groups = ["developers", "admins"]
}

resource "appgatesdp_policy" "looping_policy" {
  name  = "terraform policy"
  notes = "terraform policy notes"
  tags = [
    "terraform",
    "api-created"
  ]
  disabled = false

  expression = <<-EOF
var result = false;
%{for group in local.groups~}
if/*claims.user.groups*/(claims.user.groups && claims.user.groups.indexOf("${group}") >= 0)/*end claims.user.groups*/ { return true; }
${endfor~}
return result;
EOF
}

```

## Argument Reference

The following arguments are supported:


* `disabled`: (Optional) If true, the Policy will be disregarded during authorization.
* `expression`: (Required) A JavaScript expression that returns boolean. Criteria Scripts may be used by calling them as functions.
* `entitlements`: (Optional) List of Entitlement IDs in this Policy.
* `entitlement_links`: (Optional) List of Entitlement tags in this Policy.
* `ringfence_rules`: (Optional) List of Ringfence Rule IDs in this Policy.
* `ringfence_rule_links`: (Optional) List of Ringfence Rule tags in this Policy.
* `tamper_proofing`: (Optional) Will enable Tamper Proofing on desktop clients which will make sure the routes and ringfence configurations are not changed.
* `override_site`: (Optional) Site ID where all the Entitlements of this Policy must be deployed. This overrides Entitlement's own Site and to be used only in specific network layouts. Otherwise the assigned site on individual Entitlements will be used.
* `proxy_auto_config`: (Optional) Client configures PAC URL on the client OS.
* `trusted_network_check`: (Optional) Client suspends operations when it's in a trusted network.
* `administrative_roles`: (Optional) List of Administrative Role IDs in this Policy.
* `id`: (Required) ID of the object.
* `name`: (Required) Name of the object.
* `notes`: (Optional) Notes for the object. Used for documentation purposes.
* `tags`: (Optional) Array of tags.


### entitlements
List of Entitlement IDs in this Policy.

### entitlement_links
List of Entitlement tags in this Policy.

### ringfence_rules
List of Ringfence Rule IDs in this Policy.

### ringfence_rule_links
List of Ringfence Rule tags in this Policy.

### proxy_auto_config
Client configures PAC URL on the client OS.

* `enabled`:  (Optional)
* `url`:  (Optional) The URL to set on the Client OS. Example: https://pac.company.com/file.pac.
* `persist`:  (Optional) If true Client will leave the PAC URL configured after signing out.
### trusted_network_check
Client suspends operations when it's in a trusted network.

* `enabled`:  (Optional)
* `dns_suffix`:  (Optional) Client checks if the DNS suffix has been configured on the OS to decide whether it's on a trusted network or not.
### administrative_roles
List of Administrative Role IDs in this Policy.

### tags
Array of tags.




## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgatesdp_policy d3131f83-10d1-4abc-ac0b-7349538e8300
```
