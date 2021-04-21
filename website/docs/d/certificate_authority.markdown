---
layout: "appgate"
page_title: "APPGATE: appgatesdp_certificate_authority"
sidebar_current: "docs-appgate-datasource-certificate_authority"
description: |-
  The certificate_authority data source provides details about a specific certificate_authority.
---

# appgatesdp_certificate_authority

Get the current CA Certificate.



## Example Usage

```hcl

data "appgatesdp_certificate_authority" "ca" {
  pem = true
}

```

## Argument Reference

* pem - (Optional) (bool) Get the current CA Certificate in PEM format. defaults to false.


## Attributes Reference

* version - X.509 certificate version.
* serial  - X.509 certificate serial number.
* issuer  - The issuer name of the certificate.
* subject - The subject name of the certificate.
* valid_from  - Since when the certificate is valid from.
* valid_to - Until when the certificate is valid.
* fingerprint - SHA256 fingerprint of the certificate.
* certificate - Base64 encoded binary of the certificate. Either DER or PEM formatted depending on the request.
* subject_public_key - Base64 encoded public key of the certificate.
