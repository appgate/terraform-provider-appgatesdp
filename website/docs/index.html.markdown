---
layout: "appgatesdp"
page_title: "Provider: Appgate"
description: |-

---

# Appgate

## Example Usage

```hcl
# Configure the Appgate Provider
provider "appgatesdp" {
  username = "admin"
  password = "admin"
  url      = "https://controller.devops:444/admin"
  provider = "local"
  insecure = true
}


```

## Authentication

The Appgate provider offers a flexible means of providing credentials for
authentication. The following methods are supported, in this order, and
explained below:

- Static credentials
- Environment variables
- Config file
- Bearer Token

### Static credentials

!> **Warning:** Hard-coding credentials into any Terraform configuration is not
recommended, and risks secret leakage should this file ever be committed to a
public version control system.

Static credentials can be provided by adding an `username` and `password`
in-line in the Appgate provider block:

Usage:

```hcl
provider "appgatesdp" {
  url      = "https://appgate.controller.com:444/admin"
  username = "admin"
  password = "admin"
  provider = "local"
}
```

### Environment variables


```hcl
provider "appgatesdp" {}
```

Usage:

```sh
$ export APPGATE_ADDRESS="https://controller.devops:444/admin"
$ export APPGATE_USERNAME="admin"
$ export APPGATE_PASSWORD="admin"
$ export APPGATE_PROVIDER="local"
$ export APPGATE_INSECURE="true"
$ export APPGATE_BEARER_TOKEN="" # optional, used instead of username and password.
$ terraform plan
```

### Config file

Configure appgatesdp with a config file, can be combined with environment variables, if an `APPGATE_` environment variable is set, they take precedence over the config file.

```hcl
provider "appgatesdp" {
  config_path = var.appgate_config_file
}
```

example config file format
```json
{
    "appgate_url": "https://controller.appgate/admin",
    "appgate_username": "admin",
    "appgate_password": "admin",
    "appgate_provider": "local",
    "appgate_client_version": 15,
    "appgate_insecure": true
}

```

### Bearer Token

You can provide the Authorization Bearer token directly to the provider if you do not want to provide a username and password directly. The bearer token will subsequent be used in all resource. So its important to note that the user has the correct privileges. The bearer token can be combined with other environment variables, arguments and config file to complete the configuration of the provider. This method can be convient if you want to provision the user and authorization outside of terraform in an external program or script.


Usage:

In the example below, the token is saved to a file called `token` and exported as environment variable, you can ofcoure use it directly as environment variable.

```bash
$ cat token
eyJjbGFpbXNUb2tlbiI6ImV5SmhiR2NpT2lKU1V6VXhNaUlzSW5wcGNDSTZJa1JGUmlJc0luUjVjQ0k2SWtwWFZDSjkuZUp5dGxkZVNvOGdTaHQ5RnQ3UldlRE1SRTNFUVRoZ2hoSlhZbUlzQ0NpT004RWhNOUxzZnVuZmpuSDJBdmF2SytwTE1yQ0wvL0wwYm55VnMxR1QzWXdkU2xrd0FCL1pVUk1aN01xR3BQVWV3M0o0aDhRUlBJZ0lGYkxMNzJCWERNTUV2SGtkeGJJK2hlNVIwTWVJSHp2d2cyVDl3Z21FNUx0d3crR3FMSGc3LzVLaC9jaGhMc1Y5Y1VneGowV1JUTWVRd01VRU5Ody9CL0psRU5NQUFFWEUweEVnUVkyeUVSaWtYa3hqR3hXZ2FFeDhiQTVLNmFENHUzcy9xR1lQcUsyWVQ5KzkyaElsUW....
```

```hcl
provider "appgatesdp" {
  # this block can be empty or omitted, either provider the URL as environment variable or in a config file.
  url      = "https://appgate.controller.com:444/admin"
}
```


```bash
APPGATE_BEARER_TOKEN=`cat token` terraform apply -auto-approve
```




## Argument Reference

In addition to [generic `provider` arguments](https://www.terraform.io/docs/configuration/providers.html)
(e.g. `alias` and `version`), the following arguments are supported in the Appgate
 `provider` block:

* `config_path` - (Optional) Configure appgatesdp with a config file, if any environment variables is set, they take precedence.

* `username` - (Optional) This is the Appgate username. It must be provided, but
  it can also be sourced from the `APPGATE_USERNAME` environment variable.

* `password` - (Optional) This is the Appgate password. It must be provided, but
  it can also be sourced from the `APPGATE_PASSWORD` environment variable.

* `provider` - (Optional) This is the Appgate provider. It must be provided, but
  it can also be sourced from the `APPGATE_PROVIDER` environment variables.

* `client_version` - (Optional) This reference the appgate client SDK version, it can also be sourced from the `APPGATE_CLIENT_VERSION` environment variables. Defaults to `15`, Its not recommended to change this unless you know what you are doing.

* `insecure` - (Optional) Whether server should be accessed without verifying the TLS certificate. As the name suggests this is insecure and should not be used beyond experiments, accessing local (non-production) GHE instance etc. There is a number of ways to obtain trusted certificate for free, e.g. from Let's Encrypt. Such trusted certificate does not require this option to be enabled. Defaults to `true`.

* `debug` - (Optional) Whether HTTP request should be displayed in debug mode, combine with [TF_LOG](https://www.terraform.io/docs/internals/debugging.html) Defaults to `false`.
