Terraform Provider for Appgate
==================

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) v0.12.19
- [Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)



Building the provider
---------------------------


```sh
$ make build
```

Using the provider
---------------------------

Detailed documentation for the Appgate provider can be found in the docs directory, [here](./website/docs).

A detailed example how to deploy Appgate to AWS can be found [here](./examples/aws).



Testing the provider
---------------------------


```sh
$ make test
```

Example how to run acceptance test on an existing appgate environment.
```bash
APPGATE_ADDRESS="https://envy-10-97-168-40.devops:444/admin" \
APPGATE_USERNAME="admin" \
APPGATE_PASSWORD="admin" \
make testacc
```

test 1 acceptance test, for example
```bash
TF_ACC=1 \
APPGATE_ADDRESS="https://ec2-54-80-224-21.compute-1.amazonaws.com:444/admin" \
APPGATE_USERNAME="admin" \
APPGATE_PASSWORD="admin" \
go test -v -timeout 120m github.com/appgate/terraform-provider-appgate/appgate -run '^(TestAccApplianceBasicController)$'
```
