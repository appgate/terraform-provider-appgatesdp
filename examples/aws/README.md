# Terraform example to deploy appgate on AWS


This directory is an example how to deploy AppGate with one controller, and one gateway on AWS using a pre-existing AMI published by appgate on the marketplace.


We need to do this in 2 steps, since the terraform_appgate provider is dependent
on the phsyical resources on aws.

```
cd physical
terraform init
terraform apply \
    -var 'private_key=/path/to/a/ssh/key' \
    -var 'public_key=/path/to/a/ssh/public_key' \
    -auto-approve
cd ..
# The terraform output should print the DNS names of the controller and gateway, these
# will be used in the next step.
```

```
cd appgate-resources
# make sure you have the latest build of the terraform-provider-appgate_v* in this directory.
# which can be downloaded at
# https://code.cyxtera.com/appgate/terraform-provider-appgate/releases

terraform init

terraform plan \
    -var 'controller_dns=ec2-3-86-14-28.compute-1.amazonaws.com' \
    -var 'gateway_dns=ec2-100-26-207-49.compute-1.amazonaws.com' \
    -var 'private_key=/path/to/a/ssh/key'


terraform apply \
    -var 'controller_dns=ec2-3-86-14-28.compute-1.amazonaws.com' \
    -var 'gateway_dns=ec2-100-26-207-49.compute-1.amazonaws.com' \
    -var 'private_key=/path/to/a/ssh/key' \
    -auto-approve

```
