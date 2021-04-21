---
layout: "appgate"
page_title: "APPGATE: appgatesdp_trusted_certificate"
sidebar_current: "docs-appgate-datasource-trusted_certificate"
description: |-
  The trusted_certificate data source provides details about a specific trusted_certificate.
---

# appgatesdp_trusted_certificate

The trusted_certificate data source provides details about a specific trusted_certificate.


## Example Usage

```hcl

variable "trusted_certificate_id" {}

data "appgatesdp_trusted_certificate" "default_trusted_certificate" {
    trusted_certificate_id = "${var.trusted_certificate_id}"
}

```

## Argument Reference

* trusted_certificate_id - (Optional) ID of trusted_certificate
* trusted_certificate_name - (Optional) Name of trusted_certificate
