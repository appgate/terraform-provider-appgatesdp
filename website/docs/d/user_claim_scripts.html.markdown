---
layout: "appgate"
page_title: "APPGATE: appgate_user_claim_scripts"
sidebar_current: "docs-appgate-datasource-user_claim_scripts"
description: |-
  The user_claim_scripts data source provides details about a specific user_claim_scripts.
---

# appgate_user_claim_scripts

The user_claim_scripts data source provides details about a specific user_claim_scripts.


## Example Usage

```hcl

variable "user_claim_scripts_id" {}

data "appgate_user_claim_scripts" "default_user_claim_scripts" {
    user_claim_scripts_id = "${var.user_claim_scripts_id}"
}

```

## Argument Reference

* user_claim_scripts_id - (Optional) ID of user_claim_scripts
* user_claim_scripts_name - (Optional) Name of user_claim_scripts
