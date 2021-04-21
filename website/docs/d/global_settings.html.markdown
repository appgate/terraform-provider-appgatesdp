---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_global_settings"
sidebar_current: "docs-appgate-datasource-global_settings"
description: |-
  The global_settings data source provides details about a specific global_settings.
---

# appgatesdp_global_settings

The global_settings data source provides details about a specific global_settings.


## Example Usage

```hcl
data "appgatesdp_global_settings" "default_global_settings" {}
```

## Attributes Reference
* `claims_token_expiration` - Number ofof minutes the Claims Token is valid both for administrators and clients.
* `entitlement_token_expiration` - Number of minutes the Entitlement Token is valid for clients.
* `administration_token_expiration` - Number of minutes the administration Token is valid for administrators.
* `vpn_certificate_expiration` - Number of minutes the VPN certificates is valid for clients.
* `login_banner_message` - The configured message will be displayed on the login UI.
* `message_of_the_day` - The onfigured message will be displayed after a successful logging.
* `backup_api_enabled` - Whether the backup API is enabled or not.
* `has_backup_passphrase` - Whether there is a backup passphrase set or not. Deprecated as of 5.0. Use backupApiEnabled instead.
* `fips` -  FIPS 140-2 Compliant Tunneling.
* `geo_ip_updates` - Whether the automatic GeoIp updates are enabled or not.
* `audit_log_persistence_mode` - Audit Log persistence mode
* `app_discovery_domains` - Domains to monitor for for App Discovery feature.
* `collective_id` - A randomly generated ID during first installation to identify the Collective.
