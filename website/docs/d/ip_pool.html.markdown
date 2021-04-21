---
layout: "appgate"
page_title: "APPGATE: appgatesdp_ip_pool"
sidebar_current: "docs-appgate-datasource-ip_pool"
description: |-
  The ip_pool data source provides details about a specific ip_pool.
---

# appgatesdp_ip_pool

The ip_pool data source provides details about a specific ip_pool.


## Example Usage

```hcl

variable "ip_pool_id" {}

data "appgatesdp_ip_pool" "default_ip_pool" {
    ip_pool_id = "${var.ip_pool_id}"
}

```

## Argument Reference

* ip_pool_id - (Optional) ID of ip_pool
* ip_pool_name - (Optional) Name of ip_pool
