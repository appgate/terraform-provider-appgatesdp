---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_device_policy"
sidebar_current: "docs-appgate-resource-device-policy"
description: |-
   Create a new Device Policy.
---

# appgatesdp_device_policy

Device policy controls which are pushed to the Client (Ringfence rules, Client profile settings, etc).



## Example Usage

```hcl


resource "appgatesdp_device_policy" "test_device_policy" {
    name  = "test device policy"
    notes = "terraform policy notes"
    tags = [
        "terraform",
        "api-created"
    ]

    ringfence_rule_links = [
        "developer"
    ]
	ringfence_rules = [
		data.appgatesdp_ringfence_rule.default_ringfence_rule.id
	]
    tamper_proofing = true
    proxy_auto_config {
        enabled = true
        url     = "http://appgate.com"
        persist = false
    }
    trusted_network_check {
        enabled    = true
        dns_suffix = "aa"
    }
	client_settings {
		enabled           = true
		entitlements_list = "Hide"
		quit              = "Hide"
	}
}


```


## Argument Reference

The following arguments are supported:


* `disabled`: (Optional) If true, the Policy will be disregarded during authorization.
* `expression`: (Required) A JavaScript expression that returns boolean. Criteria Scripts may be used by calling them as functions.
* `type`: (Computed) Type of the Policy. It is informational and not enforced.
* `entitlements`: (Optional) List of Entitlement IDs in this Policy.
* `entitlement_links`: (Optional) List of Entitlement tags in this Policy.
* `policy_id`: (Computed) ID of the object.
* `name`: (Required) Name of the object.
* `notes`: (Optional) Notes for the object. Used for documentation purposes.
* `tags`: (Optional) Array of tags.
* `proxy_auto_config`: (Optional) Client configures PAC URL on the client OS.
* `trusted_network_check`: (Optional) Client suspends operations when it's in a trusted network.
* `client_settings`: (Optional) Settings that admins can apply to the Client.


### proxy_auto_config
Client configures PAC URL on the client OS.

* `enabled`:  (Optional)  default value `false`
* `url`:  (Optional) The URL to set on the Client OS. Example: https://pac.company.com/file.pac.
* `persist`:  (Optional) If true Client will leave the PAC URL configured after signing out.
### trusted_network_check
Client suspends operations when it's in a trusted network.

* `enabled`:  (Optional)  default value `false`
* `dns_suffix`:  (Optional) Client checks if the DNS suffix has been configured on the OS to decide whether it's on a trusted network or not.


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



### tags
Array of tags.




## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgatesdp_device_policy.example d3131f83-10d1-4abc-ac0b-7349538e8300
```
