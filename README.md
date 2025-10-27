# Appgate SDP Terraform Provider

This repository contains the official Terraform provider for [Appgate SDP](https://www.appgate.com/software-defined-perimeter), enabling you to manage your SDP infrastructure as code.

## 🔧 Purpose

Our goal is to provide first-class support for **the latest version of Appgate SDP**, with compatibility and maintenance extending to the **two most recent versions** as well.

### ✅ Supported Versions

We actively maintain compatibility with the latest **three** SDP versions:

| Appgate SDP Version | API Version |
| ------------------- | ----------- |
| 6.5 (latest)        | v22         |
| 6.4                 | v21         |
| 6.3                 | v20         |

> Earlier versions may still be available, but they are not guaranteed to receive further updates or support.

---

## 🤝 Contributing

We welcome contributions from the community!

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) v0.12.26
- [Go](https://golang.org/doc/install) 1.20 (to build the provider plugin)



Building the provider
---------------------------


```sh
$ make build
```

Using the provider
---------------------------

Detailed documentation for the Appgate provider can be found in the docs directory, [here](./website/docs).

Examples how to deploy Appgate to cloud platforms can be found [here](https://github.com/appgate/sdp-tf-reference-architecture).



Testing the provider
---------------------------


```sh
$ make test
```

Example how to run acceptance test on an existing Appgate environment.
```bash
APPGATE_ADDRESS="https://envy-10-97-168-40.devops:8443/admin" \
APPGATE_USERNAME="admin" \
APPGATE_PASSWORD="admin" \
make testacc
```

test 1 acceptance test, for example
```bash
TF_ACC=1 \
APPGATE_ADDRESS="https://ec2-54-80-224-21.compute-1.amazonaws.com:8443/admin" \
APPGATE_USERNAME="admin" \
APPGATE_PASSWORD="admin" \
go test -v -timeout 120m github.com/appgate/terraform-provider-appgatesdp/appgate -run '^(TestAccApplianceBasicController)$'
```