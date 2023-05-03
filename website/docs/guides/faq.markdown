---
layout: "appgatesdp"
page_title: "Frequently Asked Questions"
sidebar_current: "docs-appgatesdp-guide-faq"
description: |-
  Frequently Asked Questions
---

# Frequently Asked Questions


## Controller warning API version deprecated message.


```
API version 'application/vnd.appgate.peer-v13+json' is still in use by "Terraform/1.4.6 (+https://www.terraform.io) Terraform-Plugin-SDK/2.10.1 appgatesdp/dev - 192.168.100.142" (current version is 'application/vnd.appgate.peer-v17+json'). Please update any external scripts as support for this version will be removed.
```

If you see error messages like these, it's most likely because the provider config is missing [client version](https://registry.terraform.io/providers/appgate/appgatesdp/latest/docs#client_version) and it's computed automatically. This error message is harmless and will be generated each time you run `terraform apply|refresh|destroy`. If you want to avoid this error message, add `client_version` to your provider configuration or set the environment variable `APPGATE_CLIENT_VERSION`.


