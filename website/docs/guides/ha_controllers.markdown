---
layout: "appgatesdp"
page_title: "high availability controllers"
sidebar_current: "docs-appgatesdp-guide-high_availability_controllers"
description: |-
  Provisioning high availability controllers
---

## Provisioning high availability controllers


In appliance >= 5.4.0 we can't create more then 1 controller at the time,
it means that we must first create an appliance, seed it, then activate the controller function.


For a full example see:
- https://github.com/appgate/sdp-tf-reference-architecture/tree/main/deployment/aws/high-availability-controllers


#### Create an appliance resource

```hcl
resource "appgatesdp_appliance" "second_controller" {
  // ...
  // Configure a plain appliance that will become the second controller
  // You can configure everything here as normal expect the controller {} block
  lifecycle {
    ignore_changes = [
      # The following attributes will be defined and configured within
      # appgatesdp_appliance_controller_activation.activate_second_controller
      admin_interface,
      controller,
    ]
  }

}

```

#### Create an instance to become the second controller

The second step is to provision an instance with the appliance seed, for example an aws ec2 instance.
```hcl

data "appgatesdp_appliance_seed" "second_controller_seed" {
  depends_on = [
    appgatesdp_appliance.second_controller,
  ]
  appliance_id   = appgatesdp_appliance.second_controller.id
  password       = "cz"
  latest_version = true
}


resource "aws_instance" "appgatesdp_second_controller_instance" {
  // ...
  depends_on = [
    appgatesdp_appliance.second_controller,
  ]
}
```


#### Provision the instance with the appliance configuration.

Third step is to seed the instance
```hcl

resource "null_resource" "seed_controller" {
  depends_on = [
    appgatesdp_appliance.second_controller,
  ]


  connection {
    type        = "ssh"
    user        = "cz"
    timeout     = "25m"
    private_key = file(var.private_key)
    host        = aws_instance.second_controller.public_dns
  }

  provisioner "local-exec" {
    command = "echo ${data.appgatesdp_appliance_seed.second_controller_seed.seed_file} > seed.b64"
  }
  provisioner "file" {
    source      = "seed.b64"
    destination = "/home/cz/seed.b64"
  }
  provisioner "remote-exec" {
    inline = [
      "cat seed.b64 | base64 -d  | jq .  >> seed.json",
      // wait for the seed to get picked up and initialized
      "sleep 20",
      "echo OK",
    ]
  }
}
```

#### Enable controller function

Once the second controller is seeded with the default appliance configuration, we can enable the controller functionality.

```hcl
resource "appgatesdp_appliance_controller_activation" "activate_second_controller" {
  depends_on = [
    aws_instance.appgatesdp_second_controller_instance,
    null_resource.seed_controller,
  ]
  appliance_id = appgatesdp_appliance.second_controller.id
  controller {
    enabled = true
  }
  admin_interface {
    hostname   = aws_instance.appgatesdp_second_controller_instance.public_dns
    https_port = 8443
    https_ciphers = [
      "ECDHE-RSA-AES256-GCM-SHA384",
      "ECDHE-RSA-AES128-GCM-SHA256"
    ]
  }
}
```


