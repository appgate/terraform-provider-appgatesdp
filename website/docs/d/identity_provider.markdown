---
layout: "appgate"
page_title: "APPGATE: appgate_identity_provider"
sidebar_current: "docs-appgate-datasource-identity_provider"
description: |-
  The identity_provider data source provides details about a specific identity_provider.
---

# appgate_identity_provider

The identity_provider data source provides details about a specific identity_provider.


## Example Usage

```hcl

variable "identity_provider_id" {}

data "appgate_identity_provider" "default_identity_provider" {
    identity_provider_id = "${var.identity_provider_id}"
}

```

## Argument Reference

* identity_provider_id - (Optional) ID of identity_provider
* identity_provider_name - (Optional) Name of identity_provider
