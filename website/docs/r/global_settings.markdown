---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_global_settings"
sidebar_current: "docs-appgate-resource-global_settings"
description: |-
   Update all Global Settings.
---

# appgatesdp_global_settings

Update all Global Settings.

~> **NOTE:**  Global settings are a singleton resource to allow us to update global settings for the collective.


~> **NOTE:**  The resource documentation is based on the latest available appgate sdp appliance version, which currently is 6.0.0
Some attributes may not be available if you are running an older version, if you try to use an attribute block that is not permitted in your current version, you will be prompted by an error message.


## Example Usage

```hcl


resource "appgatesdp_global_settings" "test_global_settings" {
  login_banner_message = "Welcome"
}


```


## Argument Reference

The following arguments are supported:


* `claims_token_expiration`: (Optional) Number of minutes the Claims Token is valid both for administrators and clients.
* `entitlement_token_expiration`: (Optional) Number of minutes the Entitlement Token is valid for clients.
* `administration_token_expiration`: (Optional) Number of minutes the administration Token is valid for administrators.
* `vpn_certificate_expiration`: (Optional) Number of minutes the VPN certificates is valid for clients.
* `spa_mode`: (Optional) SPA mode.
* `login_banner_message`: (Optional) The configured message will be displayed on the login UI.
* `message_of_the_day`: (Optional) The configured message will be displayed after a successful login.
* `backup_api_enabled`: (Optional) Whether the backup API is enabled or not.
* `backup_passphrase`: (Optional) The passphrase to encrypt Appliance Backups when backup API is used.
* `fips`: (Optional) FIPS 140-2 Compliant Tunneling.
* `geo_ip_updates`: (Optional) Whether the automatic GeoIp updates are enabled or not.
* `audit_log_persistence_mode`: (Optional) Audit Log persistence mode.
* `app_discovery_domains`: (Optional) Domains to monitor for App Discovery feature.
* `registered_device_expiration_days`: (Optional) Number of days registered devices are kept in storage before being deleted.
* `collective_id`: (Optional) A randomly generated ID during first installation to identify the Collective.


### app_discovery_domains
Domains to monitor for for App Discovery feature.




## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgatesdp_global_settings.example d3131f83-10d1-4abc-ac0b-7349538e8300
```
