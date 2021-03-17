---
layout: "appgate"
page_title: "APPGATE: appgate_trusted_certificate"
sidebar_current: "docs-appgate-resource-trusted_certificate"
description: |-
   Create a new Trusted Certificate.
---

# appgate_trusted_certificate

Create a new Trusted Certificate.

## Example Usage

```hcl

resource "appgate_trusted_certificate" "cert" {
  name = "cli"
  tags = [
    "terraform",
    "api-created"
  ]
  pem = <<-EOF
-----BEGIN CERTIFICATE-----
......
-----END CERTIFICATE-----
EOF
}

```

## Argument Reference

The following arguments are supported:


* `pem`: (Required) A certificate in PEM format.
* `details`: (Optional) X509 certificate details.
* `id`: (Required) ID of the object.
* `name`: (Required) Name of the object.
* `notes`: (Optional) Notes for the object. Used for documentation purposes.
* `tags`: (Optional) Array of tags.


### tags
Array of tags.




## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgate_trusted_certificate d3131f83-10d1-4abc-ac0b-7349538e8300
```
