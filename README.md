# AppGate SDP Terraform Provider

This repository contains the official Terraform provider for [AppGate SDP](https://www.appgate.com/products/zero-trust-network-access), enabling you to manage your SDP infrastructure as code.

## 🔧 Purpose

Our goal is to provide first-class support for **the latest version of AppGate SDP**, with compatibility and maintenance extending to the **two most recent versions** as well.

Earlier versions of AppGate SDP may still be compatible, but this is not guaranteed.

---

## 🤝 Contributing

We welcome contributions from the community!

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html)
- [Go](https://golang.org/doc/install)



Building the provider
---------------------------


```sh
$ make build
```

Using the provider
---------------------------

Detailed documentation for the AppGate SDP provider can be found in the docs directory, [here](./website/docs).

Examples how to deploy AppGate SDP to cloud platforms can be found [here](https://github.com/appgate/sdp-tf-reference-architecture).



Testing the provider
---------------------------


```sh
$ make test
```

Example how to run acceptance test on an existing AppGate environment.
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
