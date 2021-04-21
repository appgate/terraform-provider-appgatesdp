---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_mfa_provider"
sidebar_current: "docs-appgate-datasource-mfa_provider"
description: |-
  The mfa_provider data source provides details about a specific mfa_provider.
---

# appgatesdp_mfa_provider

The mfa_provider data source provides details about a specific mfa_provider.


## Example Usage

```hcl

variable "mfa_provider_id" {}

data "appgatesdp_mfa_provider" "default_mfa_provider" {
    mfa_provider_id = "${var.mfa_provider_id}"
}

```

## Argument Reference

* mfa_provider_id - (Optional) ID of mfa_provider
* mfa_provider_name - (Optional) Name of mfa_provider
