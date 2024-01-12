---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_local_user"
sidebar_current: "docs-appgate-resource-local_user"
description: |-
   Create a new Local User.
---

# appgatesdp_local_user

Create a new Local User.

~> **NOTE:**  The resource documentation is based on the latest available appgate sdp appliance version, which currently is 6.2
Some attributes may not be available if you are running an older version, if you try to use an attribute block that is not permitted in your current version, you will be prompted by an error message.


## Example Usage

```hcl


resource "appgatesdp_local_user" "api_user" {
  name                  = "apiuser"
  first_name            = "john"
  last_name             = "doe"
  password              = "hunter3"
  email                 = "john.doe@test.com"
  phone                 = "+1-202-555-0172"
  failed_login_attempts = 30
  lock_start            = "2020-04-27T09:51:03Z"
  tags = [
    "terraform",
    "api-created"
  ]
}

```


## Argument Reference

The following arguments are supported:


* `first_name`: (Required) First name of the user. May be used as claim.
* `last_name`: (Required) Last name of the user. May be used as claim.
* `password`: (Required) Password for the user. Omit the field to keep the old password when updating a user.
* `email`: (Optional) E-mail address for the user. May be used as claim.
* `phone`: (Optional) Phone number for the user. May be used as claim.
* `failed_login_attempts`: (Optional) Number of wrong password login attempts since last successiful login.
* `lock_start`: (Optional) The date time when the user got locked out. A local user is locked out of the system after 5 consecutive failed login attempts. The lock is in effect for 1 minute. When the user logs in successfully, this field becomes null.
* `local_user_id`: (Optional) ID of the object.
* `name`: (Required) Name of the object.
* `notes`: (Optional) Notes for the object. Used for documentation purposes.
* `tags`: (Optional) Array of tags.


### tags
Array of tags.




## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgatesdp_local_user.example d3131f83-10d1-4abc-ac0b-7349538e8300
```
