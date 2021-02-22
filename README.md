Terraform Provider for Appgate
==================

Version compatibility
---------------------------
Our goal is to always support the latest stable release of Appgate.

The current version of the master branch supports Appgate appliance API version 13.

|                         	|  client version 12 	| **client version 13** 	| client version 14 	|
|-------------------------	|--------------------	|-------------------	    |-------------------	|
| Appgate SDP 5.1.*     	| Full support       	|                   	    |                   	|
| **Appgate SDP 5.2.***  	| Partial support    	| Full support      	    |                   	|
| Appgate SDP 5.3.*     	| Partial support   	| Partial support   	    | Partial support      	|




Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) v0.12.26
- [Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)



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
