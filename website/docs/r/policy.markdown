---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_policy"
sidebar_current: "docs-appgate-resource-policy"
description: |-
   Create a new Policy.
---

# appgatesdp_policy

Create a new Policy.

~> **NOTE:**  The resource documentation is based on the latest available appgate sdp appliance version, which currently is 5.5.0
Some attributes may not be available if you are running an older version, if you try to use an attribute block that is not permitted in your current version, you will be prompted by an error message.


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


## Argument Reference

The following arguments are supported:


* `disabled`: (Optional) If true, the Policy will be disregarded during authorization.
* `expression`: (Required) A JavaScript expression that returns boolean. Criteria Scripts may be used by calling them as functions.
* `type`: (Optional) Type of the Policy. It is informational and not enforced.
* `entitlements`: (Optional) List of Entitlement IDs in this Policy.
* `entitlement_links`: (Optional) List of Entitlement tags in this Policy.
* `ringfence_rules`: (Optional) List of Ringfence Rule IDs in this Policy.
* `ringfence_rule_links`: (Optional) List of Ringfence Rule tags in this Policy.
* `tamper_proofing`: (Optional) Will enable Tamper Proofing on desktop clients which will make sure the routes and ringfence configurations are not changed.
* `override_site`: (Optional) Site ID where all the Entitlements of this Policy must be deployed. This overrides Entitlement's own Site and to be used only in specific network layouts. Otherwise the assigned site on individual Entitlements will be used.
* `override_site_claim`: (Optional) The path of a claim that contains the UUID of an override site. It should be defined as "claims.xxx.xxx" or "claims.xxx.xxx.xxx".
* `proxy_auto_config`: (Optional) Client configures PAC URL on the client OS.
* `trusted_network_check`: (Optional) Client suspends operations when it's in a trusted network.
* `dns_settings`: (Optional) List of domain names with DNS server IPs that the Client should be using.
* `client_settings`: (Optional) Settings that admins can apply to the Client.
* `administrative_roles`: (Optional) List of Administrative Role IDs in this Policy.
* `policy_id`: (Optional) ID of the object.
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

* `enabled`:  (Optional)  default value `false`
* `url`:  (Optional) The URL to set on the Client OS. Example: https://pac.company.com/file.pac.
* `persist`:  (Optional) If true Client will leave the PAC URL configured after signing out.
### trusted_network_check
Client suspends operations when it's in a trusted network.

* `enabled`:  (Optional)  default value `false`
* `dns_suffix`:  (Optional) Client checks if the DNS suffix has been configured on the OS to decide whether it's on a trusted network or not.
### dns_settings
List of domain names with DNS server IPs that the Client should be using.

* `domain`: (Required) The domain for which the DNS servers should be used by the client.
* `servers`: (Required) undefined
### client_settings
Settings that admins can apply to the Client.

* `enabled`:  (Optional)  default value `false` Enable Client Settings for this Policy.
* `entitlements_list`:  (Optional)  Enum values: `Show,Hide`Show or hide Entitlement List on Client UI.
* `attention_level`:  (Optional)  Enum values: `Show,Low,Medium,High`Set the Attention Level automatically on Client and hide the option. "Show" will leave option to the user.
* `auto_start`:  (Optional)  Enum values: `Show,Enabled,Disabled`Set the Autostart setting automatically on Client and hide the option. "Show" will leave option to the user.
* `add_remove_profiles`:  (Optional)  Enum values: `Show,Hide`Allow adding and removing profiles on Client.
* `keep_me_signed_in`:  (Optional)  Enum values: `Show,Enabled,Disabled`Set the "Keep me signed-in" setting for credential providers automatically on Client and hide the option. "Show" will leave option to the user.
* `saml_auto_sign_in`:  (Optional)  Enum values: `Show,Enabled,Disabled`Set the "SAML/Certificate auto sign-in" setting automatically on Client and hide the option. "Show" will leave option the user.
* `quit`:  (Optional)  Enum values: `Show,Hide`Show or hide "Quit" on Client UI.
* `sign_out`:  (Optional)  Enum values: `Show,Hide`Show or hide "Sign out" on Client UI.
* `suspend`:  (Optional)  Enum values: `Show,Hide`Show or hide "Suspend" feature on Client UI.
### administrative_roles
List of Administrative Role IDs in this Policy.

### tags
Array of tags.




## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgatesdp_policy.example d3131f83-10d1-4abc-ac0b-7349538e8300
```
