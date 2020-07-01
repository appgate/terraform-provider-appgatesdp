---
layout: "appgate"
page_title: "APPGATE: appgate_site"
sidebar_current: "docs-appgate-datasource-site"
description: |-
  The site data source provides details about a specific site.
---

# appgate_site

The site data source provides details about a specific site.


## Example Usage

```hcl

variable "site_id" {}

data "appgate_site" "default_site" {
    site_id = "${var.site_id}"
}

```

## Argument Reference

* site_id - (Optional) ID of site
* site_name - (Optional) Name of site
