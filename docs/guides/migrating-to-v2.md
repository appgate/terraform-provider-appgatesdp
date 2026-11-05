---
page_title: "Migrating from 1.x to 2.0"
---

# Migrating from 1.x to 2.0

Version 2.0 is a rewrite. Schemas are now generated from the SDP Admin API OpenAPI
specification of each supported API version (v21 to v25, SDP 6.4 to 7.0), and the provider is
built on the Terraform Plugin Framework. Terraform 1.0 or later is required.

Object IDs are unchanged, so existing objects are kept. Configuration syntax and state layout
did change: follow the steps below instead of upgrading in place.

## Recommended procedure

1. Pin the 1.x provider and run `terraform plan`. It should show no changes.
2. Rewrite the configuration for 2.0 using the changes listed below.
3. Remove the resources from state, but not from SDP:
   `terraform state rm appgatesdp_condition.example` (repeat for each resource).
4. Add an `import` block for each resource. Terraform 1.5 or later is required.
   ```hcl
   import {
     to = appgatesdp_condition.example
     id = "4c07bc67-57ea-42dd-b702-c2d6c45419fc"
   }
   ```
   Singletons (`appgatesdp_global_settings`, `appgatesdp_admin_mfa_settings`) accept any import ID.
5. Upgrade the provider to `~> 2.0`, then run `terraform plan`. Adjust the configuration until
   the plan only shows the changes you intend.

## Provider configuration

| 1.x | 2.0 |
|---|---|
| `client_version = 22` | `api_version = "v22"`. `client_version` still works but is deprecated. |
| client version negotiated when unset | Unchanged. When `api_version` is unset, the newest common version is negotiated. |
| `insecure` defaulted to `true` | Defaults to `false`: TLS certificates are verified. Set `pem_filepath` to the controller CA, or `insecure = true` (`APPGATE_INSECURE=true`) to keep skipping verification. |
| `debug` | No effect. Use `TF_LOG=DEBUG`. |
| `HTTP_PROXY` used for HTTPS requests too | Standard Go proxy variables: `HTTPS_PROXY`, `HTTP_PROXY`, `NO_PROXY`. |

All other attributes and `APPGATE_*` environment variables are unchanged. The new
`APPGATE_API_VERSION` sets `api_version`.

## Syntax: nested blocks are now nested attributes

Every nested block is now an attribute that takes an object, or a list of objects:

```hcl
# 1.x
resource "appgatesdp_entitlement" "git" {
  actions {
    subtype = "tcp_up"
    action  = "allow"
    hosts   = ["10.0.0.1"]
    ports   = ["443"]
  }
}

# 2.0
resource "appgatesdp_entitlement" "git" {
  conditions = []
  actions = [{
    subtype = "tcp_up"
    action  = "allow"
    hosts   = ["10.0.0.1"]
    ports   = ["443"]
  }]
}
```

Single blocks such as the site's `vpn`, `name_resolution` and the appliance's `client_interface` become objects:
`vpn = { ... }`.

## Attribute changes

Attribute names now follow the API: the camelCase JSON property is converted to snake_case. Run
`terraform providers schema -json` or check each resource's page for the current schema.

| Change | Affected | What to do |
|---|---|---|
| `<resource>_id` attributes removed | All resources (`condition_id`, `site_id`, `policy_id`, ...) | Set `id`. A user-chosen UUID is still supported. Changing it replaces the object. |
| Data source lookup attributes renamed | All data sources (`<x>_id`, `<x>_name`) | Use `id` or `name`. |
| `type` removed from typed policies and identity providers | `appgatesdp_*_policy`, `appgatesdp_*_identity_provider` | Remove it. The provider sets it. |
| `on_boarding_two_factor` | Identity providers | Renamed `on_boarding2fa` (API `onBoarding2FA`). |
| Local file paths replaced by content | `appgatesdp_appliance_customization.file`, `appgatesdp_device_script.file`, appliance P12 `content`, portal `logo` / `background_image` | Pass base64 content, for example `file = filebase64("custom.zip")`. `device_script.content` is removed. |
| Synthetic attributes removed | `appliance_customization.detect_sha256`, `checksum_sha256` | Use the API attribute `checksum`. |
| Attributes from additional endpoints | `client_profile.url`, `replication_target.registration_token`, `certificate_authority.pem`, `appliance_seed.activated`, `appliance_seed.latest_version` | Not available yet. |
| `timeouts` blocks | All resources | Removed. They had no effect. |
| Values validated exactly as specified | For example, entitlement `risk_sensitivity` | Use the API spelling (`"High"`, not `"high"`). 1.x silently corrected the case. |
| `tags` and similar lists are no longer lowercased | All | If the controller only lowercases or reorders a list of strings, the provider keeps your values, so no diff appears. |

Optional attributes that the API fills in with a server value are now Optional and Computed.
Removing such an attribute from the configuration keeps the current value instead of resetting
it. Set the default value explicitly to reset it.

## Resources with changed behavior

| Resource | Change from 1.x |
|---|---|
| `appgatesdp_license` | Attributes follow the API's `License` object (e.g. `request` instead of `request_code`). Destroy deletes the license by ID (`DELETE /license/{id}`). Any change replaces the license. |
| `appgatesdp_blacklist_user` | Unchanged attributes; ID is still the distinguished name. Any change replaces the entry. |
| `appgatesdp_replication_source` | `registration_token` is renamed `token` (API name). The ID is `replication_source` instead of a fixed UUID. Status attributes follow the API (`status`, `last_replication_time`, ...). Available with API v24+. |
| `appgatesdp_appliance_controller_activation` | `controller` and `admin_interface` are attributes with the same schema as `appgatesdp_appliance`. The `timeouts` block is replaced by `wait_timeout` (default `45m`). As in 1.x, running out of time is a warning. |

## Behavior changes

* **Version-specific validation happens at plan time.** For example, `block_local_dns_requests`
  on a policy is rejected for `api_version = "v22"`, and `appgatesdp_replication_target` is
  rejected below v24.
* **Fewer logins.** The login token is cached and reused. 1.x logged in for nearly every
  request.
* **Deleting an appliance** still deactivates (wipes) it first if it is activated.
* **Built-in identity providers** (`local_database`, `connector`) are adopted by name on
  create and only removed from state on destroy, as in 1.x.
* **Deleting global settings or admin MFA settings** resets them to the controller defaults.
* The `state-migrate` tool, which renamed the provider from `appgate-sdp` before 0.5, has been
  removed.
