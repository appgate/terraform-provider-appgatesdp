---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_local_user"
sidebar_current: "docs-appgate-datasource-local_user"
description: |-
  The local_user data source provides details about a specific local_user.
---

# appgatesdp_local_user

The local_user data source provides details about a specific local_user.


## Example Usage

```hcl

variable "local_user_id" {}

data "appgatesdp_local_user" "default_local_user" {
    local_user_id = "${var.local_user_id}"
}

```

## Argument Reference

* local_user_id - (Optional) ID of local_user
* local_user_name - (Optional) Name of local_user
