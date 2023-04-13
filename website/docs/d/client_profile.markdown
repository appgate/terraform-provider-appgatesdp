---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_client_profile"
sidebar_current: "docs-appgate-datasource-client_profile"
description: |-
  The client_profile data source provides details about a specific client_profile.
---

# appgatesdp_client_profile

The client_profile data source provides details about a specific client_profile.


## Example Usage

```hcl

data "appgatesdp_client_profile" "portal" {
  client_profile_name = "portal"
}

```

## Argument Reference

* client_profile_id - (Optional) ID of client_profile
* client_profile_name - (Optional) Name of client_profile


## Read-Only

* url - Connection URL for the Client Profile.
