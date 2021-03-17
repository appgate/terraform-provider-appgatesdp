---
layout: "appgate"
page_title: "APPGATE: appgate_global_settings"
sidebar_current: "docs-appgate-resource-global_settings"
description: |-
   Update all Global Settings.
---

# appgate_global_settings

Global settings are a singleton resource to allow us to update global settings for the collective.

## Example Usage

```hcl

resource "appgate_global_settings" "test_global_settings" {
   login_banner_message = "Welcome"
}

```

## Argument Reference

The following arguments are supported:


* `claims_token_expiration`: (Optional) Number of minutes the Claims Token is valid both for administrators and clients.
* `entitlement_token_expiration`: (Optional) Number of minutes the Entitlement Token is valid for clients.
* `administration_token_expiration`: (Optional) Number of minutes the administration Token is valid for administrators.
* `vpn_certificate_expiration`: (Optional) Number of minutes the VPN certificates is valid for clients.
* `login_banner_message`: (Optional) The configured message will be displayed on the login UI.
* `message_of_the_day`: (Optional) The configured message will be displayed after a successful loging.
* `backup_api_enabled`: (Optional) Whether the backup API is enabled or not.
* `backup_passphrase`: (Optional) The passphrase to encrypt Appliance Backups when backup API is used.
* `fips`: (Optional) FIPS 140-2 Compliant Tunneling.
* `geo_ip_updates`: (Optional) Whether the automatic GeoIp updates are enabled or not.
* `audit_log_persistence_mode`: (Optional) Audit Log persistence mode.
* `app_discovery_domains`: (Optional) Domains to monitor for for App Discovery feature.
* `collective_id`: (Optional) A randomly generated ID during first installation to identify the Collective.


### app_discovery_domains
Domains to monitor for for App Discovery feature.
