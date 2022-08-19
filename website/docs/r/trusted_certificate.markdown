---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_trusted_certificate"
sidebar_current: "docs-appgate-resource-trusted_certificate"
description: |-
   Create a new Trusted Certificate.
---

# appgatesdp_trusted_certificate

Create a new Trusted Certificate.

~> **NOTE:**  The resource documentation is based on the latest available appgate sdp appliance version, which currently is 6.0.0
Some attributes may not be available if you are running an older version, if you try to use an attribute block that is not permitted in your current version, you will be prompted by an error message.


## Example Usage

```hcl


resource "appgatesdp_trusted_certificate" "cert" {
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
* `trusted_certificate_id`: (Optional) ID of the object.
* `name`: (Required) Name of the object.
* `notes`: (Optional) Notes for the object. Used for documentation purposes.
* `tags`: (Optional) Array of tags.


### tags
Array of tags.




## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgatesdp_trusted_certificate.example d3131f83-10d1-4abc-ac0b-7349538e8300
```
